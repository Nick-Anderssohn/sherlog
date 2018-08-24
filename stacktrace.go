package logging

import (
	"runtime"
	"strconv"
	"strings"
)

type StackTraceWrapper interface {
	GetStackTraceAsString() (string, error)
}

type StackTraceEntry struct {
	FunctionName string
	File string
	Line int
}

func (ste *StackTraceEntry) String() (string, error) {
	var buf strings.Builder
	buf.Grow(defaultStackTraceLineLen)
	_, err := buf.WriteString(ste.FunctionName)
	if err != nil {
		return "", err
	}
	_, err = buf.WriteString("(")
	if err != nil {
		return "", err
	}
	_, err = buf.WriteString(ste.File)
	if err != nil {
		return "", err
	}
	_, err = buf.WriteString(":")
	if err != nil {
		return "", err
	}
	_, err = buf.WriteString(strconv.Itoa(ste.Line))
	if err != nil {
		return "", err
	}
	_, err = buf.WriteString(")")
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func CreateStackTraceEntryFromRuntimeFrame(frame *runtime.Frame) *StackTraceEntry {
	return &StackTraceEntry{
		FunctionName: frame.Function,
		File: frame.File,
		Line: frame.Line,
	}
}

func getStackTrace(skip, maxStackTraceSize int) (stackTrace []*StackTraceEntry) {
	programCounters := make([]uintptr, 1)
	runtime.Callers(skip, programCounters)
	framePtr := runtime.CallersFrames(programCounters)
	more := true
	for i := 0; i < maxStackTraceSize; i++ {
		var frame runtime.Frame
		frame, more = framePtr.Next()
		stackTrace = append(stackTrace, CreateStackTraceEntryFromRuntimeFrame(&frame))
		if !more {
			return
		}
	}
	return
}

func StackTraceAsString(stackTrace []*StackTraceEntry) (string, error) {
	var buf strings.Builder
	buf.Grow(defaultStackTraceNumBytes)
	for _, call := range stackTrace {
		_, err := buf.WriteString("\t")
		if err != nil {
			return "", err
		}

		callStr, err := call.String()
		if err != nil {
			return "", err
		}

		_, err = buf.WriteString(callStr)
		if err != nil {
			return "", err
		}

		_, err = buf.WriteString("\n")
		if err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}