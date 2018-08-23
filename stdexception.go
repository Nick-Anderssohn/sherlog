package logging

import (
	"runtime"
	"time"
	"io"
	"fmt"
	"encoding/json"
)

// A standard exception
type StdException struct {
	stackTrace []*runtime.Frame
	maxStackTraceSize int
	message string
	timestamp time.Time
}

func NewStdException(message string) *StdException {
	return NewStdExceptionWithStackTraceSize(message, defaultStackTraceSize)
}

func NewStdExceptionWithStackTraceSize(message string, stackTraceSize int) *StdException {
	return &StdException{
		stackTrace: getStackTrace(2, stackTraceSize),
		maxStackTraceSize: stackTraceSize,
		message: message,
		timestamp: time.Now().UTC(),
	}
}

// Writes "timestamp - message" to writer.
// On failure, it will pass an error to failureHandler.
func (se *StdException) LogCompactFmt(writer io.Writer, failureHandler FailureHandler) {
	_, err := writer.Write([]byte(fmt.Sprintf("%s - %s", se.timestamp.Format(timeFmt), se.message)))

	if err != nil {
		failureHandler.HandleFail(err)
	}
}

// Packages up the exception's info into json and writes it to writer.
// On failure, it will pass an error to failureHandler.
func (se *StdException) LogAsJson(writer io.Writer, failureHandler FailureHandler) {
	jsonBytes, err := json.Marshal(map[string]interface{}{
		"Time": se.timestamp.Format(timeFmt),
		"Message": se.message,
		"StackTrace": se.stackTrace,
	})

	if err != nil {
		failureHandler.HandleFail(err)
		return
	}

	_, err = writer.Write(jsonBytes)

	if err != nil {
		failureHandler.HandleFail(err)
	}
}