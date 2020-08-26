package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

var indexTpl *template.Template

func getTime() string {
	return time.Now().Format("01/02/2006 at 15:04:05 in timezone: MST -0700 ")
}

func init() {
	fm := template.FuncMap{
		"time": getTime,
	}

	indexTpl = template.Must(template.New("").Funcs(fm).ParseGlob("templates/*.gohtml"))
}

func serveTime(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Beamer is Awesome", "Message sent by Beamer Boy")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//indexTpl.ExecuteTemplate(w, "index.gohtml", time.Now().Format("01/02/2006 at 15:04:05 in timezone: MST -0700 "))
	indexTpl.ExecuteTemplate(w, "index.gohtml", req.FormValue("value"))
	fmt.Println(req)
}

func favicon(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "site/favicon.ico")
}

func main() {

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./site"))))
	http.HandleFunc("/time/", serveTime)
	http.HandleFunc("/favicon.ico", favicon)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
