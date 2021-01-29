package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var tpls *template.Template

var badChars map[string]string

var ctx context.Context

func getTime() string {
	return time.Now().Format("01/02/2006 at 15:04:05 in timezone: MST -0700")
}

func createNavbar(params string) template.HTML {
	return template.HTML("<style>\n" +
		"    * {\n" +
		"        cursor: default;\n" +
		"    }\n\n" +
		"    body {\n" +
		"        overflow: visible;\n" +
		"    }\n\n" +
		"    body::-webkit-scrollbar {\n" +
		"        display: none;\n" +
		"    }\n\n" +
		"    #navbar {\n" +
		"        --background: #ff2020;\n" +
		"        --text-color: #ffffff;\n" +
		"        --text-hover-color: #cccccc;\n" +
		"    }\n\n" +
		"    #navbar {\n" +
		"        background-color: var(--background);\n" +
		"        padding: 20px;\n" +
		"        min-height: 32px;\n" +
		"        width: calc(100% - 20px);\n" +
		"        position: fixed;\n" +
		"        z-index: 999;\n" +
		"    }\n\n" +
		"    #navbar .element {\n" +
		"        display: inline-block;\n" +
		"        vertical-align: middle;\n" +
		"    }\n\n" +
		"    #navbar .element #sidenav-button {\n" +
		"    / / cursor: pointer;\n" +
		"        display: inline-block;\n" +
		"        margin-right: 20px;\n" +
		"    }\n\n" +
		"    #navbar .element #sidenav-button * {\n" +
		"        width: 20px;\n" +
		"        height: 3px;\n" +
		"        background-color: var(--text-color);\n" +
		"        margin: 4px 0;\n" +
		"        transition: transform 0.5s, opacity 0.5s;\n" +
		"        cursor: inherit;\n" +
		"    }\n\n" +
		"    #navbar .element #sidenav-button:hover * {\n" +
		"        background-color: var(--text-hover-color);\n" +
		"    }\n\n" +
		"    #navbar .element .text {\n" +
		"        font-size: x-large;\n" +
		"        margin-right: 20px;\n" +
		"        color: var(--text-color);\n" +
		"        text-decoration: none;\n" +
		"    / / cursor: pointer;\n" +
		"        display: inline-block;\n" +
		"    }\n\n" +
		"    #navbar #nav .text {\n" +
		"        float: right;\n" +
		"    }\n\n" +
		"    #navbar #account {\n" +
		"        float: right;\n" +
		"    }\n\n" +
		"    #sidenav-button.change #bar1 {\n" +
		"        transform: translate(0px, 7px) rotate(-45deg);\n" +
		"    }\n\n" +
		"    #sidenav-button.change #bar2 {\n" +
		"       opacity: 0;\n" +
		"   }\n\n" +
		"   #sidenav-button.change #bar3 {\n" +
		"       transform: translate(0px, -7px) rotate(45deg);\n" +
		"   }\n\n" +
		"   #navbar .element div.text:hover {\n" +
		"       color: var(--text-hover-color);\n" +
		"   }\n\n" +
		"   body {\n" +
		"       background-color: #ffffff;\n" +
		"   }\n\n" +
		"   #title-name {\n" +
		"       font-size: 48px;\n" +
		"       text-align: center;\n" +
		"   }\n\n" +
		"   #title-name {\n" +
		"       color: #0000ff;\n" +
		"       background-color: #00ddff;\n" +
		"       border: 5px solid #000000;\n" +
		"   }\n\n" +
		"   #sidenav {\n" +
		"       z-index: 999;\n" +
		"       transition: 0.5s;\n" +
		"       background-color: #222222;\n" +
		"       opacity: 80%;\n" +
		"       height: calc(100% - 72px);\n" +
		"       position: fixed;\n" +
		"       text-align: center;\n" +
		"       bottom: 0;\n" +
		"       overflow-x: hidden;\n" +
		"   }\n\n" +
		"   #sidenav * {\n" +
		"       opacity: 0;\n" +
		"       z-index: 2;\n" +
		"       transition: opacity 0.3s;\n" +
		"       color: #ffffff;\n" +
		"     //width: 0;\n" +
		"     //overflow-x: hidden;\n" +
		"   }\n\n" +
		"   #sidenav.change * {\n" +
		"       opacity: 100%;\n" +
		"     //width: auto;\n" +
		"   }\n\n" +
		"   #title-wrapper {\n" +
		"       padding-block: 20px;\n" +
		"       padding-inline: 20px;\n" +
		"       background: #b96fff;\n" +
		"   }\n\n" +
		"   #content {\n" +
		"       text-align: center;\n" +
		"       background: #000000;\n" +
		"   }\n\n" +
		"   #content #mission {\n" +
		"       background: #424242;\n" +
		"       height: 300px;\n" +
		"   }\n\n" +
		"   .content-wrapper {\n" +
		"       display: inline-block;\n" +
		"       transform: translateY(50px);\n" +
		"   }\n\n" +
		"   #mission .text {\n" +
		"       font: 24px \"Comic Sans MS\", sans-serif;\n" +
		"   }\n\n" +
		"   .text-wrapper {\n" +
		"       width: calc(50% - 10px);\n" +
		"       vertical-align: middle;\n" +
		"       margin-right: 5px;\n" +
		"   }\n\n" +
		"   #mission .text-wrapper {\n" +
		"       float: left;\n" +
		"   }\n\n" +
		"   #mission .text-wrapper .text {\n" +
		"       float: right;\n" +
		"       transform: translateY(50%);\n" +
		"   }\n\n" +
		"   #about-us .text-wrapper {\n" +
		"       float: right;\n" +
		"   }\n\n" +
		"   #about-us .text-wrapper .text {\n" +
		"       float: left;\n" +
		"       transform: translateY(50%);\n" +
		"   }\n\n" +
		"   .img-wrapper {\n" +
		"       width: calc(50% - 10px);\n" +
		"       margin-left: 5px;\n" +
		"   }\n\n" +
		"   #mission .img-wrapper {\n" +
		"       float: right;\n" +
		"   }\n\n\n" +
		"   #mission .img-wrapper .img {\n" +
		"       float: left;\n" +
		"   }\n\n" +
		"   #about-us .img-wrapper {\n" +
		"       float: left\n" +
		"   }\n\n" +
		"   #about-us .img-wrapper .img {\n" +
		"       float: right;\n" +
		"   }\n\n" +
		"   .clickable {\n" +
		"       cursor: pointer;\n" +
		"   }\n\n" +
		"   #about-us {\n" +
		"       height: 300px;\n" +
		"   }\n\n" +
		"   #about-us .text {\n" +
		"       font: 24px \"Comic Sans MS\", sans-serif;\n" +
		"       color: #fff;\n" +
		"   }\n\n" +
		"   .button-wrapper {\n" +
		"       margin: 20px;\n" +
		"   }\n\n" +
		"   .button-wrapper div {\n" +
		"       width: max-content;\n" +
		"   }\n\n" +
		"   #sidenav {\n" +
		"       font: 18px \"Times New Roman\";\n" +
		"   }\n" +
		"   \n" +
		"   .clickable * {\n" +
		"       cursor: inherit;\n" +
		"   }\n\n</style>")
}

func init() {
	ctx = context.Background()

	fm := template.FuncMap{
		"time":   getTime,
		"navbar": createNavbar,
	}

	tpls = template.Must(template.New("").Funcs(fm).ParseGlob("templates/*.gotpls"))
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
	http.ServeFile(w, req, "site/"+req.URL.Path)
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
	tpls.ExecuteTemplate(w, "homepage.gohtml", getUser(w, req).Username)
}

func proofos(w http.ResponseWriter, req *http.Request) {
	if strings.HasSuffix(req.URL.Path, "/proofos/") {
		tpls.ExecuteTemplate(w, "proofofskills.gohtml", nil)
		return
	}
	http.ServeFile(w, req, "site/"+req.URL.Path[9:])
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
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	check(err)

	usersdb = client.Database("accountdb").Collection("users")
	sessionsdb = client.Database("accountdb").Collection("sessions")
	profilePicturesdb = client.Database("accountdb").Collection("profilepictures")

	threadsdb = client.Database("forumdb").Collection("threads")
	commentsdb = client.Database("forumdb").Collection("comments")
	votesdb = client.Database("forumdb").Collection("votes")

	playersdb = client.Database("gamedb").Collection("players")

	timer := time.AfterFunc(time.Second, cleaner)
	defer timer.Stop()

	mux := http.NewServeMux()
	fileMux := http.NewServeMux()

	fileMux.HandleFunc("/", fileHandle)

	//http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./site"))))
	//mux.HandleFunc("/", site)
	mux.HandleFunc("/", sendHome)
	mux.HandleFunc("/home/", home)
	mux.HandleFunc("/proofos/", proofos)
	mux.HandleFunc("/aboutus/", func(w http.ResponseWriter, req *http.Request) {
		tpls.ExecuteTemplate(w, "aboutus.gohtml", nil)
	})
	mux.HandleFunc("/testpage/", func(w http.ResponseWriter, req *http.Request) {
		check(tpls.ExecuteTemplate(w, "testpage.gohtml", nil))
	})
	mux.HandleFunc("/time/", serveTime)
	mux.HandleFunc("/favicon.ico", favicon)
	//mux.HandleFunc("/sendfile/", sendfile)
	//mux.HandleFunc("/recieved/", returnRecieved)
	mux.HandleFunc("/cookies/", cookie)
	mux.HandleFunc("/signup/", signUp)
	mux.HandleFunc("/login/", login)
	mux.HandleFunc("/logout/", logout)
	mux.HandleFunc("/signup/checkusername/", checkUsername)
	mux.HandleFunc("/account/", account)
	mux.HandleFunc("/account/profilepicture/", profilePicture)
	mux.HandleFunc("/account/delete/", deleteAccount)
	mux.HandleFunc("/forum/", forum)
	mux.HandleFunc("/thread/", forumThread)
	mux.HandleFunc("/createthread/", createThread)
	mux.HandleFunc("/forum/vote/", vote)
	mux.HandleFunc("/forum/comment/", createComment)
	mux.HandleFunc("/webgame/", webgame)
	mux.HandleFunc("/webgame/start/", gamestart)
	mux.HandleFunc("/webgame/training/", training)
	mux.HandleFunc("/webgame/details/", gamedetails)
	mux.HandleFunc("/webgame/raid/", gameraid)

	mux.HandleFunc("/test", test)
	mux.HandleFunc("/test2", test2)

	mux.HandleFunc("/ping/", ping)
	mux.HandleFunc("/log/", getLog)

	handlerfunc := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		//fmt.Println(req.URL.Path)
		c, err := req.Cookie("session")
		//fmt.Println("cerr", err)
		//fmt.Println(req.URL.Path)
		if err == nil {
			if c.Value == "" {
				http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: 1, Path: "/"})
				http.Error(w, "Empty Session", http.StatusForbidden)
				return
			}
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

	handleTimer := time.AfterFunc(time.Microsecond, func() {
		log.Fatal(http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { http.Redirect(w, req, "https://"+req.Host, 307) })))
	})
	defer handleTimer.Stop()

	log.Fatal(http.ListenAndServeTLS(":443", "TLS/cert.pem", "TLS/privkey.pem", handlerfunc))
}
