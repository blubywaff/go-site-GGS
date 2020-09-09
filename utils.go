package main

import (
	"fmt"
)

var logs []string

func check(err error) {
	if err != nil {
		logger(err.Error())
		fmt.Println(err)
	}
}

func logger(issue string) {
	logs = append(logs, issue)
}
