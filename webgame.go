package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
)

func webgame(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/login/", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {
		webgameAjax(w, req)
		return
	}

	username := getUser(w, req).Username
	if !containsPlayer(bson.D{{"Username", username}}) {
		tpls.ExecuteTemplate(w, "newplayer.gohtml", nil)
		return
	}
	if !containsPlayer(bson.D{{"Username", username}, {"HasTrained", true}}) {
		tpls.ExecuteTemplate(w, "gamestart.gohtml", nil)
		return
	}
	tpls.ExecuteTemplate(w, "webgame.gohtml", nil)
	//playerT := readPlayer(bson.D{{"Username", username}, {"IsTraining", true}})
	//player := readPlayer(bson.D{{"Username", username}, {"IsTraining", false}})

}

func training(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/login/", http.StatusSeeOther)
		return
	}
	username := getUser(w, req).Username
	writePlayer(Player{false, username, []Ship{}, Base{}})
	tpls.ExecuteTemplate(w, "trainingground.gohtml", nil)
}

func gamestart(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/login/", http.StatusSeeOther)
		return
	}
	username := getUser(w, req).Username
	writePlayer(Player{false, username, []Ship{}, Base{}})
	tpls.ExecuteTemplate(w, "gamestart.gohtml", nil)
}

func gamedetails(w http.ResponseWriter, req *http.Request) {
	tpls.ExecuteTemplate(w, "gamedetails.gohtml", nil)
}

func webgameAjax(w http.ResponseWriter, req *http.Request) {
	bytes, err := ioutil.ReadAll(req.Body)
	check(err)
	username := getUser(w, req).Username
	requests := strings.Split(string(bytes), "|")
	act := requests[0]
	if act == "init" {
		if containsPlayer(bson.D{{"Username", username}}) {
			fmt.Fprint(w, "error-exists")
			return
		}
		writePlayer(Player{false, username, []Ship{}, Base{}})
	} else if act == "real" {
		if !containsPlayer(bson.D{{"Username", username}}) {
			fmt.Fprint(w, "error-exists")
			return
		}
		updatePlayer(bson.D{{"Username", username}}, bson.D{{"$set", bson.D{{"HasTrained", true}}}})
	} else if act == "turret" {
		// TODO fix to calculate and enforce cost of turret
		act = requests[1]
		base := getBase(username)
		if act == "add" {
			turret := Turret{}
			err = json.Unmarshal([]byte(requests[2]), &turret)
			if !check(err) {
				fmt.Fprint(w, "error-decode")
				return
			}
			if base.hasTurretByPosition(turret.Position) {
				fmt.Fprint(w, "error-exists")
				return
			}
			expense := turretAddCost(base)
			affords := canAfford(base, expense)
			if !affords[0] {
				fmt.Fprint(w, "error-water")
				return
			} else if !affords[1] {
				fmt.Fprint(w, "error-metal")
				return
			} else if !affords[2] {
				fmt.Fprint(w, "error-fuel")
				return
			} else if !affords[3] {
				fmt.Fprint(w, "error-people")
				return
			}
			doCosts(base, expense)
			id := uuid.New()
			turret.ID = id.String()
			turret.Level = 1
			writeTurret(username, turret)
		} else if act == "remove" {
			turret := Turret{}
			err = json.Unmarshal([]byte(requests[2]), &turret)
			if !check(err) {
				fmt.Fprint(w, "error-decode")
				return
			}
			if !base.hasTurretByID(turret.ID) {
				fmt.Fprint(w, "error-exists")
				return
			}
			removeTurret(username, turret.ID)
		} else if act == "change" {
			turret := Turret{}
			err = json.Unmarshal([]byte(requests[2]), &turret)
			if !check(err) {
				fmt.Fprint(w, "error-decode")
				return
			}
			if base.hasTurretByID(turret.ID) {
				fmt.Fprint(w, "error-exists")
				return
			}
			original := base.getTurretByID(turret.ID)
			if original == turret {
				fmt.Fprint(w, "error-same")
				return
			}
			if turret.Position != original.Position && turret.Level != original.Level {
				fmt.Fprint(w, "error-multiple")
				return
			} else if turret.Position != original.Position {
				curpostur := base.getTurretByPosition(turret.Position)
				if curpostur != (Turret{}) {
					updateTurret(username, curpostur.ID, bson.D{{"$set", bson.D{{"Position", original.Position}}}})
				}
				updateTurret(username, turret.ID, bson.D{{"$set", bson.D{{"Position", turret.Position}}}})

			} else if turret.Level != original.Level {
				expense := turretLevelCost(turret)
				affords := canAfford(base, expense)
				if !affords[0] {
					fmt.Fprint(w, "error-water")
					return
				} else if !affords[1] {
					fmt.Fprint(w, "error-metal")
					return
				} else if !affords[2] {
					fmt.Fprint(w, "error-fuel")
					return
				} else if !affords[3] {
					fmt.Fprint(w, "error-people")
					return
				}
				doCosts(base, expense)
				updateTurret(username, turret.ID, bson.D{{"$set", bson.D{{"Level", turret.Level}}}})
			}
		}
	} else if act == "ship" {
		act = requests[1]
		player := getPlayer(username)
		base := getBase(username)
		if act == "add" {
			ship := Ship{}
			err = json.Unmarshal([]byte(requests[2]), &ship)
			if !check(err) {
				fmt.Fprint(w, "error-decode")
				return
			}
			expense := shipAddCost(player)
			affords := canAfford(base, expense)
			if !affords[0] {
				fmt.Fprint(w, "error-water")
				return
			} else if !affords[1] {
				fmt.Fprint(w, "error-metal")
				return
			} else if !affords[2] {
				fmt.Fprint(w, "error-fuel")
				return
			} else if !affords[3] {
				fmt.Fprint(w, "error-people")
				return
			}
			doCosts(base, expense)
			id := uuid.New().String()
			ship.ID = id
			ship.Level = 1
			ship.Crew = 1
			ship.Defense = 1
			ship.Strength = 1
			writeShip(username, ship)

		} else if act == "change" {
			ship := Ship{}
			err = json.Unmarshal([]byte(requests[2]), &ship)
			if !check(err) {
				fmt.Fprint(w, "error-decode")
				return
			}
			if base.hasTurretByID(ship.ID) {
				fmt.Fprint(w, "error-exists")
				return
			}
			original := player.getShipByID(ship.ID) // TODO player get ship
			if original == ship {
				fmt.Fprint(w, "error-same")
				return
			}
			if ship.Level != original.Level {
				expense := shipLevelCost(ship)
				affords := canAfford(base, expense)
				if !affords[0] {
					fmt.Fprint(w, "error-water")
					return
				} else if !affords[1] {
					fmt.Fprint(w, "error-metal")
					return
				} else if !affords[2] {
					fmt.Fprint(w, "error-fuel")
					return
				} else if !affords[3] {
					fmt.Fprint(w, "error-people")
					return
				}
				doCosts(base, expense)
				updateShip(username, ship.ID, bson.D{{"$set", bson.D{{"Level", ship.Level}}}})
			}
		}
	}
	fmt.Fprint(w, "done")
}

// Water Metal Fuel People
func canAfford(base Base, costs []int) []bool {
	return []bool{base.Water >= costs[0], base.Metal >= costs[1], base.Fuel >= costs[2], base.People >= costs[3]}
}

func doCosts(base Base, costs []int) {
	updateBase(base.Owner, bson.D{{"$inc", bson.D{{"Water", costs[0]}, {"Metal", costs[1]}, {"Fuel", costs[2]}, {"People", costs[3]}}}})
}

// Formula: #turrets^2 * turretMultiplier * turretBaseCost
func turretAddCost(base Base) []int {
	multiplier := 1
	baseCost := 1
	costs := []int{0, 0, 0, 0}
	costs[1] = int(math.Pow(float64(len(base.Turrets)), float64(2))) * multiplier * baseCost
	return costs
}

// Formula: level^2 * turretMultiplier * turretBaseCost
// Formula: level^1.2 * turretMultiplier * turretBaseCost
func turretLevelCost(turret Turret) []int {
	multiplier := 1
	baseCost := 1
	costs := []int{0, 0, 0, 0}
	costs[1] = int(math.Pow(float64(turret.Level), float64(2))) * multiplier * baseCost
	costs[2] = int(math.Pow(float64(turret.Level), 1.2)) * multiplier * baseCost
	return costs
}

// Formula: #ships^2 * shipMultiplier * shipBaseCost
func shipAddCost(player Player) []int {
	multiplier := 1
	baseCost := 1
	costs := []int{0, 0, 0, 0}
	costs[1] = int(math.Pow(float64(len(player.Ships)), float64(2))) * multiplier * baseCost
	return costs
}

// Formula: level^2 * shipMultiplier * shipBaseCost
// Formula: level^1.2 * shipMultiplier * shipBaseCost
func shipLevelCost(ship Ship) []int {
	multiplier := 1
	baseCost := 1
	costs := []int{0, 0, 0, 0}
	costs[1] = int(math.Pow(float64(ship.Level), float64(2))) * multiplier * baseCost
	costs[2] = int(math.Pow(float64(ship.Level), 1.2)) * multiplier * baseCost
	return costs
}
