package exlogger

import (
	"time"
	"github.com/Nick-Anderssohn/sherlog"
)

// I recommend you create your own logger package in your project to hold the singleton instance
// of a sherlog logger

var Logger sherlog.Logger

func init() {
	var err error
	// I want sherlog to use pacific time instead of UTC (which is the default)
	sherlog.SherlogLocation, err = time.LoadLocation("America/Los_Angeles")
	if err != nil {
		// If logging fails to get setup, I don't even want my program to start.
		panic(err)
	}

	// I want all log messages to go into one rolling log file. I want the file to roll every midnight.
	Logger, err = sherlog.NewNightlyRollingFileLogger("nightly_rolling_log.log")
	// If logging fails to get setup, I don't even want my program to start.
	if err != nil {
		panic(err)
	}
}
