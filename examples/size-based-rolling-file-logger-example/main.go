package main

import (
	"errors"

	"github.com/Nick-Anderssohn/sherlog/examples/exception-returner"
	"github.com/Nick-Anderssohn/sherlog/examples/size-based-rolling-file-logger-example/exlogger"
)

func main() {
	// the rolling logger should create multiple files as they fill up. we set it to be limited to 5 messages per file.
	for i := 0; i < 4; i++ {
		err := exception_returner.ExampleReturnOpsError()
		if err != nil {
			exlogger.Logger.Log(err)
		}

		err = exception_returner.ExampleReturnError()
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
