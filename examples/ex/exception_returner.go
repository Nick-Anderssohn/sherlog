package ex

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Nick-Anderssohn/sherlog"
)

// CustomLogLevel will implement sherlog.Level so that it will integrate seamlessly with sherlog
type CustomLogLevel int

const (
	// WeirdLogLevel starts at 100 so that it doesn't intersect with any default sherlog
	// log level enums.
	WeirdLogLevel CustomLogLevel = 100
)

var customLevelLabels = map[CustomLogLevel]string{
	WeirdLogLevel: "WEIRD_LOG_LEVEL",
}

// Create the 2 needed functions to implement sherlog.Level so that this custom log level can be used
// by sherlog

// GetLevelId to implement interface sherlog.Level
func (logLevel CustomLogLevel) GetLevelId() int {
	return int(logLevel)
}

// GetLabel to implement sherlog.Level
func (logLevel CustomLogLevel) GetLabel() string {
	return customLevelLabels[logLevel]
}

// ExampleReturnError will return an error with log level ERROR
func ExampleReturnError() error {
	err := funcThatReturnsError()
	if err != nil {
		// funcThatReturnsError already called sherlog.AsError, but
		// if I called sherlog.AsError again, it wouldn't matter. The AsError would detect that the
		// error already has a stack trace and would not overwrite it. However, it would overwrite the
		// log level if it had a different one. In this case the log level is already ERROR, so it would be overwritten
		// to the same value anyways. All of the AsSomeLogLevel functions work in this manner.

		// So, these end up resulting in the same thing:
		// return sherlog.AsError(err)
		return err
	}
	return nil
}

/*
An example of a function that has a bug and returns an error.
*/
func funcThatReturnsError() error {
	notAnInt := "I am not an int"
	_, err := strconv.Atoi(notAnInt)
	if err != nil {
		// If I get an error here, it probably means there is a bug in the code. So, I will give it
		// a log level of ERROR. It will also create a stack trace. I can keep the return type of this
		// function as error because sherlog exceptions implement the error interface.
		return sherlog.AsError(err)
	}
	fmt.Println("Doing other cool stuff here if there was no error")
	return nil
}

/*
ExampleReturnOpsError is an example of a function that fails to connect to a database, so it returns on OPS_ERROR.
*/
func ExampleReturnOpsError() error {
	err := queryDB()
	if err != nil {
		// since an error from a database query has potential to be something like a connection issue that
		// ops could deal with, we give it the log level OPS_ERROR (and a stack trace of course).
		// Prefer OPS_ERROR over ERROR if there is a chance that it could be an ops issue.
		return sherlog.AsOpsError(err)
	}
	return nil
}

// Pretend this is a function from some other library that tries to query a database, but fails due to a connection
// issue. Because that third party does not use sherlog, it returns an error without a stack trace
func queryDB() error {
	return errors.New("could not connect to postgres")
}

// ExampleReturnCustomLeveledException returns one of our custom exceptions for the example.
func ExampleReturnCustomLeveledException() error {
	// Oh shoots, I did something that must return an error with my own custom log level.
	return sherlog.NewLeveledException("i did something weird", WeirdLogLevel)
}

// ....You get the idea... we can do the same sort of stuff for any log level
