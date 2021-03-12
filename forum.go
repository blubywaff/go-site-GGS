package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strconv"
	"time"
)

func forum(w http.ResponseWriter, req *http.Request) {
	tpls.ExecuteTemplate(w, "forumhome.gohtml", getForumData())
}

func forumThread(w http.ResponseWriter, req *http.Request) {
	threadIDQ, ok := req.URL.Query()["thread"]
	if !ok {
		http.Error(w, "No Thread ID", http.StatusBadRequest)
	}
	threadID := threadIDQ[0]

	/*if len(threadID) < 36 {
		http.Error(w, "Invalid Thread ID", http.StatusBadRequest)
		return
	}*/

	if !containsThread(bson.D{{Key: "ID", Value: threadID}}) {
		http.Redirect(w, req, "/forum/", http.StatusSeeOther)
	}

	tpls.ExecuteTemplate(w, "thread.gohtml", getThread(threadID).getFull())
}

func forumData(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		fmt.Fprint(w, "You can't send stuff here, silly!")
		return
	}
	timeWindow := req.URL.Query()["time"]
	num := req.URL.Query()["num"]
	page := req.URL.Query()["page"]
	if len(timeWindow) != 1 || len(num) != 1 || len(page) != 1 {
		http.Error(w, "Invalid Parameters", http.StatusBadRequest)
		return
	}

	single := func(i int, e error) int { check(e); return i }
	s := func(i []byte, e error) []byte { check(e); return i }

	tw := single(strconv.Atoi(timeWindow[0]))
	n := single(strconv.Atoi(num[0]))
	p := single(strconv.Atoi(page[0]))

	if tw < 0 || tw > 4 {
		http.Error(w, "Invalid Time Parameter", http.StatusBadRequest)
		return
	}
	if n < 0 {
		http.Error(w, "Invalid Number Parameter", http.StatusBadRequest)
		return
	} else if n > 50 {
		http.Error(w, "Do not request more than 50 at a time", http.StatusBadRequest)
		return
	}
	if p < 0 {
		http.Error(w, "Where are you trying to go?", http.StatusBadRequest)
		return
	} else if p >= 10000 {
		http.Error(w, "Those records are on the high shelf!", http.StatusTeapot)
		return
	}

	fmt.Fprint(w, string(s(json.Marshal(getForums(tw, n, p)))))
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
		http.Redirect(w, req, "/thread/?thread="+id.String(), http.StatusSeeOther)
		return
	}
	tpls.ExecuteTemplate(w, "createthread.gohtml", nil)
}

func vote(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(w, req) {
		http.Error(w, "Not Logged In!", http.StatusUnauthorized)
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
			http.Redirect(w, req, "/thread/?thread="+id, http.StatusSeeOther)
			return
		}
		updateVote(username, id, okthread, vote)
		voteOn(id, okthread, 2*vote)
	} else {
		writeVote(username, id, okthread, vote)
		voteOn(id, okthread, vote)
	}
	fmt.Fprint(w, "Vote Successful!")
}

func voteOn(id string, isThread bool, vote int) {
	if isThread {
		updateThread(bson.D{{Key: "ID", Value: id}}, bson.D{{Key: "$inc", Value: bson.D{{Key: "Score", Value: vote}}}})
	} else {
		updateComment(bson.D{{Key: "ID", Value: id}}, bson.D{{Key: "$inc", Value: bson.D{{Key: "Score", Value: vote}}}})
	}
}

func createComment(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Post comment via thread page", http.StatusBadRequest)
		return
	}

	if !alreadyLoggedIn(w, req) {
		http.Error(w, "Must have session", http.StatusUnauthorized)
		return
	}

	content := req.FormValue("content")
	username := getUser(w, req).Username
	threadIDQ, okthread := req.URL.Query()["thread"]
	commentIDQ, okcomment := req.URL.Query()["comment"]
	var id string
	if okthread && okcomment {
		http.Error(w, "Both Thread and Comment Provided", http.StatusBadRequest)
		return
	} else if !okthread && !okcomment {
		http.Error(w, "No Resource Provided", http.StatusBadRequest)
		return
	} else if okthread {
		id = threadIDQ[0]
	} else if okcomment {
		id = commentIDQ[0]
	}
	commentID := uuid.New().String()

	fmt.Println(content)

	writeComment(Comment{username, content, time.Now(), 0, []string{}, commentID})
	addComment(id, okthread, commentID)
	fmt.Fprint(w, "Comment Added!")
	return
}

func readComments(w http.ResponseWriter, req *http.Request) {
	threadIDQ, okthread := req.URL.Query()["thread"]
	commentIDQ, okcomment := req.URL.Query()["comment"]
	var id string
	var item interface{}
	if okthread && okcomment {
		http.Error(w, "Both Thread and Comment Provided", http.StatusBadRequest)
		return
	} else if !okthread && !okcomment {
		http.Error(w, "No Resource Provided", http.StatusBadRequest)
		return
	} else if okthread {
		id = threadIDQ[0]
		item = getThread(id)
	} else if okcomment {
		id = commentIDQ[0]
		item = getThread(id)
	}
	//comment := getComment(id)
	var comments []Comment
	itemM, ok := item.(map[string]interface{})
	fmt.Println(itemM)
	fmt.Println(item)
	fmt.Println(ok)
	for _, k := range itemM["Replies"].([]string) {
		comments = append(comments, getComment(k))
	}
	fmt.Println(comments)
	fmt.Println(item)
	jsonC, _ := json.Marshal(comments)
	fmt.Fprint(w, string(jsonC))
}

func fullComment(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		http.Error(w, "You can't send stuff here, silly!", http.StatusBadRequest)
		return
	}
	id := req.URL.Query()["ID"]
	if len(id) < 1 || len(id) > 1 {
		http.Error(w, "Must have only one ID request", http.StatusBadRequest)
	}
	comment := getComment(id[0])
	jsonP, _ := json.Marshal(comment.getFull())
	fmt.Fprint(w, jsonP)

}

func addComment(rootid string, rootisthread bool, commentid string) {
	if rootisthread {
		updateThread(bson.D{{Key: "ID", Value: rootid}}, bson.D{{Key: "$push", Value: bson.D{{Key: "Replies", Value: commentid}}}})
	} else {
		updateComment(bson.D{{Key: "ID", Value: rootid}}, bson.D{{Key: "$push", Value: bson.D{{Key: "Replies", Value: commentid}}}})
	}
}
