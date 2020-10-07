package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func webgame(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/login/", http.StatusSeeOther)
		return
	}
	username := getUser(w, req).Username
	if !containsPlayer(bson.D{{"Username", username}, {"IsTraining", true}}) {
		tpls.ExecuteTemplate(w, "newplayer.gohtml", nil)
		return
	}
	if !containsPlayer(bson.D{{"Username", username}, {"IsTraining", false}}) {
		tpls.ExecuteTemplate(w, "gamestart.gohtml", nil)
		return
	}
	tpls.ExecuteTemplate(w, "webgame.gohtml", nil)
	//playerT := readPlayer(bson.D{{"Username", username}, {"IsTraining", true}})
	//player := readPlayer(bson.D{{"Username", username}, {"IsTraining", false}})

}

func training(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/login/", http.StatusSeeOther)
		return
	}
	username := getUser(w, req).Username
	writePlayer(Player{true, username, []Ship{}, Base{}})
	tpls.ExecuteTemplate(w, "trainingground.gohtml", nil)
}

func gamestart(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/login/", http.StatusSeeOther)
		return
	}
	username := getUser(w, req).Username
	writePlayer(Player{false, username, []Ship{}, Base{}})
	tpls.ExecuteTemplate(w, "gamestart.gohtml", nil)
}

func gamedetails(w http.ResponseWriter, req *http.Request) {
	tpls.ExecuteTemplate(w, "gamedetails.gohtml", nil)
}
