package main

import (
	"net/http"
	"github.com/google/uuid"
	"time"
	"strconv"
	"go.mongodb.org/mongo-driver/bson"
)

func forum(w http.ResponseWriter, req *http.Request) {
	tpls.ExecuteTemplate(w, "forum.gohtml", getForumData())
}

func forumThread(w http.ResponseWriter, req *http.Request) {
	threadIDQ, ok := req.URL.Query()["thread"]
	if !ok {
		http.Error(w, "No Thread ID", http.StatusBadRequest)
	}
	threadID := threadIDQ[0]
	
	if len(threadID) < 36 {
		http.Error(w, "Invalid Thread ID", http.StatusBadRequest)
		return
	}
	tpls.ExecuteTemplate(w, "thread.gohtml", getThread(threadID).getFull())
}

func createThread(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		title := req.FormValue("title")
		body := req.FormValue("body")
		id := uuid.New()

		writeThread(Thread{getUser(w, req).Username, title, time.Now(), id.String(), body, 0, []string{}})
		http.Redirect(w, req, "/thread/?thread=" + id.String(), http.StatusSeeOther)
		return
	}
	tpls.ExecuteTemplate(w, "createthread.gohtml", nil)
}

func vote(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	username := getUser(w, req).Username
	threadIDQ, okthread := req.URL.Query()["thread"]
	commentIDQ, okcomment := req.URL.Query()["comment"]
	voteQ, okvote := req.URL.Query()["vote"]
	if !okvote {
		http.Error(w, "No vote provided", http.StatusBadRequest)
	}
	var id string
	if okthread && okcomment {
		http.Error(w, "Both Thread and Comment Provided", http.StatusBadRequest)
	} else if !okthread && !okcomment {
		http.Error(w, "No Resource Provided", http.StatusBadRequest)
	} else if okthread {
		id = threadIDQ[0]
	} else if okcomment {
		id = commentIDQ[0]
	}
	voteS := voteQ[0]
	vote, err := strconv.Atoi(voteS)
	if !check(err) || vote < -1 || vote > 1 {
		http.Error(w, "Invalid Vote", http.StatusBadRequest)
	}
	if vote == 0 {
		if containsVote(username, id, okthread) {
			voteOn(id, okthread, -getVote(username, id, okthread).Vote)
			removeVote(username, id, okthread)
		}
		return
	}
	if containsVote(username, id, okthread) {
		if getVote(username, id, okthread).Vote == vote {
			removeVote(username, id, okthread)
			voteOn(id, okthread, -vote)
			return
		}
		updateVote(username, id, okthread, vote)
		voteOn(id, okthread, 2 * vote)
	} else {
		writeVote(username, id, okthread, vote)
		voteOn(id, okthread, vote)
	}
	http.Redirect(w, req, "/thread/?thread=" + id, http.StatusSeeOther)
}

func voteOn(id string, isThread bool, vote int) {
	if isThread {
		updateThread(bson.D{{Key: "ID", Value: id}}, bson.D{{Key: "$set", Value: bson.D{{Key: "Score", Value: readThread(bson.D{{Key: "ID", Value: id}}).Score + vote}}}})
	} else {
		updateComment(bson.D{{Key: "ID", Value: id}}, bson.D{{Key: "$set", Value: bson.D{{Key: "Score", Value: readComment(bson.D{{Key: "ID", Value: id}}).Score + vote}}}})
	}
}

func createComment(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		content := req.FormValue("content")
		username := getUser(w, req).Username
		threadIDQ, okthread := req.URL.Query()["thread"]
		commentIDQ, okcomment := req.URL.Query()["comment"]
		var id string
		if okthread && okcomment {
			http.Error(w, "Both Thread and Comment Provided", http.StatusBadRequest)
		} else if !okthread && !okcomment {
			http.Error(w, "No Resource Provided", http.StatusBadRequest)
		} else if okthread {
			id = threadIDQ[0]
		} else if okcomment {
			id = commentIDQ[0]
		}
		uuid := uuid.New().String()

		writeComment(Comment{username, content, time.Now(), 0, []string{}, uuid})
		addComment(id, okthread, uuid)
		http.Redirect(w, req, "/forum/" + id, http.StatusSeeOther)
		return
	}
	tpls.ExecuteTemplate(w, "createcomment.gohtml", nil)
}

func addComment(rootid string, rootisthread bool, commentid string) {
	if rootisthread {
		updateThread(bson.D{{Key: "ID", Value: rootid}}, bson.D{{Key: "$push", Value: bson.D{{Key: "Replies", Value: commentid}}}})
	} else {
		updateComment(bson.D{{Key: "ID", Value: rootid}}, bson.D{{Key: "$push", Value: bson.D{{Key: "Replies", Value: commentid}}}})
	}
}