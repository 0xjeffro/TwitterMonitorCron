package main

import (
	"TwitterMonitorCron/twitter"
	"TwitterMonitorCron/utils"
	"os"
	"strings"
)

func main() {
	utils.CheckEnvironment()
	targetStr := os.Getenv("TARGET_USER_ID")
	targets := strings.Split(targetStr, "|")
	for _, target := range targets {
		twitter.GetTweets(target)
	}
}
