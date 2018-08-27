package main

import (
	"sherlog"
	"sherlog/examples/nightly-rolling-file-logger-example/exlogger"
	"time"
)

func main() {
	var done bool

	// kill this example after 25 hours
	go func() {
		<-time.After(25 * time.Hour)
		done = true
	}()

	for !done {
		err := sherlog.NewInfo("Still testing")
		exlogger.Logger.Log(err)
		time.Sleep(time.Minute)
	}
}
