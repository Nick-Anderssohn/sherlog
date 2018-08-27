package exlogger

import (
	"sherlog"
)

// I recommend you create your own logger package in your project to hold the singleton instance
// of a sherlog logger

var Logger sherlog.Logger

func init() {
	var err error
	// I want all log messages to go into one rolling log file. I want the file to roll to another file every 5 messages.
	Logger, err = sherlog.NewRollingFileLoggerWithSizeLimit("rolling_log.log", 5)

	// If logging fails to get setup, I don't even want my program to start.
	if err != nil {
		panic(err)
	}
}
