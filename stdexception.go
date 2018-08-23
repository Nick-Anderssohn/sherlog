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
// Returns returns the logged message or an error if there is one.
func (se *StdException) LogCompactFmt(writer io.Writer) ([]byte, error) {
	logMessage := []byte(fmt.Sprintf("%s - %s", se.timestamp.Format(timeFmt), se.message))
	_, err := writer.Write(logMessage)
	if err != nil {
		return nil, err
	}
	return logMessage, nil
}

// Packages up the exception's info into json and writes it to writer.
// Returns returns the logged message or an error if there is one.
func (se *StdException) LogAsJson(writer io.Writer) (jsonBytes []byte, err error){
	jsonBytes, err = json.Marshal(map[string]interface{}{
		"Time": se.timestamp.Format(timeFmt),
		"Message": se.message,
		"StackTrace": se.stackTrace,
	})

	if err != nil {
		return
	}

	_, err = writer.Write(jsonBytes)
	if err != nil {
		return nil, err
	}

	return
}