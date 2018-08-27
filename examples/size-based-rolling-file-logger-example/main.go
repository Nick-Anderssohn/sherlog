package main

import (
	"errors"

	"github.com/Nick-Anderssohn/sherlog/examples/ex"
	"github.com/Nick-Anderssohn/sherlog/examples/size-based-rolling-file-logger-example/exlogger"
)

func main() {
	// the rolling logger should create multiple files as they fill up. we set it to be limited to 5 messages per file.
	for i := 0; i < 4; i++ {
		err := ex.ExampleReturnOpsError()
		if err != nil {
			exlogger.Logger.Log(err)
		}

		err = ex.ExampleReturnError()
		if err != nil {
			exlogger.Logger.Log(err)
		}

		err = ex.ExampleReturnCustomLeveledException()
		if err != nil {
			exlogger.Logger.Log(err)
		}

		err = errors.New("test an accidental non-sherlog error to see that it is handled correctly")
		exlogger.Logger.Log(err)
	}
}
