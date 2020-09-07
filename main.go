package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"strings"
	"os"
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
	tpls.ParseGlob("pages/*.gohtml")

	badChars = map[string]string{
		"#": "%23",
		" ": "%20",
	}

	dbSessionsCleaned = time.Now()
	timer := time.AfterFunc(time.Second, cleaner)
	defer timer.Stop()
}

func serveTime(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Beamer is Awesome", "Message sent by Beamer Boy")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//indexTpl.ExecuteTemplate(w, "index.gohtml", time.Now().Format("01/02/2006 at 15:04:05 in timezone: MST -0700 "))
	tpls.ExecuteTemplate(w, "time.gohtml", req.FormValue("value"))
}

func favicon(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "assets/favicon.ico")
}

func cookie(w http.ResponseWriter, req *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "general",
		Value:  "grevious",
		MaxAge: 60,
	})
}

func index(w http.ResponseWriter, req *http.Request) {
	redirect(w, req, "home")
}

func redirect(w http.ResponseWriter, req *http.Request, dest string) {
	c, err := req.Cookie("session")
	if err != nil {
		panic(err)
	}
	c.MaxAge = sessionLength
	http.SetCookie(w, c)
	http.Redirect(w, req, dest, http.StatusSeeOther)
}

func sessionCookie(w http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("session")
	if err != nil {
		fmt.Println(err)
		return
	}
	c.MaxAge = sessionLength
	http.SetCookie(w, c)
}

func site(w http.ResponseWriter, req *http.Request) {
	//sessionCookie(w, req)
	http.ServeFile(w, req, "site/" + req.URL.Path)
}

func sUp(h http.HandlerFunc) http.HandlerFunc {
	//fmt.Println("hello")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session")
		//fmt.Println("cerr", err)
		if err == nil {
			http.SetCookie(w, &http.Cookie{Name: "session", Value: c.Value, MaxAge: sessionLength, Path: "/"})
			updateSession(c.Value, time.Now().Format(dbTimeFormat))
		}
		h.ServeHTTP(w, r)
	})
}

func sendHome(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "/home/", http.StatusSeeOther)
}

func home(w http.ResponseWriter, req *http.Request) {
	tpls.ExecuteTemplate(w, "homepage.gohtml", nil)
}

func proofos(w http.ResponseWriter, req *http.Request) {
	if strings.HasSuffix(req.URL.Path, "/proofos/") {
		tpls.ExecuteTemplate(w, "proofofskills.gohtml", nil)
		return
	}
	http.ServeFile(w, req, "site/" + req.URL.Path[9:])
}

func fileHandle(w http.ResponseWriter, req *http.Request) {
	file := "assets/" + strings.Split(req.URL.Path, "/")[len(strings.Split(req.URL.Path, "/"))-1]
	http.ServeFile(w, req, file)
}

func main() {

	var err error
	usersdb, err = sql.Open("mysql", "root:manhin0717@tcp(localhost:3306)/testdb?charset=utf8")
	check(err)
	defer usersdb.Close()
	//err = usersdb.Ping()
	check(err)
	
	timer := time.AfterFunc(time.Second, cleaner)
	defer timer.Stop()

	mux := http.NewServeMux()
	fileMux := http.NewServeMux()

	fileMux.HandleFunc("/",  fileHandle)

	
	//http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./site"))))
	//mux.HandleFunc("/", site)
	mux.HandleFunc("/", sendHome)
	mux.HandleFunc("/home/", home)
	mux.HandleFunc("/proofos/", proofos)
	mux.HandleFunc("/time/", serveTime)
	mux.HandleFunc("/favicon.ico", favicon)
	mux.HandleFunc("/sendfile/", sendfile)
	mux.HandleFunc("/recieved/", returnRecieved)
	mux.HandleFunc("/cookies/", cookie)
	mux.HandleFunc("/signup/", signUp)
	mux.HandleFunc("/login/", login)
	mux.HandleFunc("/logout/", logout)

	mux.HandleFunc("/test", test)
	mux.HandleFunc("/test2", test2)


	log.Fatal(http.ListenAndServe(":80", http.HandlerFunc(func (w http.ResponseWriter, req *http.Request) {
		fmt.Println(req.URL.Path)
		c, err := req.Cookie("session")
		//fmt.Println("cerr", err)
		//fmt.Println(req.URL.Path)
		if err == nil {
			http.SetCookie(w, &http.Cookie{Name: "session", Value: c.Value, MaxAge: sessionLength, Path: "/"})
			updateSession(c.Value, time.Now().Format(dbTimeFormat))
		}
		end := strings.Split(req.URL.Path, ".")[len(strings.Split(req.URL.Path, "."))-1]
		if !strings.Contains(req.URL.Path, "/recieved/") && end != "gohtml" && end != "css" && end != "js" && end != "html" {
			file := "./assets/" + strings.Split(req.URL.Path, "/")[len(strings.Split(req.URL.Path, "/"))-1]
			fmt.Println("file", file)
			this, err := os.Stat(file)
			fmt.Println("osstat", err, this.Name(), os.IsExist(err))
			if this.Name() != "assets" && (os.IsExist(err) || (!os.IsExist(err) && err == nil)) {
				fmt.Println("into ere")
				fileMux.ServeHTTP(w, req)
				return
			}
		}
		mux.ServeHTTP(w, req)
	})))
	
}
