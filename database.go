package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

var dbTimeFormat = "01/02/2006 15:04:05"

// search this for that in column of table
func contains(data *sql.DB, find string, col string, table string) bool {
	rs, err := data.Query("select " + col + " from " + table + " where " + col + "='" + find + "';")
	if rs == nil {
		return false
	}
	defer rs.Close()
	if err != nil {
		logger(err.Error())
		fmt.Println(err)
		return false
	}
	rs.Next()
	var result string
	err = rs.Scan(&result)
	return err == nil && result == find
}

// search this for this in column of table and return from findcol
func find(data *sql.DB, find string, col string, table string, resultCol string) string {
	rs, err := data.Query("select " + resultCol + " from " + table + " where " + col + "='" + find + "';")
	defer rs.Close()
	if err != nil {
		logger(err.Error())
		fmt.Println(err)
		return ""
	}
	rs.Next()
	var result string
	err = rs.Scan(&result)
	return result
}

// remove from this when column in table equals this
func remove(data *sql.DB, col string, table string, remove string) {
	/*stmt, err := data.Prepare("delete from " + table + " where " + col + "='" + remove + "';")
	check(err)
	defer stmt.Close()
	r, err := stmt.Exec()
	check(err)
	_, err = r.RowsAffected()
	check(err)*/
	_, err := data.Exec("delete from " + table + " where " + col + "='" + remove + "';")
	check(err)
}

func writeUser(username string, password string, firstname string, lastname string) {
	/*stmt, err := usersdb.Prepare("insert into users values ('" + username + "', '" + password + "', '" + firstname + "', '" + lastname + "';")
	check(err)
	defer stmt.Close()
	r, err := stmt.Exec()
	check(err)
	_, err = r.RowsAffected()
	check(err)*/
	_, err := usersdb.Exec("insert into users values ('" + username + "', '" + password + "', '" + firstname + "', '" + lastname + "');")
	check(err)
}

func writeSession(sessionID string, username string, lastActivity string) {
	/*stmt, err := usersdb.Prepare("insert into sessions values ('" + sessionID + "', '" + username + "', '" + lastActivity + "';")
	check(err)
	defer stmt.Close()
	r, err := stmt.Exec()
	check(err)
	_, err = r.RowsAffected()
	check(err)*/
	_, err := usersdb.Exec("insert into sessions values ('" + sessionID + "', '" + username + "', '" + lastActivity + "');")
	check(err)
}

func updateSession(sessionID string, lastActivity string) {
	_, err := usersdb.Exec("update sessions set lastActivity = '" + lastActivity + "' where sessionID = '" + sessionID + "';")
	//fmt.Println(err)
	check(err)
}