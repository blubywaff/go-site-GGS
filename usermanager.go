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

type profile struct{
	Username string `bson:"Username"`
	Picture []byte `bson:"Picture"`
}

var usersdb *mongo.Collection
var profilePicturesdb *mongo.Collection

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

func readProfilePicture(username string) []byte {
	profile := profile{}
	err := profilePicturesdb.FindOne(context.Background(), bson.D{{Key: "Username", Value: username}}).Decode(&profile)
	if !check(err) {
		return []byte{}
	}
	return profile.Picture
}

func writeProfilePicture(username string, file []byte) {
	_, err := profilePicturesdb.InsertOne(context.Background(), profile{username, file})
	check(err)
}

func removeProfilePicture(username string) {
	res := profilePicturesdb.FindOneAndDelete(context.Background(), bson.D{{Key: "Username", Value: username}})
	check(res.Err())
}

func updateProfilePicture(username string, file []byte) {
	_, err := profilePicturesdb.UpdateOne(context.Background(), bson.D{{Key: "Username", Value: username}}, bson.D{{Key: "$set", Value: bson.D{{Key: "Picture", Value: file}}}})
	check(err)
}

func containsProfilePicture(username string) bool {
	profile := profile{}
	err := profilePicturesdb.FindOne(context.Background(), bson.D{{Key: "Username", Value: username}}).Decode(&profile)
	check(err)
	return err == nil
}