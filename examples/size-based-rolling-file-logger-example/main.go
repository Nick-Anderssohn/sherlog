package main

import (
	"fmt"
	"sherlog"
	"sherlog/examples/exception-returner"
	"sherlog/examples/size-based-rolling-file-logger-example/exlogger"
	"errors"
)

func main() {
	// the rolling logger should create multiple files as they fill up. we set it to be limited to 5 messages per file.
	for i := 0; i < 4; i++ {
		err := exception_returner.ExampleReturnError()
		if err != nil {
			exlogger.Logger.Log(err)
		}

		err = exception_returner.ExampleReturnOpsError()
		if err != nil {
			exlogger.Logger.Log(err)
		}

		err = ExampleReturnWarning()
		if err != nil {
			exlogger.Logger.Log(err)
		}

		err = exception_returner.ExampleReturnCustomLeveledException()
		if err != nil {
			exlogger.Logger.Log(err)
		}

		err = errors.New("test an accidental non-sherlog error to see that it is handled correctly")
		exlogger.Logger.Log(err)
	}
}


func ExampleReturnWarning() error {
	potentialWarning := doSomethingThatTakesLongerThanExpected()
	fmt.Println("I'm doing other stuff here without checking the error because I know it's either nil or just a warning")
	return potentialWarning
}

func doSomethingThatTakesLongerThanExpected() error {
	fmt.Println("Wow for some reason something happened here that took 100x longer than it normally does")
	thingTookLongerThanNormal := true
	if thingTookLongerThanNormal {
		return sherlog.NewWarning("thing took longer than it normally does")
	}
	return nil
}