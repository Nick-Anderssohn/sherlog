package logging

import (
	"time"
	"io"
	"encoding/json"
	"strings"
)

// A standard exception
type StdException struct {
	stackTrace []*StackTraceEntry
	maxStackTraceSize int
	message string
	timestamp time.Time
}

func NewStdException(message string) *StdException {
	return NewStdExceptionWithStackTraceSize(message, defaultStackTraceNumLines)
}

func NewStdExceptionWithStackTraceSize(message string, stackTraceNumLines int) *StdException {
	return &StdException{
		stackTrace:        getStackTrace(2, stackTraceNumLines),
		maxStackTraceSize: stackTraceNumLines,
		message:           message,
		timestamp:         time.Now().UTC(),
	}
}

func (se *StdException) GetStackTraceAsString() (string, error) {
	return StackTraceAsString(se.stackTrace)
}

func (se *StdException) Log(writer io.Writer) ([]byte, error) {
	var buf strings.Builder
	_, err := buf.WriteString(se.Error())
	if err != nil {
		return nil, err
	}
	_, err = buf.WriteString(":")
	if err != nil {
		return nil, err
	}
	stackTraceStr, err := se.GetStackTraceAsString()
	if err != nil {
		return nil, err
	}
	_, err = buf.WriteString(stackTraceStr)
	if err != nil {
		return nil, err
	}
	logMessage := []byte(buf.String())
	_, err = writer.Write(logMessage)
	if err != nil {
		return nil, err
	}
	return logMessage, nil
}

// Writes "timestamp - message" to writer.
// Returns returns the logged message or an error if there is one.
func (se *StdException) LogCompactFmt(writer io.Writer) ([]byte, error) {
	logMessage := []byte(se.Error())
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

// returns "timestamp - message"
// Unfortunately, this function has to panic if there is an error.
// It cannot return the error because it needs this method signature in order to
// implement the built in error interface
func (se *StdException) Error() string {
	var buf strings.Builder
	_, err := buf.WriteString(se.timestamp.Format(timeFmt))
	if err != nil {
		panic(err)
	}
	_, err = buf.WriteString(" - ")
	if err != nil {
		panic(err)
	}
	_, err = buf.WriteString(se.message)
	if err != nil {
		panic(err)
	}
	return buf.String()
}