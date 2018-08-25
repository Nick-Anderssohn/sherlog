package sherlock

import (
	"io"
	"encoding/json"
	"strings"
)

type Level interface {
	GetLevelId() int
	GetLabel() string
}

/*
An exception with a level such as ERROR or WARNING.
StdException is embedded.
Implements error, Loggable, StackTraceWrapper, and LeveledLoggable.
 */
type LeveledException struct {
	// If we really really wanted to, we could save about 100 ns/op in the constructor if we changed
	// the embedded StdException into a pointer/not-embedded field (so stdException *StdException)
	StdException
	level Level
}

func (le *LeveledException) GetLevel() Level {
	return le.level
}

/*
Creates a new LeveledException. A stack trace is created immediately. Stack trace depth is limited to 64 by default.

The stack trace does not get converted to a string until GetStackTraceAsString is called. Waiting to do this
until it actually gets logged vastly improves performance. I have noticed a performance of about
2000 ns/op for this function on my desktop with Intel® Core™ i7-6700 CPU @ 3.40GHz × 8
running Ubuntu 18.04.1. This is about 5 times faster than creating an exception in Java.
*/
func NewLeveledException(message string, level Level) *LeveledException {
	return newLeveledException(message, level, defaultStackTraceNumLines, 5)
}

/*
Creates a new LeveledException. A stack trace is created immediately. Stack trace depth is limited to 64 by default.

The stack trace does not get converted to a string until GetStackTraceAsString is called. Waiting to do this
until it actually gets logged vastly improves performance. I have noticed a performance of about
2000 ns/op  for this function on my desktop with Intel® Core™ i7-6700 CPU @ 3.40GHz × 8
running Ubuntu 18.04.1. This is about 5 times faster than creating an exception in Java.
*/
func NewLeveledExceptionWithStackTraceSize(message string, level Level, maxStackTraceDepth int) *LeveledException {
	return newLeveledException(message, level, maxStackTraceDepth, 5)
}

func newLeveledException(message string, level Level, maxStackTraceDepth, skip int) *LeveledException {
	return &LeveledException{
		StdException: *newStdException(message, maxStackTraceDepth, skip),
		level:        level,
	}
}

/*
Writes to the writer a string formatted as:

	yyyy-mm-dd hh:mm:ss - LEVEL - message:
		sherlock.exampleFunc(exampleFile.go:18)
		sherlock.exampleFunc2(exampleFile2.go:46)
		sherlock.exampleFunc3(exampleFile2.go:177)

Time is UTC.
Returns the string that was logged or an error if there was one.
*/
func (le *LeveledException) Log(writer io.Writer) ([]byte, error) {
	var buf strings.Builder
	buf.WriteString(le.createCompactMessage())
	buf.WriteString(":\n")
	buf.WriteString(le.GetStackTraceAsString())
	logMessage := []byte(buf.String())

	_, err := writer.Write(logMessage)
	if err != nil {
		return nil, err
	}

	return logMessage, nil
}

/*
Writes to the writer a string formatted as:

	yyyy-mm-dd hh:mm:ss - LEVEL - message

Time is UTC.
Note that it does not have the stack trace.
Returns the string that was logged or an error if there was one.
*/
func (le *LeveledException) LogNoStack(writer io.Writer) ([]byte, error) {
	logMessage := []byte(le.createCompactMessage())
	_, err := writer.Write(logMessage)
	if err != nil {
		return nil, err
	}
	return logMessage, nil
}

/*
Packages up the exception's info into json and writes it to writer.
Returns returns the logged message or an error if there was one.
*/
func (le *LeveledException) LogAsJson(writer io.Writer) ([]byte, error) {
	jsonBytes, err := le.toJsonBytes()

	if err != nil {
		return nil, err
	}

	_, err = writer.Write(jsonBytes)
	if err != nil {
		return nil, err
	}

	return jsonBytes, err
}

/*
Returns the message and stack trace in a string formatted like this:

	LEVEL - message:
		sherlock.exampleFunc(exampleFile.go:18)
		sherlock.exampleFunc2(exampleFile2.go:46)
		sherlock.exampleFunc3(exampleFile2.go:177)

Leaves out the timestamp so that LeveledException will print nicely with log.Println
*/
func (le *LeveledException) Error() string {
	var buf strings.Builder
	buf.WriteString(le.level.GetLabel())
	buf.WriteString(" - ")
	buf.WriteString(le.message)
	buf.WriteString(":\n")
	buf.WriteString(le.GetStackTraceAsString())
	return buf.String()
}

/*
Returns the timestamp and message as:

	yyyy-mm-dd hh:mm:ss - LEVEL - message

Time is UTC.
*/
func (le *LeveledException) createCompactMessage() string {
	var buf strings.Builder
	buf.WriteString(le.timestamp.Format(timeFmt))
	buf.WriteString(" - ")
	buf.WriteString(le.level.GetLabel())
	buf.WriteString(" - ")
	buf.WriteString(le.message)
	return buf.String()
}

/*
Returns the timestamp, message, and stack trace as:

	yyyy-mm-dd hh:mm:ss - LEVEL - message:
		sherlock.exampleFunc(exampleFile.go:18)
		sherlock.exampleFunc2(exampleFile2.go:46)
		sherlock.exampleFunc3(exampleFile2.go:177)

Time is UTC.
*/
func (le *LeveledException) createLogMessage() string {
	var buf strings.Builder
	buf.WriteString(le.createCompactMessage())
	buf.WriteString(":\n")
	buf.WriteString(le.GetStackTraceAsString())
	return buf.String()
}

func (le *LeveledException) toJsonBytes() ([]byte, error) {
	return json.Marshal(le.toJsonMap())
}

func (le *LeveledException) toJsonMap() map[string]interface{} {
	jsonMap := le.StdException.toJsonMap()
	jsonMap["Level"] = le.level
	return jsonMap
}