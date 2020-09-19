package main

import (
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"fmt"
	"io/ioutil"
)

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
	return containsSession(bson.D{{"SessionID", c.Value}})
}

func signUp(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		username := req.FormValue("username")
		email := req.FormValue("email")
		firstname := req.FormValue("firstname")
		lastname := req.FormValue("lastname")
		password := req.FormValue("password")

		if containsUser(bson.D{{"Username", username}}) {
			http.Error(w, "Username already taken", http.StatusForbidden)
			return
		}

		sID := uuid.New()
		http.SetCookie(w, &http.Cookie{Name: "session", Value: sID.String(), MaxAge: sessionLength, Path: "/"})
		bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		writeUser(user{username, email, bs, firstname, lastname})
		writeSession(session{sID.String(), username, time.Now().Format(dbTimeFormat)})
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

		if !containsUser(bson.D{{"Username", un}}) {
			http.Error(w, "Invalid Username", http.StatusForbidden)
			return
		}
		err := bcrypt.CompareHashAndPassword(readUser(bson.D{{"Username", un}}).Password, []byte(p))
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
		writeSession(session{sID.String(), un, time.Now().Format(dbTimeFormat)})
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

	removeSession(bson.D{{"SessionID", c.Value}})

	http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1, Path: "/"})

}

func cleanSessions() {
	cur, err := sessionsdb.Find(context.Background(), bson.D{})
	check(err)
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		session := session{}
		err := cur.Decode(&session)
		check(err)
		timeLast, _ := time.Parse(dbTimeFormat, session.LastActivity)
		nowtxt := time.Now().Format(dbTimeFormat)
		now, _ := time.Parse(dbTimeFormat, nowtxt)
		if now.Sub(timeLast) > (time.Second * time.Duration(sessionLength)) {
			removeSession(bson.D{{"SessionID", session.SessionID}})
		}
	}
}

func cleaner() {
	for _ = range time.Tick(time.Second * time.Duration(5)) {
		//fmt.Println(now)
		cleanSessions()
	}
}

func checkUsername(w http.ResponseWriter, req *http.Request) {
	bytes, err := ioutil.ReadAll(req.Body)
	check(err)
	fmt.Fprint(w, containsUser(bson.D{{"Username", string(bytes)}}))
}
