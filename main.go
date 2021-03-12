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
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var SUPER_DEBUG_MODE_OVERRIDE bool

var tpls *template.Template

var badChars map[string]string

var ctx context.Context

func getTime() string {
	return time.Now().Format("01/02/2006 at 15:04:05 in timezone: MST -0700")
}

func shortifyNumber(number int) string {
	magnitude := (len(strconv.Itoa(number)) - 1) / 3
	if number < 0 {
		magnitude = (len(strconv.Itoa(number)) - 2) / 3
	}

	magmap := map[int]string{
		0:  "",  // normal
		1:  "k", // thousand
		2:  "m", // million
		3:  "b", // billion
		4:  "t", // trillion
		5:  "q", // quadrillion
		6:  "Q", // quintillion
		7:  "s", // sextillion
		8:  "S", // septillion
		9:  "o", // octillion
		10: "n", // nonillion
		11: "d", // decillion
		12: "u", // undecillion
		13: "D", // duodecillion
		14: "T", // tredecillion
	}
	for k, v := range magmap {
		if magnitude == k {
			return fmt.Sprintf("%5.1f"+v, float64(number)/math.Pow(1000.0, float64(magnitude)))
		}
	}
	return strconv.Itoa(number)
}

func userifyTime(timeIn time.Time) string {
	dif := time.Now().Sub(timeIn)
	sec := dif.Seconds()
	unit := 0.0
	res := ""
	skip := false

	if sec < 60 {
		unit = sec
		res = fmt.Sprintf("%.f seconds ago", sec)
		skip = true
	}
	min := dif.Minutes()
	if min < 60 && !skip {
		unit = min
		res = fmt.Sprintf("%.f minutes ago", min)
		skip = true
	}
	hour := dif.Hours()
	if hour < 24 && !skip {
		unit = hour
		res = fmt.Sprintf("%.f hours ago", hour)
		skip = true
	}
	if !skip {
		day := hour / 24
		if day < 31 {
			unit = day
			res = fmt.Sprintf("%.f days ago", day)
			skip = true
		}
		mon := day / 31
		if !skip && mon < 12 {
			unit = mon
			res = fmt.Sprintf("%.f months ago", mon)
			skip = true
		}
		yr := day / 365
		if !skip && yr < 10 {
			unit = yr
			res = fmt.Sprintf("%.f years ago", yr)
			skip = true
		}
		dc := yr / 10
		if !skip && dc < 10 {
			unit = dc
			res = fmt.Sprintf("%.f decades ago", dc)
			skip = true
		}
		ct := dc / 10
		if !skip && ct < 10 {
			unit = ct
			res = fmt.Sprintf("%.f centuries ago", ct)
			skip = true
		}
	}

	if unit > 0 && unit < 2 {
		res = res[0:len(res)-5] + res[len(res)-4:]
	}

	if res == "" {
		res = timeIn.String()
	}
	return res
}

func init() {
	//SUPER_DEBUG_MODE_OVERRIDE = os.Args[1] == "-debug"
	if SUPER_DEBUG_MODE_OVERRIDE {
		fmt.Println("Debug: Enabled")
	}
	fmt.Println(SUPER_DEBUG_MODE_OVERRIDE)
	ctx = context.Background()

	fm := template.FuncMap{
		"time":     getTime,
		"shortify": shortifyNumber,
		"timeify":  userifyTime,
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
	mux.HandleFunc("/forum/comment/read/", readComments)
	mux.HandleFunc("/webgame/", webgame)
	mux.HandleFunc("/webgame/start/", gamestart)
	mux.HandleFunc("/webgame/training/", training)
	mux.HandleFunc("/webgame/details/", gamedetails)
	mux.HandleFunc("/webgame/raid/", gameraid)
	mux.HandleFunc("/forum/data/", forumData)

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
		//log.Fatal(http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { http.Redirect(w, req, "https://"+req.Host, 307) })))
	})
	defer handleTimer.Stop()

	log.Fatal(http.ListenAndServeTLS(":443", "TLS/cert.pem", "TLS/privkey.pem", handlerfunc))
}
