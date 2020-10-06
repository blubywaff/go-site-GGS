package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func webgame(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/login/", http.StatusSeeOther)
	}
	username := getUser(w, req)
	if !containsPlayer(bson.D{{"Username", username}, {"IsTraining", true}}) {
		tpls.ExecuteTemplate(w, "newplayer.gohtml", nil)
		return
	}
	if !containsPlayer(bson.D{{"Username", username}, {"IsTraining", false}}) {
		tpls.ExecuteTemplate(w, "gamestart.gohtml", nil)
		return
	}
	//playerT := readPlayer(bson.D{{"Username", username}, {"IsTraining", true}})
	//player := readPlayer(bson.D{{"Username", username}, {"IsTraining", false}})

}

func training(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/login/", http.StatusSeeOther)
	}
	// TODO add code to create db entry for training
	tpls.ExecuteTemplate(w, "trainingground.gohtml", nil)
}

func gamestart(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/login/", http.StatusSeeOther)
	}
	// TODO create db entry for new base
	tpls.ExecuteTemplate(w, "gamestart.gohtml", nil)
}

func gamedetails(w http.ResponseWriter, req *http.Request) {
	tpls.ExecuteTemplate(w, "gamedetails.gohtml", nil)
}
