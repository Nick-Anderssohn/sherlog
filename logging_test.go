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

const testMessage = "Test Message"

func init() {
	for i := 0; i < defaultStackTraceNumLines; i++ {
		testStackTrace = append(testStackTrace, &testSte)
	}
}

// ***************** Benchmarks *******************

func BenchmarkStackTraceAsString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stackTraceAsString(testStackTrace)
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

func BenchmarkNewLeveledException(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewLeveledException("Test Message", EnumError)
	}
}