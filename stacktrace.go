package logging

import "runtime"

type StackTraceEntry struct {
	FunctionName string
	File string
	Line int
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