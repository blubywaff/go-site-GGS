package main

import (
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
	"time"
	"database/sql"
	"fmt"
)

var usersdb *sql.DB

var dbSessionsCleaned time.Time

const sessionLength int = 60

func test(w http.ResponseWriter, req *http.Request) {
	c, err := req.Cookie("test")
	if err == nil {
		http.SetCookie(w, &http.Cookie{Name: "test", Value: c.Value, MaxAge: 60,})
	}
	//fmt.Println(c, err)
	if req.Method == http.MethodPost {
		//fmt.Println(req.FormValue("username"))
		http.SetCookie(w, &http.Cookie{Name: "test", Value: "00000000-0000-0000-0000-000000000000", MaxAge: 60,})
		http.Redirect(w, req, "/test2", http.StatusSeeOther)
	}
	tpls.ExecuteTemplate(w, "login.gohtml", nil)
}

func test2(w http.ResponseWriter, req *http.Request) {
	c, _ := req.Cookie("test")
	c.Value = "oof"
	//fmt.Println("err", err)
	//fmt.Println("cVal", c.Value)
	//http.SetCookie(w, c)
	tpls.ExecuteTemplate(w, "signup.gohtml", nil)
}

func alreadyLoggedIn(w http.ResponseWriter, req *http.Request) bool {
	c, err := req.Cookie("session")
	if err != nil {
		return false
	}
	http.SetCookie(w, &http.Cookie{Name: "session", Value: c.Value, MaxAge: 60, Path: "/"})
	return contains(usersdb, c.Value, "sessionID", "sessions")
}

func signUp(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		p := req.FormValue("password")
		f := req.FormValue("firstname")
		l := req.FormValue("lastname")

		if contains(usersdb, un, "username", "users") {
			http.Error(w, "Username already taken", http.StatusForbidden)
			return
		}

		sID := uuid.New()
		http.SetCookie(w, &http.Cookie{Name: "session", Value: sID.String(), MaxAge: sessionLength, Path: "/"})
		bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		writeUser(un, string(bs), f, l)
		writeSession(sID.String(), un, time.Now().Format(dbTimeFormat))
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	tpls.ExecuteTemplate(w, "signup.gohtml", nil)
}

func login(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {

		un := req.FormValue("username")
		p := req.FormValue("password")

		if !contains(usersdb, un, "username", "users") {
			http.Error(w, "Invalid Username", http.StatusForbidden)
			return
		}
		err := bcrypt.CompareHashAndPassword([]byte(find(usersdb, un, "username", "users", "password")), []byte(p))
		if err != nil {
			http.Error(w, "Password and username do not match", http.StatusForbidden)
			return
		}

		sID := uuid.New()
		http.SetCookie(w, &http.Cookie{Name: "session", Value: sID.String(), MaxAge: sessionLength, Path: "/"})
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		writeSession(sID.String(), un, time.Now().Format(dbTimeFormat))
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return

	}

	tpls.ExecuteTemplate(w, "login.gohtml", nil)

}

func logout(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	c, err := req.Cookie("session")

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	remove(usersdb, "sessionID", "sessions", c.Value)

	http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1, Path: "/"})

}

func cleanSessions() {
	r, _ := usersdb.Query("select sessionID, lastActivity from sessions;")
	defer r.Close()
	for r.Next() {
		var sid string
		var last string
		r.Scan(&sid, &last)
		timeLast, _ := time.Parse(dbTimeFormat, last)
		nowtxt := time.Now().Format(dbTimeFormat)
		now, _ := time.Parse(dbTimeFormat, nowtxt)
		if now.Sub(timeLast) > (time.Second * time.Duration(sessionLength)) {
			_, err := usersdb.Exec("delete from sessions where sessionID='" + sid + "';")
			check(err)
		}
	}
}

func cleaner() {
	for now := range time.Tick(time.Second * time.Duration(5)) {
		fmt.Println(now)
		cleanSessions()
	}
}
