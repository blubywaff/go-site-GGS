package main

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"context"
)

var dbTimeFormat = "01/02/2006 15:04:05.000"

type session struct{
	SessionID string `bson:"SessionID"`
	Username string `bson:"Username"`
	LastActivity string `bson:"LastActivity"`
}

var sessionsdb *mongo.Collection

func writeSession(session session) {
	_, err := sessionsdb.InsertOne(context.Background(), session)
	check(err)
}

func removeSession(filter bson.D) {
	res := sessionsdb.FindOneAndDelete(context.Background(), filter)
	check(res.Err())
}

func updateSession(filter bson.D, update bson.D) {
	_, err := sessionsdb.UpdateOne(context.Background(), filter, update)
	check(err)
}

func readSession(filter bson.D) session {
	session := session{}
	err := sessionsdb.FindOne(context.Background(), filter).Decode(&session)
	check(err)
	return session
}

func containsSession(filter bson.D) bool {
	session := session{}
	err := sessionsdb.FindOne(context.Background(), filter).Decode(&session)
	check(err)
	return err == nil
}