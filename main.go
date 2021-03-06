package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/svg"
	"golang.org/x/crypto/acme/autocert"
	"twisc.nchu.edu.tw/handler"
)

func main() {
	// parse flag
	port := flag.Int("p", 8086, "Port number (default: 8086)")
	debug := flag.Bool("debug", false, "Weather activate debug mode")
	flag.Parse()

	// web server
	mux := http.NewServeMux()
	staticFolder := []string{"/assets", "/.well-known/pki-validation"}

	if debug != nil && *debug {
		fmt.Println("debugging mode!")
		for _, dir := range staticFolder {
			fileServer := http.FileServer(http.Dir("." + dir))
			mux.Handle(dir+"/", http.StripPrefix(dir, neuter(fileServer)))
		}
	} else {
		// minify static files
		m := minify.New()
		m.AddFunc("text/css", css.Minify)
		m.AddFunc("image/svg+xml", svg.Minify)
		m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)

		for _, dir := range staticFolder {
			fileServer := http.FileServer(http.Dir("." + dir))
			mux.Handle(dir+"/", http.StripPrefix(dir, m.Middleware(neuter(fileServer))))
		}
	}

	mux.HandleFunc("/api/", handler.ApiHandler)
	mux.HandleFunc("/error/", handler.ErrorWebHandler)
	mux.HandleFunc("/manage/", handler.ManageWebHandler)
	mux.HandleFunc("/", handler.BasicWebHandler)

	// TLS Manager
	tls := &autocert.Manager{
		Cache:      autocert.DirCache("./"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("twisc.nchu.edu.tw", "www.twisc.nchu.edu.tw"),
	}

	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", *port),
		Handler:   mux,
		TLSConfig: tls.TLSConfig(),
	}

	// https://stackoverflow.com/questions/32325343/go-tcp-too-many-open-files-debug
	// try to solve "too many open files debug" bug
	http.DefaultClient.Timeout = time.Minute * 10

	if *port == 443 {
		fmt.Println("https://localhost")
		go http.ListenAndServe(":80", http.HandlerFunc(redirect))
		if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Fatalln("ListenAndServe: ", err)
		}
	} else {
		fmt.Printf("http://localhost:%d\n", *port)
		if err := server.ListenAndServe(); err != nil {
			log.Fatalln("ListenAndServe: ", err)
		}
	}
}

func redirect(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

func neuter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			handler.Forbidden(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
