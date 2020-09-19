package main

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"context"
)

type user struct{
	Username string `bson:"Username"`
	Email string `bson:"Email"`
	Password []byte `bson:"Password"`
	Firstname string `bson:"Firstname"`
	Lastname string `bson:"Lastname"`
}

const PRIVATE = 0
const RESTRICTED = 1
const OPEN = 2

var usersdb *mongo.Collection

func writeUser(user user) {
	_, err := usersdb.InsertOne(context.Background(), user)
	check(err)
}

func removeUser(filter bson.D) {
	res := usersdb.FindOneAndDelete(context.Background(), filter)
	check(res.Err())
}

func updateUser(filter bson.D, user user) {
	_, err := usersdb.UpdateOne(context.Background(), filter, user)
	check(err)
}

func readUser(filter bson.D) user {
	user := user{}
	err := usersdb.FindOne(context.Background(), filter).Decode(&user)
	check(err)
	return user
}

func containsUser(filter bson.D) bool {
	user := user{}
	err := usersdb.FindOne(context.Background(), filter).Decode(&user)
	//fmt.Println("contains", err)
	check(err)
	return err == nil
}