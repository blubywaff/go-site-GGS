package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var playersdb *mongo.Collection

type Player struct {
	IsTraining bool   `bson:"IsTraining"`
	Username   string `bson:"Username"`
	Ships      []Ship `bson:"Ships"`
	Base       Base   `bson:"Base"`
}

type Base struct {
	Strength int `bson:"Strength"`
	Power    int `bson:"Power"`
	Water    int `bson:"Water"`
	Metal    int `bson:"Metal"`
	Fuel     int `bson:"Fuel"`
}

type Ship struct {
	Type     string `bson:"Type"`
	IsMain   bool   `bson:"IsMain"`
	Strength int    `bson:"Strength"`
	Defense  int    `bson:"Defense"`
	Crew     int    `bson:"Crew"`
}

func getPlayer(username string) Player {
	return readPlayer(bson.D{{"Username", username}})
}

func getShips(username string) []Ship {
	return getPlayer(username).Ships
}

func getBase(username string) Base {
	return getPlayer(username).Base
}

func aggregatePlayersdb(pipeline mongo.Pipeline) []bson.M {
	cursor, err := playersdb.Aggregate(context.Background(), pipeline)
	if !check(err) {
		return []bson.M{}
	}
	var m []bson.M
	err = cursor.All(context.Background(), &m)
	check(err)
	return m
}

func readPlayer(filter bson.D) Player {
	player := Player{}
	err := playersdb.FindOne(context.Background(), filter).Decode(&player)
	check(err)
	return player
}

func writePlayer(player Player) {
	_, err := playersdb.InsertOne(context.Background(), player)
	check(err)
}

func removePlayer(filter bson.D) {
	res := playersdb.FindOneAndDelete(context.Background(), filter)
	check(res.Err())
}

func updatePlayer(filter bson.D, update bson.D) {
	_, err := playersdb.UpdateOne(context.Background(), filter, update)
	check(err)
}

func containsPlayer(filter bson.D) bool {
	player := Player{}
	err := playersdb.FindOne(context.Background(), filter).Decode(&player)
	check(err)
	return err == nil
}
