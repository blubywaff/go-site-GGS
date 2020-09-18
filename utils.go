package main

import (
	"fmt"
)

var logs []string

func check(err error) bool {
	if err != nil {
		logger(err.Error())
		fmt.Println(err)
	}
	return err == nil
}

func logger(issue string) {
	logs = append(logs, issue)
}
