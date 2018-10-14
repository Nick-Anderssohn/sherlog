package sherlog

import (
	"encoding/json"
	"io"
	"strings"
)

/*
Level is an interface used to specify the log level on an exception/error.
LevelId is meant to be something along the lines of an enum, so
that we don't have to switch based off of the string value of the
log level. Label is the string representation.
*/
type Level interface {
	GetLevelId() int
	GetLabel() string
}

/*
LevelWrapper is something that holds a modifiable log level.
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
LeveledException is an exception with a level such as ERROR or WARNING.
StdException is embedded.
Implements error, LevelWrapper, Loggable, StackTraceWrapper, and LeveledLoggable.
*/
type LeveledException struct {
	// If we really really wanted to, we could save about 100 ns/op in the constructor if we changed
	// the embedded StdException into a pointer/not-embedded field (so stdException *StdException)
	StdException
	level Level
}

/*
GetLevel returns the level.
*/
func (le *LeveledException) GetLevel() Level {
	return le.level
}

/*
SetLevel sets the level.
*/
func (le *LeveledException) SetLevel(level Level) {
	le.level = level
}

/*
NewLeveledException creates a new LeveledException. A stack trace is created immediately. Stack trace depth is limited to 64 by default.

The stack trace does not get converted to a string until GetStackTraceAsString is called. Waiting to do this
until it actually gets logged vastly improves performance. I have noticed a performance of about
2000 ns/op for this function on my desktop with Intel® Core™ i7-6700 CPU @ 3.40GHz × 8
running Ubuntu 18.04.1. This is about 5x faster than creating an exception in Java.
*/
func NewLeveledException(message string, level Level) error {
	return newLeveledException(message, level, defaultStackTraceDepth, 5)
}

/*
NewLeveledExceptionWithStackTraceSize creates a new LeveledException. A stack trace is created immediately. Stack trace depth is limited to maxStackTraceDepth.

The stack trace does not get converted to a string until GetStackTraceAsString is called. Waiting to do this
until it actually gets logged vastly improves performance. I have noticed a performance of about
2000 ns/op  for this function on my desktop with Intel® Core™ i7-6700 CPU @ 3.40GHz × 8
running Ubuntu 18.04.1. This is about 5x faster than creating an exception in Java.
*/
func NewLeveledExceptionWithStackTraceSize(message string, level Level, maxStackTraceDepth int) error {
	return newLeveledException(message, level, maxStackTraceDepth, 5)
}

func newLeveledException(message string, level Level, maxStackTraceDepth, skip int) error {
	return &LeveledException{
		StdException: *newStdException(message, maxStackTraceDepth, skip),
		level:        level,
	}
}

/*
Log writes to the writer a string formatted as:

	yyyy-mm-dd hh:mm:ss - LEVEL - message:
		sherlog.exampleFunc(exampleFile.go:18)
		sherlog.exampleFunc2(exampleFile2.go:46)
		sherlog.exampleFunc3(exampleFile2.go:177)

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
LogNoStack writes to the writer a string formatted as:

	yyyy-mm-dd hh:mm:ss - LEVEL - message

Note that it does not have the stack trace.
Returns the string that was logged or an error if there was one.
*/
func (le *LeveledException) LogNoStack(writer io.Writer) error {
	for _, msg := range le.messageChain {
		writer.Write([]byte(msg))
		writer.Write([]byte("\nCaused by:\n"))
	}
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
LogAsJson packages up the exception's info into json and writes it to writer.

The json is formatted like this
	{
	   "Level":"INFO",
	   "Message":"I'm informative!",
	   "StackTrace":[
		  {
			 "FunctionName":"github.com/Nick-Anderssohn/sherlog.TestLogJson",
			 "File":"/home/nick/go/src/github.com/Nick-Anderssohn/sherlog/scratch_test.go",
			 "Line":68
		  },
		  {
			 "FunctionName":"testing.tRunner",
			 "File":"/usr/local/go/src/testing/testing.go",
			 "Line":777
		  },
		  {
			 "FunctionName":"runtime.goexit",
			 "File":"/usr/local/go/src/runtime/asm_amd64.s",
			 "Line":2361
		  }
	   ],
	   "StackTraceStr":"\tgithub.com/Nick-Anderssohn/sherlog.TestLogJson(/home/nick/go/src/github.com/Nick-Anderssohn/sherlog/scratch_test.go:68)\n\ttesting.tRunner(/usr/local/go/src/testing/testing.go:777)\n\truntime.goexit(/usr/local/go/src/runtime/asm_amd64.s:2361)\n",
	   "Time":"2018-10-03 07:51:14"
	}

Returns an error if there was one.
*/
func (le *LeveledException) LogAsJson(writer io.Writer) error {
	jsonBytes, err := le.ToJsonBytes()

	if err != nil {
		return err
	}

	_, err = writer.Write(jsonBytes)

	return err
}

/*
Error returns the message and stack trace in a string formatted like this:

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

/*
ToJsonBytes returns the bytes for a json blob that looks like this:

	{
	   "Level":"INFO",
	   "Message":"I'm informative!",
	   "StackTrace":[
		  {
			 "FunctionName":"github.com/Nick-Anderssohn/sherlog.TestLogJson",
			 "File":"/home/nick/go/src/github.com/Nick-Anderssohn/sherlog/scratch_test.go",
			 "Line":68
		  },
		  {
			 "FunctionName":"testing.tRunner",
			 "File":"/usr/local/go/src/testing/testing.go",
			 "Line":777
		  },
		  {
			 "FunctionName":"runtime.goexit",
			 "File":"/usr/local/go/src/runtime/asm_amd64.s",
			 "Line":2361
		  }
	   ],
	   "StackTraceStr":"\tgithub.com/Nick-Anderssohn/sherlog.TestLogJson(/home/nick/go/src/github.com/Nick-Anderssohn/sherlog/scratch_test.go:68)\n\ttesting.tRunner(/usr/local/go/src/testing/testing.go:777)\n\truntime.goexit(/usr/local/go/src/runtime/asm_amd64.s:2361)\n",
	   "Time":"2018-10-03 07:51:14"
	}
*/
func (le *LeveledException) ToJsonBytes() ([]byte, error) {
	return json.Marshal(le.ToJsonMap())
}

/*
ToJsonMap creates a map[string]interface{} that, when compiled to json, looks like this:

	{
	   "Level":"INFO",
	   "Message":"I'm informative!",
	   "StackTrace":[
		  {
			 "FunctionName":"github.com/Nick-Anderssohn/sherlog.TestLogJson",
			 "File":"/home/nick/go/src/github.com/Nick-Anderssohn/sherlog/scratch_test.go",
			 "Line":68
		  },
		  {
			 "FunctionName":"testing.tRunner",
			 "File":"/usr/local/go/src/testing/testing.go",
			 "Line":777
		  },
		  {
			 "FunctionName":"runtime.goexit",
			 "File":"/usr/local/go/src/runtime/asm_amd64.s",
			 "Line":2361
		  }
	   ],
	   "StackTraceStr":"\tgithub.com/Nick-Anderssohn/sherlog.TestLogJson(/home/nick/go/src/github.com/Nick-Anderssohn/sherlog/scratch_test.go:68)\n\ttesting.tRunner(/usr/local/go/src/testing/testing.go:777)\n\truntime.goexit(/usr/local/go/src/runtime/asm_amd64.s:2361)\n",
	   "Time":"2018-10-03 07:51:14"
	}
*/
func (le *LeveledException) ToJsonMap() map[string]interface{} {
	jsonMap := le.StdException.ToJsonMap()
	jsonMap["Level"] = le.level.GetLabel()
	return jsonMap
}
