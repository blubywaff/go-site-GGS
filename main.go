package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var tpls *template.Template

func getTime() string {
	return time.Now().Format("01/02/2006 at 15:04:05 in timezone: MST -0700 ")
}

func init() {
	fm := template.FuncMap{
		"time": getTime,
	}

	tpls = template.Must(template.New("").Funcs(fm).ParseGlob("templates/*.gohtml"))
}

func serveTime(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Beamer is Awesome", "Message sent by Beamer Boy")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//indexTpl.ExecuteTemplate(w, "index.gohtml", time.Now().Format("01/02/2006 at 15:04:05 in timezone: MST -0700 "))
	tpls.ExecuteTemplate(w, "time.gohtml", req.FormValue("value"))
	fmt.Println(req)
}

func favicon(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "site/favicon.ico")
}

func sendfile(w http.ResponseWriter, req *http.Request) {

	//var s string

	if req.Method == http.MethodPost {
		f, h, err := req.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		bs, err := ioutil.ReadAll(f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dst, err := os.Create(filepath.Join("./recieved/", h.Filename))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		err = ioutil.WriteFile(filepath.Join("./recieved/", h.Filename), bs, os.ModeAppend)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		if req.FormValue("sendback") == "on" {
			http.ServeFile(w, req, filepath.Join("./recieved/", h.Filename))
		return
		}

	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpls.ExecuteTemplate(w, "sendfile.gohtml", nil)

}

func returnRecieved(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.URL.Path)
	//http.ServeFile(w, req, filepath.Join(".", req.URL.Path))
	http.ServeFile(w, req, string(req.URL.Path))
}

func main() {

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./site"))))
	http.HandleFunc("/time/", serveTime)
	http.HandleFunc("/favicon.ico", favicon)
	http.HandleFunc("/sendfile", sendfile)
	http.HandleFunc("/recieved/", returnRecieved)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
