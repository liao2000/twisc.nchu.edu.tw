# twisc.nchu.edu.tw

國立中興大學資通安全研究與教學中心網站 [www.twisc.nchu.edu.tw](http://www.twisc.nchu.edu.tw)

![](https://i.imgur.com/Dptp0qQ.png)

## Dependencies

+ Go 1.16
+ SQLite3
+ Front-end javascript

    + jQuery (v3.5.1)

    + Text editor: [CkEditor](https://ckeditor.com/)

    + Date format (jQuery dependency): [jquery-dateFormat](https://github.com/phstc/jquery-dateFormat)

## Files
+ twisc.nchu.edu.tw/
    + .git/
    + assets/  (static files)
        + img/
        + js/
        + style/
        + upload/ (client's uploaded files)
    + config/
    + db/ (sqlite database)
        + calendar.db
        + calendar.db.sql (*only schema*)
        + main.db
        + main.db.sql (*only schema*)
    + handler/
        + basic.go (`/*`)
        + error.go (`/error/*` HTTP403 & 404)
        + api.go (`/api/*` for Ajax)
        + manage.go (`/manage/*`)
        + article.go
        + calendar.go
        + files.go (manage the uploaded files)
        + login.go
        + renderer.go
        + util.go
    + html/
    + go.mod
    + go.sum
    + __main.go__ (main program)

## Quick start

1. Create database from /db/*.sql
    + /db/calendar.db.sql
    + /db/main.db.sql
    ```sh
    cat calendar.db.sql | sqlite3 calendar.db
    cat main.db.sql | sqlite3 main.db
    ```
2. run
    ```sh
    # build
    $ go build main.go

    # run at port 8086
    $ ./main

    # run at port 443
    $ ./main -p 443

    # disable minify static files
    $ ./main --debug
    ```
3. default user
    + id: root
    + password: 00000000