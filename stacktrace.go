package sherlock

import (
	"runtime"
	"strconv"
	"strings"
)

type StackTraceWrapper interface {
	GetStackTrace() []*StackTraceEntry
	GetStackTraceAsString() string
}

type StackTraceEntry struct {
	FunctionName string
	File string
	Line int
}

func (ste *StackTraceEntry) String() string {
	var buf strings.Builder
	buf.Grow(defaultStackTraceLineLen)
	buf.WriteString(ste.FunctionName)
	buf.WriteString("(")
	buf.WriteString(ste.File)
	buf.WriteString(":")
	buf.WriteString(strconv.Itoa(ste.Line))
	buf.WriteString(")")
	return buf.String()
}

func createStackTraceEntryFromRuntimeFrame(frame *runtime.Frame) *StackTraceEntry {
	return &StackTraceEntry{
		FunctionName: frame.Function,
		File: frame.File,
		Line: frame.Line,
	}
}

/*
skip is the number of calls to skip recording at the top of our stack trace
maxStackSize limits the number of callers to record in the stack trace
 */
func getStackTrace(skip, maxStackTraceSize int) (stackTrace []*StackTraceEntry) {
	programCounters := make([]uintptr, maxStackTraceSize)
	runtime.Callers(skip, programCounters)
	framePtr := runtime.CallersFrames(programCounters)

	for i, more := 0, true; i < maxStackTraceSize && more; i++ {
		var frame runtime.Frame
		frame, more = framePtr.Next()

		if frame.Function == "" {
			return
		}

		stackTrace = append(stackTrace, createStackTraceEntryFromRuntimeFrame(&frame))
	}
	return
}

/*
Returns the stack trace in the following format:
	sherlock.exampleFunc(exampleFile.go:18)
	sherlock.exampleFunc2(exampleFile2.go:46)
	sherlock.exampleFunc3(exampleFile2.go:177)
*/
func stackTraceAsString(stackTrace []*StackTraceEntry) string {
	var buf strings.Builder
	buf.Grow(defaultStackTraceNumBytes)
	for _, call := range stackTrace {
		buf.WriteString("\t")
		buf.WriteString(call.String())
		buf.WriteString("\n")
	}
	return buf.String()
}