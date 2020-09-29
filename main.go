package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
	"fmt"
	"strings"
	"os"
	"io"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"context"
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
		logger(err.Error())
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
			updateSession(bson.D{{Key: "SessionID", Value: c.Value}}, bson.D{{Key: "$set", Value: bson.D{{Key: "LastActivity", Value: time.Now().Format(dbTimeFormat)}}}})
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

func ping(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "OK")
}

func getLog(w http.ResponseWriter, req *http.Request) {
	tpls.ExecuteTemplate(w, "log.gohtml", logs)
}

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	check(err)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	check(err)

	usersdb = client.Database("accountdb").Collection("users")
	sessionsdb = client.Database("accountdb").Collection("sessions")
	profilePicturesdb = client.Database("accountdb").Collection("profilepictures")

	threadsdb = client.Database("forumdb").Collection("threads")
	commentsdb = client.Database("forumdb").Collection("comments")
	votesdb = client.Database("forumdb").Collection("votes")
	
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
	//mux.HandleFunc("/sendfile/", sendfile)
	//mux.HandleFunc("/recieved/", returnRecieved)
	mux.HandleFunc("/cookies/", cookie)
	mux.HandleFunc("/signup/", signUp)
	mux.HandleFunc("/login/", login)
	mux.HandleFunc("/logout/", logout)
	mux.HandleFunc("/signup/checkusername", checkUsername)
	mux.HandleFunc("/account/", account)
	mux.HandleFunc("/account/profilepicture", profilePicture)
	mux.HandleFunc("/forum/", forum)
	mux.HandleFunc("/thread/", forumThread)
	mux.HandleFunc("/createthread/", createThread)
	mux.HandleFunc("/forum/vote/", vote)
	mux.HandleFunc("/forum/comment/", createComment)

	mux.HandleFunc("/test", test)
	mux.HandleFunc("/test2", test2)

	mux.HandleFunc("/ping/", ping)
	mux.HandleFunc("/log/", getLog)

	handlerfunc := http.HandlerFunc(func (w http.ResponseWriter, req *http.Request) {
		//fmt.Println(req.URL.Path)
		c, err := req.Cookie("session")
		//fmt.Println("cerr", err)
		//fmt.Println(req.URL.Path)
		if err == nil {
			http.SetCookie(w, &http.Cookie{Name: "session", Value: c.Value, MaxAge: sessionLength, Path: "/"})
			updateSession(bson.D{{Key: "SessionID", Value: c.Value}}, bson.D{{Key: "$set", Value: bson.D{{Key: "LastActivity", Value: time.Now().Format(dbTimeFormat)}}}})
		}
		end := strings.Split(req.URL.Path, ".")[len(strings.Split(req.URL.Path, "."))-1]
		if !strings.Contains(req.URL.Path, "/recieved/") && end != "gohtml" && end != "css" && end != "js" && end != "html" {
			file := "./assets/" + strings.Split(req.URL.Path, "/")[len(strings.Split(req.URL.Path, "/"))-1]
			//fmt.Println("file", file)
			this, err := os.Stat(file)
			//fmt.Println("osstat", err, this.Name(), os.IsExist(err))
			if this != nil && this.Name() != "assets" && (os.IsExist(err) || (!os.IsExist(err) && err == nil)) {
				//fmt.Println("into ere")
				fileMux.ServeHTTP(w, req)
				return
			}
		}
		mux.ServeHTTP(w, req)
	})

	log.Fatal(http.ListenAndServeTLS(":443", "TLS/cert.pem", "TLS/privkey.pem", handlerfunc))
}
