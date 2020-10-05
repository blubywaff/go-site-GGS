package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var playersdb *mongo.Collection

type Player struct {
	IsTraining bool
	Username   string
	Ships      []Ship
	Base       Base
}

type Base struct {
	ID       string
	Strength int
	Power    int
	Water    int
	Metal    int
	Fuel     int
}

type Ship struct {
	Type     string
	IsMain   bool
	Strength int
	Defense  int
	Crew     int
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
