package sherlog

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

/*
StdException is the most basic exception that sherlog offers.
Implements error, Loggable, and StackTraceWrapper.
*/
type StdException struct {
	stackTrace        []*StackTraceEntry
	stackTraceStr     string
	maxStackTraceSize int
	message           string
	timestamp         *time.Time
	messageChain      []string

	// NonLoggedMsg can be optionally used to attach a secondary message that won't be logged.
	NonLoggedMsg string
}

type prependable interface {
	prependMsg(msg string)
}

/*
NewStdException creates a new exception. A stack trace is created immediately. Stack trace depth is limited to 64 by default.

The stack trace does not get converted to a string until GetStackTraceAsString is called. Waiting to do this
until it actually gets logged vastly improves performance. I have noticed a performance of about
1800 ns/op to 2000 ns/op for this function on my desktop with Intel® Core™ i7-6700 CPU @ 3.40GHz × 8
running Ubuntu 18.04.1. This is about 5x faster than creating an exception in Java.
*/
func NewStdException(message string) error {
	// Skip the top 4 functions in the stack trace so that the caller of this function is shown at the top
	return newStdException(message, defaultStackTraceDepth, 4)
}

/*
NewStdExceptionWithStackTraceSize creates a new exception. A stack trace is created immediately. stackTraceNumLines allows
you to limit the depth of the stack trace.

The stack trace does not get converted to a string until GetStackTraceAsString is called. Waiting to do this
until it actually gets logged vastly improves performance. I have noticed a performance of about
1800 ns/op to 2000 ns/op for this function on my desktop with Intel® Core™ i7-6700 CPU @ 3.40GHz × 8
running Ubuntu 18.04.1. This is about 5x faster than creating an exception in Java.
*/
func NewStdExceptionWithStackTraceSize(message string, stackTraceNumLines int) error {
	// Skip the top 4 functions in the stack trace so that the caller of this function is shown at the top
	return newStdException(message, stackTraceNumLines, 4)
}

func newStdException(message string, stackTraceNumLines, skip int) *StdException {
	timestamp := time.Now().In(Location)
	return &StdException{
		stackTrace:        getStackTrace(skip, stackTraceNumLines),
		maxStackTraceSize: stackTraceNumLines,
		message:           message,
		timestamp:         &timestamp,
	}
}

/*
prependMsg adds a message to your error:
	timestamp - yourNewMsg
	Caused by:
	Your existing error....
*/
func (se *StdException) prependMsg(msg string) {
	se.messageChain = append([]string{msg}, se.messageChain...)
}

/*
prependMsg adds a message to your error:
	timestamp - yourNewMsg
	Caused by:
	Your existing error....
*/
func PrependMsg(err error, msg string) error {
	if err == nil {
		return nil
	}
	var buf strings.Builder
	buf.WriteString(time.Now().Format(timeFmt))
	buf.WriteString(" - ")
	buf.WriteString(msg)
	msg = buf.String()
	if val, hasPrependFunc := err.(prependable); hasPrependFunc {
		val.prependMsg(msg)
	} else {
		err = fmt.Errorf("%s\nCaused by\n%s", msg, err.Error())
	}

	return err
}

/*
GetStackTrace returns the stack trace as slice of *StackTraceEntry.
*/
func (se *StdException) GetStackTrace() []*StackTraceEntry {
	return se.stackTrace
}

/*
GetStackTraceAsString returns the stack trace in a string formatted as:

	sherlog.exampleFunc(exampleFile.go:18)
	sherlog.exampleFunc2(exampleFile2.go:46)
	sherlog.exampleFunc3(exampleFile2.go:177)

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
Log writes to the writer a string formatted as:

	yyyy-mm-dd hh:mm:ss - message:
		sherlog.exampleFunc(exampleFile.go:18)
		sherlog.exampleFunc2(exampleFile2.go:46)
		sherlog.exampleFunc3(exampleFile2.go:177)

Returns the string that was logged or an error if there was one.
*/
func (se *StdException) Log(writer io.Writer) error {
	err := se.LogNoStack(writer)
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(":\n"))
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(se.GetStackTraceAsString()))
	return err
}

/*
LogNoStack writes to the writer a string formatted as:

	yyyy-mm-dd hh:mm:ss - message

Note that it does not have the stack trace.
Returns the string that was logged or an error if there was one.
*/
func (se *StdException) LogNoStack(writer io.Writer) error {
	for _, msg := range se.messageChain {
		writer.Write([]byte(msg))
		writer.Write([]byte("\nCaused by:\n"))
	}
	_, err := writer.Write([]byte(se.timestamp.Format(timeFmt)))
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(" - "))
	if err != nil {
		return err
	}
	_, err = writer.Write([]byte(se.message))
	return err
}

/*
LogAsJson packages up the exception's info into json and writes it to writer.

The json is formatted like this
	{
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
func (se *StdException) LogAsJson(writer io.Writer) error {
	jsonBytes, err := se.ToJsonBytes()
	if err != nil {
		return err
	}

	_, err = writer.Write(jsonBytes)

	return err
}

/*
Error returns the message and stack trace in a string formatted like this:

	message:
		sherlog.exampleFunc(exampleFile.go:18)
		sherlog.exampleFunc2(exampleFile2.go:46)
		sherlog.exampleFunc3(exampleFile2.go:177)

Leaves out the timestamp so that StdException will print nicely with log.Println
*/
func (se *StdException) Error() string {
	var buf strings.Builder
	buf.WriteString(" - ")
	buf.WriteString(se.message)
	buf.WriteString(":\n")
	buf.WriteString(se.GetStackTraceAsString())
	return buf.String()
}

/*
ToJsonBytes returns the bytes for a json blob that looks like this:

	{
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
func (se *StdException) ToJsonBytes() ([]byte, error) {
	return json.Marshal(se.ToJsonMap())
}

/*
ToJsonMap creates a map[string]interface{} that, when compiled to json, looks like this:

	{
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
func (se *StdException) ToJsonMap() map[string]interface{} {
	return map[string]interface{}{
		"Time":          se.timestamp.Format(timeFmt),
		"Message":       se.message,
		"StackTrace":    se.stackTrace,
		"StackTraceStr": se.GetStackTraceAsString(),
	}
}

func (se *StdException) GetMessage() string {
	return se.message
}