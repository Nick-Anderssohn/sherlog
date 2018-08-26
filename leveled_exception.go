package sherlog

import (
	"io"
	"encoding/json"
	"strings"
)

/*
An interface used to specify the log level on an exception/error.
LevelId is meant to be something along the lines of an enum, so
that we don't have to switch based off of the string value of the
log level. Label is the string representation.
 */
type Level interface {
	GetLevelId() int
	GetLabel() string
}

/*
Something that holds a modifiable log level.
*/
type LevelWrapper interface {
	GetLevel() Level
	SetLevel(level Level)
}

func isLevelWrapper(err error) bool {
	_, isLeveled := err.(LevelWrapper)
	return isLeveled
}

/*
An exception with a level such as ERROR or WARNING.
StdException is embedded.
Implements error, LevelWrapper, Loggable, StackTraceWrapper, and LeveledLoggable.
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

func (le *LeveledException) SetLevel(level Level) {
	le.level = level
}

/*
Creates a new LeveledException. A stack trace is created immediately. Stack trace depth is limited to 64 by default.

The stack trace does not get converted to a string until GetStackTraceAsString is called. Waiting to do this
until it actually gets logged vastly improves performance. I have noticed a performance of about
2000 ns/op for this function on my desktop with Intel® Core™ i7-6700 CPU @ 3.40GHz × 8
running Ubuntu 18.04.1. This is about 5x faster than creating an exception in Java.
*/
func NewLeveledException(message string, level Level) *LeveledException {
	return newLeveledException(message, level, defaultStackTraceDepth, 5)
}

/*
Creates a new LeveledException. A stack trace is created immediately. Stack trace depth is limited to maxStackTraceDepth.

The stack trace does not get converted to a string until GetStackTraceAsString is called. Waiting to do this
until it actually gets logged vastly improves performance. I have noticed a performance of about
2000 ns/op  for this function on my desktop with Intel® Core™ i7-6700 CPU @ 3.40GHz × 8
running Ubuntu 18.04.1. This is about 5x faster than creating an exception in Java.
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
		sherlog.exampleFunc(exampleFile.go:18)
		sherlog.exampleFunc2(exampleFile2.go:46)
		sherlog.exampleFunc3(exampleFile2.go:177)

Time is UTC.
Returns the string that was logged or an error if there was one.
*/
func (le *LeveledException) Log(writer io.Writer) error {
	err := le.LogNoStack(writer)
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(":\n"))
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(le.GetStackTraceAsString()))
	return err
}

/*
Writes to the writer a string formatted as:

	yyyy-mm-dd hh:mm:ss - LEVEL - message

Time is UTC.
Note that it does not have the stack trace.
Returns the string that was logged or an error if there was one.
*/
func (le *LeveledException) LogNoStack(writer io.Writer) error {
	_, err := writer.Write([]byte(le.timestamp.Format(timeFmt)))
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(" - "))
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(le.level.GetLabel()))
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(" - "))
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(le.message))
	return err
}

/*
Packages up the exception's info into json and writes it to writer.
Returns returns the logged message or an error if there was one.
*/
func (le *LeveledException) LogAsJson(writer io.Writer) error {
	jsonBytes, err := le.toJsonBytes()

	if err != nil {
		return err
	}

	_, err = writer.Write(jsonBytes)

	return err
}

/*
Returns the message and stack trace in a string formatted like this:

	LEVEL - message:
		sherlog.exampleFunc(exampleFile.go:18)
		sherlog.exampleFunc2(exampleFile2.go:46)
		sherlog.exampleFunc3(exampleFile2.go:177)

Leaves out the timestamp so that LeveledException will print nicely with log.Println
*/
func (le *LeveledException) Error() string {
	var buf strings.Builder
	buf.WriteString(" - ")
	buf.WriteString(le.level.GetLabel())
	buf.WriteString(" - ")
	buf.WriteString(le.message)
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