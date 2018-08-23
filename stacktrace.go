package logging

import "runtime"

func getStackTrace(skip, maxStackTraceSize int) (stackTrace []*runtime.Frame) {
	programCounters := make([]uintptr, 1)
	runtime.Callers(skip, programCounters)
	framePtr := runtime.CallersFrames(programCounters)
	more := true
	for i := 0; i < maxStackTraceSize; i++ {
		var frame runtime.Frame
		frame, more = framePtr.Next()
		stackTrace = append(stackTrace, &frame)
		if !more {
			return
		}
	}
	return
}