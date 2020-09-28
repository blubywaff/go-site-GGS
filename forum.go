package main

import (
	"net/http"
	"github.com/google/uuid"
	"time"
	"strconv"
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

	tpls.ExecuteTemplate(w, "thread.gohtml", getThread(threadID))
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
	if !check(err) {
		http.Error(w, "Invalid Vote", http.StatusBadRequest)
	}

	if containsVote(username, id, okthread) {
		if getVote(username, id, okthread).Vote == vote {
			removeVote(username, id, okthread)
			return
		}
		updateVote(username, id, okthread, vote)
	} else {
		writeVote(username, id, okthread, vote)
	}
}