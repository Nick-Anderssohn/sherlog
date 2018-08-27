package main

import (
	"time"
	"github.com/Nick-Anderssohn/sherlog"
	"github.com/Nick-Anderssohn/sherlog/examples/nightly-rolling-file-logger-example/exlogger"
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
