package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var threadsdb *mongo.Collection
var commentsdb *mongo.Collection
var votesdb *mongo.Collection

type Thread struct {
	Poster   string    `bson:"Poster"`
	Title    string    `bson:"Title"`
	PostTime time.Time `bson:"PostTime"`
	ID       string    `bson:"ID"`
	Body     string    `bson:"Body"`
	Score    int       `bson:"Score"`
	Replies  []string  `bson:"Replies"`
}

/*type ThreadInfo struct {
	Poster string `bson:"Poster"`
	Title string `bson:"Title"`
	PostTime time.Time `bson:"PostTime"`
	ID string `bson:"ID"`
	Score int `bson:"Score"`
}*/

type Comment struct {
	Poster   string    `bson:"Poster"`
	Content  string    `bson:"Content"`
	PostTime time.Time `bson:"PostTime"`
	Score    int       `bson:"Score"`
	Replies  []string  `bson:"Replies"`
	ID       string    `bson:"ID"`
}

type FullThread struct {
	Poster   string
	Title    string
	PostTime time.Time
	ID       string
	Body     string
	Score    int
	Replies  []FullComment
}

type FullComment struct {
	Poster   string
	Content  string
	PostTime time.Time
	Score    int
	Replies  []FullComment
	ID       string
}

func (thread Thread) getFull() FullThread {
	full := FullThread{
		thread.Poster,
		thread.Title,
		thread.PostTime,
		thread.ID,
		thread.Body,
		thread.Score,
		[]FullComment{},
	}
	fulls := []FullComment{}
	for _, c := range thread.Replies {
		fulls = append(fulls, getComment(c).getFull())
	}
	full.Replies = fulls
	return full
}

func (comment Comment) getFull() FullComment {
	full := FullComment{
		comment.Poster,
		comment.Content,
		comment.PostTime,
		comment.Score,
		[]FullComment{},
		comment.ID,
	}
	fulls := []FullComment{}
	for _, c := range comment.Replies {
		fulls = append(fulls, getComment(c).getFull())
	}
	full.Replies = fulls
	return full
}

type ForumData struct {
	Top []Thread `bson:"Top"`
}

type Vote struct {
	Post     string `bson:"Post"`
	IsThread bool   `bson:"IsThread"`
	Vote     int    `bson:"Vote"`
}

type Votes struct {
	Username string `bson:"Username"`
	Votes    []Vote `bson:"Votes"`
}

func getForumData() ForumData {
	if SUPER_DEBUG_MODE_OVERRIDE {
		fd := ForumData{}
		fd.Top = append(fd.Top, Thread{
			Poster:   "DEBUG_ACCOUNT",
			Title:    "DEBUG !",
			PostTime: time.Time{},
			ID:       "DEBUG_IMPOSSIBLE_ID_THIS_SHOULD_STILL_WORK",
			Body:     "hello this is a test debug message, it should be impossible to see this unless the server is in debug mode",
			Score:    -21,
			Replies:  nil,
		}, Thread{
			Poster:   "DEBUG_ACCOUNT",
			Title:    "DEBUG POST 2",
			PostTime: time.Time{},
			ID:       "DEBUG_IMPOSSIBLE_NUMBER_2",
			Body:     "hi hi hi hi hi hi this is a useleess post hi hi hi hi hi hi hi hi hi hi hi hi hi hi hi hi hi hi hi hi hi",
			Score:    7777777777777777777,
			Replies:  nil,
		})
		return fd
	}

	cursor, err := threadsdb.Aggregate(ctx, mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"PostTime", bson.D{
					{"$gte", time.Now().Add(time.Hour * -24)},
				}},
			}},
		},
		bson.D{
			{"$sort", bson.D{
				{"Score", -1},
			}},
		},
		bson.D{
			{"$limit", 10},
		},
	})
	if !check(err) {
		return ForumData{}
	}
	var m []bson.M
	err = cursor.All(ctx, &m)
	if !check(err) {
		return ForumData{}
	}
	var result ForumData
	for _, bm := range m {
		var t Thread
		bb, _ := bson.Marshal(bm)
		_ = bson.Unmarshal(bb, &t)
		result.Top = append(result.Top, t)
	}
	return result
}

func getForums(timeWindow int, num int, page int) ForumData {
	fmt.Println(timeWindow, num, page)
	cursor, err := threadsdb.Aggregate(ctx, mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"PostTime", bson.D{
					{"$gte", func() time.Time {
						if timeWindow == 0 {
							return time.Now().Add(time.Hour * -24)
						} else if timeWindow == 1 {
							return time.Now().Add(time.Hour * -24 * 7)
						} else if timeWindow == 2 {
							return time.Now().Add(time.Hour * -24 * 30)
						} else if timeWindow == 3 {
							return time.Now().Add(time.Hour * -24 * 365)
						} else {
							return time.Time{}
						}
					}()},
				}},
			}},
		},
		bson.D{
			{"$sort", bson.D{
				{"Score", -1},
				{"PostTime", -1},
				{"_id", -1},
			}},
		},
		bson.D{
			{"$limit", num * (page + 1)},
		},
		bson.D{
			{"$sort", bson.D{
				{"Score", 1},
				{"PostTime", 1},
				{"_id", 1},
			}},
		},
		bson.D{
			{"$limit", num},
		},
		bson.D{
			{"$sort", bson.D{
				{"Score", -1},
				{"PostTime", -1},
				{"_id", -1},
			}},
		},
	})
	if !check(err) {
		return ForumData{}
	}
	var m []bson.M
	err = cursor.All(ctx, &m)
	if !check(err) {
		return ForumData{}
	}
	var result ForumData
	for _, bm := range m {
		var t Thread
		bb, _ := bson.Marshal(bm)
		_ = bson.Unmarshal(bb, &t)
		result.Top = append(result.Top, t)
	}
	return result
}

func getThread(id string) Thread {
	thread := Thread{}
	err := threadsdb.FindOne(ctx, bson.D{{Key: "ID", Value: id}}).Decode(&thread)
	check(err)
	return thread
}

func readThread(filter bson.D) Thread {
	thread := Thread{}
	err := threadsdb.FindOne(ctx, filter).Decode(&thread)
	check(err)
	return thread
}

func writeThread(thread Thread) {
	_, err := threadsdb.InsertOne(ctx, thread)
	check(err)
}

func removeThread(filter bson.D) {
	res := threadsdb.FindOneAndDelete(ctx, filter)
	check(res.Err())
}

func updateThread(filter bson.D, update bson.D) {
	_, err := threadsdb.UpdateOne(ctx, filter, update)
	check(err)
}

func containsThread(filter bson.D) bool {
	thread := Thread{}
	err := threadsdb.FindOne(ctx, filter).Decode(&thread)
	check(err)
	return err == nil
}

func getComment(id string) Comment {
	comment := Comment{}
	err := commentsdb.FindOne(ctx, bson.D{{Key: "ID", Value: id}}).Decode(&comment)
	check(err)
	return comment
}

func readComment(filter bson.D) Comment {
	comment := Comment{}
	err := commentsdb.FindOne(ctx, filter).Decode(&comment)
	check(err)
	return comment
}

func writeComment(comment Comment) {
	_, err := commentsdb.InsertOne(ctx, comment)
	check(err)
}

func removeComment(filter bson.D) {
	res := commentsdb.FindOneAndDelete(ctx, filter)
	check(res.Err())
}

func updateComment(filter bson.D, update bson.D) {
	_, err := commentsdb.UpdateOne(ctx, filter, update)
	check(err)
}

func containsComment(filter bson.D) bool {
	comment := Comment{}
	err := commentsdb.FindOne(ctx, filter).Decode(&comment)
	check(err)
	return err == nil
}

func getVotes(username string) Votes {
	votes := Votes{}
	err := votesdb.FindOne(ctx, bson.D{{Key: "Username", Value: username}}).Decode(&votes)
	check(err)
	return votes
}

func readVotes(filter bson.D) Votes {
	votes := Votes{}
	err := votesdb.FindOne(ctx, filter).Decode(&votes)
	check(err)
	return votes
}

func writeVotes(votes Votes) {
	_, err := votesdb.InsertOne(ctx, votes)
	check(err)
}

func removeVotes(filter bson.D) {
	res := votesdb.FindOneAndDelete(ctx, filter)
	check(res.Err())
}

func updateVotes(filter bson.D, update bson.D) {
	_, err := votesdb.UpdateOne(ctx, filter, update)
	check(err)
}

func containsVotes(filter bson.D) bool {
	votes := Votes{}
	err := votesdb.FindOne(ctx, filter).Decode(&votes)
	check(err)
	return err == nil
}

func containsVote(username string, post string, isThread bool) bool {
	//return containsVotes(bson.D{{Key: "Username", Value: username}, {Key: "Votes.$.Post", Value: post}, {Key: "Votes.$.IsThread", Value: isThread}})
	return getVote(username, post, isThread) != Vote{}
}

func getVote(username string, post string, isThread bool) Vote {
	cursor, err := votesdb.Aggregate(ctx, mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"Username", username},
				{"Votes.Post", post},
				{"Votes.IsThread", isThread},
			}},
		},
		bson.D{
			{Key: "$project", Value: bson.D{
				{"_id", false},
				{"Username", false},
			}},
		},
		bson.D{
			{"$unwind", bson.D{
				{"path", "$Votes"},
			}},
		},
		bson.D{
			{"$match", bson.D{
				{"Votes.IsThread", isThread},
				{"Votes.Post", post},
			}},
		},
	})
	if !check(err) {
		return Vote{}
	}
	var m []bson.M
	err = cursor.All(ctx, &m)
	fmt.Println(m)
	if !check(err) {
		return Vote{}
	}
	if len(m) < 1 {
		fmt.Println("no vote")
		return Vote{}
	}
	fmt.Println(m[0]["Votes"].(bson.M)["IsThread"])
	mVote := m[0]["Votes"].(bson.M)
	return Vote{mVote["Post"].(string), mVote["IsThread"].(bool), int(mVote["Vote"].(int32))}
}

func writeVote(username string, post string, isThread bool, vote int) {
	updateVotes(bson.D{{Key: "Username", Value: username}}, bson.D{{Key: "$push", Value: bson.D{{Key: "Votes", Value: Vote{post, isThread, vote}}}}})
}

func updateVote(username string, post string, isThread bool, vote int) {
	updateVotes(bson.D{{Key: "Username", Value: username}, {Key: "Votes.Post", Value: post}}, bson.D{{Key: "$set", Value: bson.D{{Key: "Votes.$.Vote", Value: vote}}}})
}

func removeVote(username string, post string, isThread bool) {
	updateVotes(bson.D{{Key: "Username", Value: username}}, bson.D{{Key: "$pull", Value: bson.D{{Key: "Votes", Value: bson.D{{Key: "IsThread", Value: isThread}, {Key: "Post", Value: post}}}}}})
}

func readVote(username string, post string, isThread bool) int {
	return getVote(username, post, isThread).Vote
}
