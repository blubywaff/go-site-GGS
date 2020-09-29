package main

import (
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"context"
)

var threadsdb *mongo.Collection
var commentsdb *mongo.Collection
var votesdb *mongo.Collection

type Thread struct {
	Poster string `bson:"Poster"`
	Title string `bson:"Title"`
	PostTime time.Time `bson:"PostTime"`
	ID string `bson:"ID"`
	Body string `bson:"Body"`
	Score int `bson:"Score"`
	Replies []string `bson:"Replies"`
}

/*type ThreadInfo struct {
	Poster string `bson:"Poster"`
	Title string `bson:"Title"`
	PostTime time.Time `bson:"PostTime"`
	ID string `bson:"ID"`
	Score int `bson:"Score"`
}*/

type Comment struct {
	Poster string `bson:"Poster"`
	Content string `bson:"Content"`
	PostTime time.Time `bson:"PostTime"`
	Score int `bson:"Score"`
	Replies []string `bson:"Replies"`
	ID string `bson:"ID"`
}

type FullThread struct {
	Poster string 
	Title string 
	PostTime time.Time 
	ID string 
	Body string 
	Score int 
	Replies []FullComment
}

type FullComment struct {
	Poster string 
	Content string 
	PostTime time.Time 
	Score int 
	Replies []FullComment
	ID string
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

type FormData struct {
	Top []Thread `bson:"Top"`
}

type Vote struct {
	Post string `bson:"Post"`
	IsThread bool `bson:"IsThread"`
	Vote int `bson:"Vote"`
}

type Votes struct {
	Username string `bson:"Username"`
	Votes []Vote `bson:"Votes"`
}

/*
	cursor, err :=threadsdb.Aggregate(context.Background(), mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"PostTime", bson.D{
					{"$gte", time.Now().Add(time.Hour * -24)},
				}},
			}},
			
		},
		bson.D{
			{"$group", bson.D{
				{"_id", "ID"},
				{"score", bson.D{
					{"$max", "$Score"},
				}},
			}},
		},
	})
	*/

func getForumData() FormData {
	exclude := bson.A{}
	formData := FormData{}
	for i := 0; i < 10; i++ {
		cursor, err := threadsdb.Aggregate(context.Background(), mongo.Pipeline{
			bson.D{
				{Key: "$match", Value: bson.D{
					{Key: "PostTime", Value: bson.D{
						{Key: "$gte", Value: time.Now().Add(time.Hour * -24)},
					}},
					{Key: "ID", Value: bson.D{
						{Key: "$nin", Value: exclude},
					}},
				}},
			},
			bson.D{
				{Key: "$group", Value: bson.D{
					{Key: "_id", Value: "$ID"},
					{Key: "score", Value: bson.D{
						{Key: "$max", Value: "$Score"},
					}},
				}},
			},
		})
		if !check(err) {
			break
		}
		var m []bson.M
		err = cursor.All(context.Background(), &m)
		check(err)
		if len(m) < 1 {
			break
		}
		exclude = append(exclude, m[0]["_id"].(string))
		formData.Top = append(formData.Top, getThread(m[0]["_id"].(string)))
	}
	return formData
}

func getThread(id string) Thread {
	thread := Thread{}
	err := threadsdb.FindOne(context.Background(), bson.D{{Key: "ID", Value: id}}).Decode(&thread)
	check(err)
	return thread
}

func readThread(filter bson.D) Thread {
	thread := Thread{}
	err := threadsdb.FindOne(context.Background(), filter).Decode(&thread)
	check(err)
	return thread
}

func writeThread(thread Thread) {
	_, err := threadsdb.InsertOne(context.Background(), thread)
	check(err)
}

func removeThread(filter bson.D) {
	res := threadsdb.FindOneAndDelete(context.Background(), filter)
	check(res.Err())
}

func updateThread(filter bson.D, update bson.D) {
	_, err := threadsdb.UpdateOne(context.Background(), filter, update)
	check(err)
}

func containsThread(filter bson.D) bool {
	thread := Thread{}
	err := threadsdb.FindOne(context.Background(), filter).Decode(&thread)
	check(err)
	return err == nil
}



func getComment(id string) Comment {
	comment := Comment{}
	err := commentsdb.FindOne(context.Background(), bson.D{{Key: "ID", Value: id}}).Decode(&comment)
	check(err)
	return comment
}

func readComment(filter bson.D) Comment {
	comment := Comment{}
	err := commentsdb.FindOne(context.Background(), filter).Decode(&comment)
	check(err)
	return comment
}

func writeComment(comment Comment) {
	_, err := commentsdb.InsertOne(context.Background(), comment)
	check(err)
}

func removeComment(filter bson.D) {
	res := commentsdb.FindOneAndDelete(context.Background(), filter)
	check(res.Err())
}

func updateComment(filter bson.D, update bson.D) {
	_, err := commentsdb.UpdateOne(context.Background(), filter, update)
	check(err)
}

func containsComment(filter bson.D) bool {
	comment := Comment{}
	err := commentsdb.FindOne(context.Background(), filter).Decode(&comment)
	check(err)
	return err == nil
}



func getVotes(username string) Votes {
	votes := Votes{}
	err := votesdb.FindOne(context.Background(), bson.D{{Key: "Username", Value: username}}).Decode(&votes)
	check(err)
	return votes
}

func readVotes(filter bson.D) Votes {
	votes := Votes{}
	err := votesdb.FindOne(context.Background(), filter).Decode(&votes)
	check(err)
	return votes
}

func writeVotes(votes Votes) {
	_, err := votesdb.InsertOne(context.Background(), votes)
	check(err)
}

func removeVotes(filter bson.D) {
	res := votesdb.FindOneAndDelete(context.Background(), filter)
	check(res.Err())
}

func updateVotes(filter bson.D, update bson.D) {
	_, err := votesdb.UpdateOne(context.Background(), filter, update)
	check(err)
}

func containsVotes(filter bson.D) bool {
	votes := Votes{}
	err := votesdb.FindOne(context.Background(), filter).Decode(&votes)
	check(err)
	return err == nil
}



func containsVote(username string, post string, isThread bool) bool {
	//return containsVotes(bson.D{{Key: "Username", Value: username}, {Key: "Votes.$.Post", Value: post}, {Key: "Votes.$.IsThread", Value: isThread}})
	return getVote(username, post, isThread) != Vote{}
}

func getVote(username string, post string, isThread bool) Vote {
	votes := getVotes(username).Votes
	for i := 0; i < len(votes); i++ {
		vote := votes[i]
		if vote.IsThread == isThread && vote.Post == post {
			return vote
		}
	}
	return Vote{}
}

func writeVote(username string, post string, isThread bool, vote int) {
	updateVotes(bson.D{{Key: "Username", Value: username}}, bson.D{{Key: "$push", Value: bson.D{{Key: "Votes", Value: Vote{post, isThread, vote,}}}}})
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