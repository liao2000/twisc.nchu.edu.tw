package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// PageData is a type for filling HTML template
type PageData struct {
	Title   string
	Isindex bool
	IsLogin bool
	Main    template.HTML
	Time    int64
	Year    int
}

func initPageData() *PageData {
	data := new(PageData)
	data.Isindex = false // default value
	data.Time = time.Now().Unix() >> 10
	data.Year, _, _ = time.Now().Date()
	return data
}

// BasicWebHandler is a handler for handling url whose prefix is /
func BasicWebHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// Handle static file
	staticFiles := []string{"/robot.txt", "/sitemap.xml", "/favicon.ico"}
	for _, f := range staticFiles {
		if r.URL.Path == f {
			http.StripPrefix("/", http.FileServer(http.Dir("./"))).ServeHTTP(w, r)
			return
		}
	}

	// Handle simple web and check login info
	data := initPageData()

	user := CheckLoginBySession(w, r)
	data.IsLogin = user != nil

	var simpleWeb = map[string]string{
		"/about":             "關於",
		"/technic":           "技術研發",
		"/academy":           "學術活動",
		"/research":          "研究成果",
		"/official-document": "辦法表格",
	}

	// Handle simple web
	if title, ok := simpleWeb[r.URL.Path]; ok {
		data.Title = title
	} else {
		// Handle non simple web
		switch r.URL.Path {
		case "/":
			data.Title = "國立中興大學資通安全研究與教學中心"
			data.Isindex = true
			data.Main = RenderIndexPage()
		case "/news":
			data.Title = "最新消息"

			if id := strings.Join(r.Form["id"], ""); id != "" {
				aid, err := strconv.ParseInt(id, 10, 64)

				if err != nil {
					NotFound(w, r)
					return
				}

				uid := ""
				if data.IsLogin {
					uid = user.ID
				}

				artInfo := GetArticleByAid(aid, uid)

				// avoid /news?id=xxx
				if artInfo == nil {
					NotFound(w, r)
					return
				}

				data.Title = artInfo.Title + " | 國立中興大學資通安全研究與教學中心"
				data.Main = RenderPublicArticle(artInfo)
			} else {
				data.Title += " | 國立中興大學資通安全研究與教學中心"
				data.Main, _ = getHTML(r.URL.Path)
			}
		case "/login":
			if CheckLoginBySession(w, r) != nil {
				http.Redirect(w, r, "/manage", 302)
				return
			}
			data.Title = "登入"
		case "/logout":
			ret := struct {
				Err string `json:"err"`
			}{}
			if err := Logout(w, r); err != nil {
				ret.Err = "登出失敗，重試，或使用瀏覽器清除 Cookie"
				json.NewEncoder(w).Encode(w)
				return
			}

			http.Redirect(w, r, "/", 302)
			return
		default:
			NotFound(w, r)
			return
		}
	}

	if r.URL.Path != "/" && r.URL.Path != "/news" {
		data.Title += " | 國立中興大學資通安全研究與教學中心"
		data.Main, _ = getHTML(r.URL.Path)
	}

	t, _ := template.ParseFiles("./html/layout.gohtml")
	t.Execute(w, data)
}