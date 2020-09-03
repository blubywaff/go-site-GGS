package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

var tpls *template.Template

var badChars map[string]string

func getTime() string {
	return time.Now().Format("01/02/2006 at 15:04:05 in timezone: MST -0700")
}

func init() {
	fm := template.FuncMap{
		"time": getTime,
	}

	tpls = template.Must(template.New("").Funcs(fm).ParseGlob("templates/*.gohtml"))

	badChars = map[string]string{
		"#": "%23",
		" ": "%20",
	}

	dbSessionsCleaned = time.Now()

}

func serveTime(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Beamer is Awesome", "Message sent by Beamer Boy")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//indexTpl.ExecuteTemplate(w, "index.gohtml", time.Now().Format("01/02/2006 at 15:04:05 in timezone: MST -0700 "))
	tpls.ExecuteTemplate(w, "time.gohtml", req.FormValue("value"))
}

func favicon(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "site/favicon.ico")
}

func cookie(w http.ResponseWriter, req *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "general",
		Value:  "grevious",
		MaxAge: 60,
	})
}

func main() {

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./site"))))
	http.HandleFunc("/time/", serveTime)
	http.HandleFunc("/favicon.ico", favicon)
	http.HandleFunc("/sendfile", sendfile)
	http.HandleFunc("/recieved/", returnRecieved)
	http.HandleFunc("/cookies/", cookie)
	http.HandleFunc("/signup/", signUp)
	http.HandleFunc("/login/", login)
	http.HandleFunc("/logout/", logout)

	log.Fatal(http.ListenAndServe(":80", nil))
}
