package main

import (
	"TwitterMonitorCron/twitter"
	"TwitterMonitorCron/utils"
)

func main() {
	utils.CheckEnvironment()
	twitter.GetTweets()
}
