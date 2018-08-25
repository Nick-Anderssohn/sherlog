package sherlock

import (
	"testing"
)

var testSte = StackTraceEntry{
	FunctionName: "ChickenDinner.testFunc",
	File: "testFile",
	Line: 7,
}

var testStackTrace []*StackTraceEntry

func init() {
	for i := 0; i < defaultStackTraceNumLines; i++ {
		testStackTrace = append(testStackTrace, &testSte)
	}
}

func BenchmarkStackTraceAsString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StackTraceAsString(testStackTrace)
	}
}

func BenchmarkGetStackTrace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getStackTrace(2, defaultStackTraceNumLines)
	}
}

func BenchmarkNewStdException(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewStdException("Test Message")
	}
}