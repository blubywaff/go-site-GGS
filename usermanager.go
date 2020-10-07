package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type user struct {
	Username  string `bson:"Username"`
	Email     string `bson:"Email"`
	Password  []byte `bson:"Password"`
	Firstname string `bson:"Firstname"`
	Lastname  string `bson:"Lastname"`
}

type profile struct {
	Username string `bson:"Username"`
	Picture  []byte `bson:"Picture"`
}

var usersdb *mongo.Collection
var profilePicturesdb *mongo.Collection

func writeUser(user user) {
	_, err := usersdb.InsertOne(ctx, user)
	check(err)
}

func removeUser(filter bson.D) {
	res := usersdb.FindOneAndDelete(ctx, filter)
	check(res.Err())
}

func updateUser(filter bson.D, user user) {
	_, err := usersdb.UpdateOne(ctx, filter, user)
	check(err)
}

func readUser(filter bson.D) user {
	user := user{}
	err := usersdb.FindOne(ctx, filter).Decode(&user)
	check(err)
	return user
}

func containsUser(filter bson.D) bool {
	user := user{}
	err := usersdb.FindOne(ctx, filter).Decode(&user)
	//fmt.Println("contains", err)
	check(err)
	return err == nil
}

func readProfilePicture(username string) []byte {
	profile := profile{}
	err := profilePicturesdb.FindOne(ctx, bson.D{{Key: "Username", Value: username}}).Decode(&profile)
	if !check(err) {
		return []byte{}
	}
	return profile.Picture
}

func writeProfilePicture(username string, file []byte) {
	_, err := profilePicturesdb.InsertOne(ctx, profile{username, file})
	check(err)
}

func removeProfilePicture(username string) {
	res := profilePicturesdb.FindOneAndDelete(ctx, bson.D{{Key: "Username", Value: username}})
	check(res.Err())
}

func updateProfilePicture(username string, file []byte) {
	_, err := profilePicturesdb.UpdateOne(ctx, bson.D{{Key: "Username", Value: username}}, bson.D{{Key: "$set", Value: bson.D{{Key: "Picture", Value: file}}}})
	check(err)
}

func containsProfilePicture(username string) bool {
	profile := profile{}
	err := profilePicturesdb.FindOne(ctx, bson.D{{Key: "Username", Value: username}}).Decode(&profile)
	check(err)
	return err == nil
}
