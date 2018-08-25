package sherlock

import (
	"time"
	"io"
	"strings"
	"encoding/json"
)

/*
The most basic exception that sherlock offers.
Implements error, Loggable, and StackTraceWrapper.
*/
type StdException struct {
	stackTrace []*StackTraceEntry
	stackTraceStr string
	maxStackTraceSize int
	message string
	timestamp *time.Time // UTC
}

/*
Creates a new exception. A stack trace is created immediately. Stack trace depth is limited to 64 by default.

The stack trace does not get converted to a string until GetStackTraceAsString is called. Waiting to do this
until it actually gets logged vastly improves performance. I have noticed a performance of about
1800 ns/op to 2000 ns/op for this function on my desktop with Intel® Core™ i7-6700 CPU @ 3.40GHz × 8
running Ubuntu 18.04.1. This is about 5x faster than creating an exception in Java.
*/
func NewStdException(message string) *StdException {
	// Skip the top 4 functions in the stack trace so that the caller of this function is shown at the top
	return newStdException(message, defaultStackTraceNumLines, 4)
}

/*
Creates a new exception. A stack trace is created immediately. stackTraceNumLines allows
you to limit the depth of the stack trace.

The stack trace does not get converted to a string until GetStackTraceAsString is called. Waiting to do this
until it actually gets logged vastly improves performance. I have noticed a performance of about
1800 ns/op to 2000 ns/op for this function on my desktop with Intel® Core™ i7-6700 CPU @ 3.40GHz × 8
running Ubuntu 18.04.1. This is about 5x faster than creating an exception in Java.
*/
func NewStdExceptionWithStackTraceSize(message string, stackTraceNumLines int) *StdException {
	// Skip the top 4 functions in the stack trace so that the caller of this function is shown at the top
	return newStdException(message, stackTraceNumLines, 4)
}

func newStdException(message string, stackTraceNumLines, skip int) *StdException {
	timestamp := time.Now().UTC()
	return &StdException{
		stackTrace:        getStackTrace(skip, stackTraceNumLines),
		maxStackTraceSize: stackTraceNumLines,
		message:           message,
		timestamp:         &timestamp,
	}
}

/*
Returns the stack trace as slice of *StackTraceEntry.
*/
func (se *StdException) GetStackTrace() []*StackTraceEntry {
	return se.stackTrace
}

/*
Returns the stack trace in a string formatted as:

		sherlock.exampleFunc(exampleFile.go:18)
		sherlock.exampleFunc2(exampleFile2.go:46)
		sherlock.exampleFunc3(exampleFile2.go:177)

Uses the cached stack trace string if one is available.
If it has to convert the stack trace to a string, it will cache it for later.
*/
func (se *StdException) GetStackTraceAsString() string {
	if se.stackTraceStr == "" {
		se.stackTraceStr = stackTraceAsString(se.stackTrace)
	}

	return se.stackTraceStr
}

/*
Writes to the writer a string formatted as:

	yyyy-mm-dd hh:mm:ss - message:
		sherlock.exampleFunc(exampleFile.go:18)
		sherlock.exampleFunc2(exampleFile2.go:46)
		sherlock.exampleFunc3(exampleFile2.go:177)

Time is UTC.
Returns the string that was logged or an error if there was one.
*/
func (se *StdException) Log(writer io.Writer) ([]byte, error) {
	var buf strings.Builder
	buf.WriteString(se.createCompactMessage())
	buf.WriteString(":\n")
	buf.WriteString(se.GetStackTraceAsString())
	logMessage := []byte(buf.String())

	_, err := writer.Write(logMessage)
	if err != nil {
		return nil, err
	}

	return logMessage, nil
}

/*
Writes to the writer a string formatted as:

	yyyy-mm-dd hh:mm:ss - message

Time is UTC.
Note that it does not have the stack trace.
Returns the string that was logged or an error if there was one.
*/
func (se *StdException) LogNoStack(writer io.Writer) ([]byte, error) {
	logMessage := []byte(se.createCompactMessage())

	_, err := writer.Write(logMessage)
	if err != nil {
		return nil, err
	}
	return logMessage, nil
}

/*
Packages up the exception's info into json and writes it to writer.
Returns the logged message or an error if there was one.
*/
func (se *StdException) LogAsJson(writer io.Writer) (jsonBytes []byte, err error){
	jsonBytes, err = se.toJsonBytes()
	if err != nil {
		return
	}

	_, err = writer.Write(jsonBytes)
	if err != nil {
		return nil, err
	}

	return
}

/*
Returns the message and stack trace in a string formatted like this:

	message:
		sherlock.exampleFunc(exampleFile.go:18)
		sherlock.exampleFunc2(exampleFile2.go:46)
		sherlock.exampleFunc3(exampleFile2.go:177)

Leaves out the timestamp so that StdException will print nicely with log.Println
*/
func (se *StdException) Error() string {
	var buf strings.Builder
	buf.WriteString(se.message)
	buf.WriteString(":\n")
	buf.WriteString(se.GetStackTraceAsString())
	return buf.String()
}

/*
Returns the timestamp and message as

	yyyy-mm-dd hh:mm:ss - message

Time is UTC.
*/
func (se *StdException) createCompactMessage() string {
	var buf strings.Builder
	buf.WriteString(se.timestamp.Format(timeFmt))
	buf.WriteString(" - ")
	buf.WriteString(se.message)
	return buf.String()
}

/*
Returns the timestamp, message, and stack trace as:

	yyyy-mm-dd hh:mm:ss - message:
		sherlock.exampleFunc(exampleFile.go:18)
		sherlock.exampleFunc2(exampleFile2.go:46)
		sherlock.exampleFunc3(exampleFile2.go:177)

Time is UTC.
*/
func (se *StdException) createLogMessage() string {
	var buf strings.Builder
	buf.WriteString(se.createCompactMessage())
	buf.WriteString(":\n")
	buf.WriteString(se.GetStackTraceAsString())
	return buf.String()
}

func (se *StdException) toJsonBytes() ([]byte, error) {
	return json.Marshal(se.toJsonMap())
}

func (se *StdException) toJsonMap() map[string]interface{} {
	return map[string]interface{}{
		"Time": se.timestamp.Format(timeFmt),
		"Message": se.message,
		"StackTrace": se.stackTrace,
		"StackTraceStr": se.GetStackTraceAsString(),
	}
}