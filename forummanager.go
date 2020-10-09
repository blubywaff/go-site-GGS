package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var threadsdb *mongo.Collection
var votesdb *mongo.Collection

type Thread struct {
	Poster   string    `bson:"Poster"`
	Title    string    `bson:"Title"`
	PostTime time.Time `bson:"PostTime"`
	ID       string    `bson:"ID"`
	Body     string    `bson:"Body"`
	Score    int       `bson:"Score"`
	Replies  []Comment `bson:"Replies"`
}

type Comment struct {
	Poster   string    `bson:"Poster"`
	Content  string    `bson:"Content"`
	PostTime time.Time `bson:"PostTime"`
	Score    int       `bson:"Score"`
	Replies  []Comment `bson:"Replies"`
	ID       string    `bson:"ID"`
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

func readComment(id, rootId string) Comment {
	//TODO Aggregation
	cursor, err := threadsdb.Aggregate(ctx, mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"ID", rootId},
			}},
		},
		bson.D{
			{"$redact", bson.D{
				{"$cond", bson.D{
					{"if", bson.D{
						{"ID", id},
					}},
					{"then", "$$KEEP"},
					{"else", "DESCEND"},
				}},
			}},
		},
	})
	if !check(err) {
		return Comment{}
	}
	var m []bson.M
	err = cursor.All(ctx, &m)
	if !check(err) {
		return Comment{}
	}
	fmt.Println(m)
	var result Comment
	bb, _ := bson.Marshal(m[0])
	_ = bson.Unmarshal(bb, &result)
	return result
}

func writeComment(comment Comment, rootId, trueRoot string) {
	if rootId == trueRoot {

	}
	_, err := threadsdb.Aggregate(ctx, mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"ID", trueRoot},
			}},
		},
		bson.D{
			{"$redact", bson.D{
				{"$cond", bson.D{
					{"if", bson.D{
						{"ID", rootId},
					}},
					{"then", "$$KEEP"},
					{"else", "DESCEND"},
				}},
			}},
		},
		bson.D{
			{"$push", bson.D{
				{"Replies", comment},
			}},
		},
	})
	check(err)
}

func removeComment(id, rootId, trueRoot string) {
	_, err := threadsdb.Aggregate(ctx, mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"ID", trueRoot},
			}},
		},
		bson.D{
			{"$redact", bson.D{
				{"$cond", bson.D{
					{"if", bson.D{
						{"ID", rootId},
					}},
					{"then", "$$KEEP"},
					{"else", "DESCEND"},
				}},
			}},
		},
		bson.D{
			{"$pull", bson.D{
				{"Replies", bson.D{
					{"ID", id},
				}},
			}},
		},
	})
	check(err)
}

func updateComment(id, rootId, trueRoot string, update bson.D) {
	_, err := threadsdb.Aggregate(ctx, mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"ID", trueRoot},
			}},
		},
		bson.D{
			{"$redact", bson.D{
				{"$cond", bson.D{
					{"if", bson.D{
						{"ID", rootId},
					}},
					{"then", "$$KEEP"},
					{"else", "DESCEND"},
				}},
			}},
		},
		update,
	})
	check(err)
}

func containsComment(id, rootID string) bool {
	return readComment(id, rootID).ID == id
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
