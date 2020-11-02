package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var playersdb *mongo.Collection

type Player struct {
	HasTrained bool   `bson:"HasTrained"`
	Username   string `bson:"Username"`
	Ships      []Ship `bson:"Ships"`
	Base       Base   `bson:"Base"`
}

type Base struct {
	Owner   string   `bson:"Owner"`
	Power   int      `bson:"Power"`
	Water   int      `bson:"Water"`
	Metal   int      `bson:"Metal"`
	Fuel    int      `bson:"Fuel"`
	Planets []int    `bson:"Planets"`
	Turrets []Turret `bson:"Turrets"`
}

type Ship struct {
	ID    string `bson:"ID"`
	Level int    `bson:"Level"`
	Crew  int    `bson:"Crew"`
}

type Turret struct {
	ID       string   `bson:"ID"`
	Level    int      `bson:"Level"`
	Position Position `bson:"Position"`
}

type Position struct {
	X int `bson:"X"`
	Y int `bson:"Y"`
}

func (b Base) hasTurretByPosition(pos Position) bool {
	return b.getTurretByPosition(pos) != Turret{}
}

func (b Base) hasTurretByID(id string) bool {
	return b.getTurretByID(id) != Turret{}
}

func (b Base) getTurretByPosition(pos Position) Turret {
	for _, turret := range b.Turrets {
		if turret.Position == pos {
			return turret
		}
	}
	return Turret{}
}

func (b Base) getTurretByID(id string) Turret {
	for _, turret := range b.Turrets {
		if turret.ID == id {
			return turret
		}
	}
	return Turret{}
}

func (p Player) hasShipByID(id string) bool {
	return p.getShipByID(id) != Ship{}
}

func (p Player) getShipByID(id string) Ship {
	for _, ship := range p.Ships {
		if ship.ID == id {
			return ship
		}
	}
	return Ship{}
}

func getPlayer(username string) Player {
	return readPlayer(bson.D{{"Username", username}})
}

func aggregatePlayersdb(pipeline mongo.Pipeline) []bson.M {
	cursor, err := playersdb.Aggregate(ctx, pipeline)
	if !check(err) {
		return []bson.M{}
	}
	var m []bson.M
	err = cursor.All(ctx, &m)
	check(err)
	return m
}

func readPlayer(filter bson.D) Player {
	player := Player{}
	err := playersdb.FindOne(ctx, filter).Decode(&player)
	check(err)
	return player
}

func writePlayer(player Player) {
	_, err := playersdb.InsertOne(ctx, player)
	check(err)
}

func removePlayer(filter bson.D) {
	res := playersdb.FindOneAndDelete(ctx, filter)
	check(res.Err())
}

func updatePlayer(filter bson.D, update bson.D) {
	_, err := playersdb.UpdateOne(ctx, filter, update)
	check(err)
}

func containsPlayer(filter bson.D) bool {
	player := Player{}
	err := playersdb.FindOne(ctx, filter).Decode(&player)
	check(err)
	return err == nil
}

func getShips(username string) []Ship {
	//TODO fix this
	return getPlayer(username).Ships
}

func getShip(username string, shipID string) Ship {
	cursor, err := playersdb.Aggregate(ctx, mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"Username", username},
			}},
		},
		bson.D{
			{"$project", bson.D{
				{"Ships", true},
			}},
		},
		bson.D{
			{"$unwind", bson.D{
				{"path", "$Ships"},
			}},
		},
		bson.D{
			{"$match", bson.D{
				{"ID", shipID},
			}},
		},
	})
	if !check(err) {
		return Ship{}
	}
	var m []bson.M
	err = cursor.All(ctx, &m)
	if !check(err) {
		return Ship{}
	}
	if len(m) < 1 {
		fmt.Println("No ship")
		return Ship{}
	}
	mShip := m[0]["Ships"].(bson.M)
	bb, err := bson.Marshal(mShip)
	check(err)
	var ship Ship
	err = bson.Unmarshal(bb, &ship)
	return ship
}

func updateShip(username string, shipID string, update bson.D) {
	updatePlayer(bson.D{{"Username", username}, {"Ships.ID", shipID}}, update)
}

func writeShip(username string, ship Ship) {
	updatePlayer(bson.D{{"Username", username}}, bson.D{{"$push", bson.D{{"Ships", ship}}}})
}

func removeShip(username string, shipID string) {
	updatePlayer(bson.D{{"Username", username}}, bson.D{{"$pull", bson.D{{"Ships", bson.D{{"ID", shipID}}}}}})
}

func getBase(username string) Base {
	return getPlayer(username).Base
}

func updateBase(username string, update bson.D) {
	updatePlayer(bson.D{{"Username", username}}, update)
}

func writeBase(username, string, base Base) {
	updatePlayer(bson.D{{"Username", username}}, bson.D{{"$set", bson.D{{"Base", base}}}})
}

func getTurrets(username string) []Turret {
	//TODO fix this
	return getPlayer(username).Base.Turrets
}

func getTurret(username string, turretID string) Turret {
	cursor, err := playersdb.Aggregate(ctx, mongo.Pipeline{
		bson.D{
			{"$match", bson.D{
				{"Username", username},
			}},
		},
		bson.D{
			{"$project", bson.D{
				{"Turrets", "$Base.Turrets"},
			}},
		},
		bson.D{
			{"$project", bson.D{
				{"Turrets", true},
			}},
		},
		bson.D{
			{"$unwind", bson.D{
				{"path", "$Turrets"},
			}},
		},
		bson.D{
			{"$match", bson.D{
				{"ID", turretID},
			}},
		},
	})
	if !check(err) {
		return Turret{}
	}
	var m []bson.M
	err = cursor.All(ctx, &m)
	if !check(err) {
		return Turret{}
	}
	if len(m) < 1 {
		fmt.Println("No turret")
		return Turret{}
	}
	mTurret := m[0]["Turrets"].(bson.M)
	bb, err := bson.Marshal(mTurret)
	check(err)
	var turret Turret
	err = bson.Unmarshal(bb, &turret)
	return turret
}

func updateTurret(username string, turretID string, update bson.D) {
	updatePlayer(bson.D{{"Username", username}, {"Base.Turrets.ID", turretID}}, update)
}

func writeTurret(username string, turret Turret) {
	updatePlayer(bson.D{{"Username", username}}, bson.D{{"$push", bson.D{{"Base.Turrets", turret}}}})
}

func removeTurret(username string, turretID string) {
	updatePlayer(bson.D{{"Username", username}}, bson.D{{"$pull", bson.D{{"Base.Turrets", bson.D{{"ID", turretID}}}}}})
}
