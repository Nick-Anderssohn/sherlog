package sherlog

import (
	"testing"
)

var testSte = StackTraceEntry{
	FunctionName: "ChickenDinner.testFunc",
	File:         "testFile",
	Line:         7,
}

var testStackTrace []*StackTraceEntry

const testMessage = "Test Message"

func init() {
	for i := 0; i < defaultStackTraceDepth; i++ {
		testStackTrace = append(testStackTrace, &testSte)
	}
}

// ***************** Tests ************************

func TestLeveledExceptionImplementsDesiredInterfaces(t *testing.T) {
	// error, LevelWrapper, Loggable, StackTraceWrapper, and LeveledLoggable.
	var exception interface{} = NewLeveledException("Wub Wub", EnumInfo)
	_, implements := exception.(error)
	errorIfFalse(implements, t, "not an error")
	_, implements = exception.(LevelWrapper)
	errorIfFalse(implements, t, "not a LevelWrapper")
	_, implements = exception.(Loggable)
	errorIfFalse(implements, t, "not a Loggable")
	_, implements = exception.(LeveledLoggable)
	errorIfFalse(implements, t, "not a LeveledLoggable")
}

func TestStdExceptionImplementsDesiredInterfaces(t *testing.T) {
	// error, Loggable, StackTraceWrapper.
	var exception interface{} = NewStdException("Wub Wub")
	_, implements := exception.(error)
	errorIfFalse(implements, t, "not an error")
	_, implements = exception.(Loggable)
	errorIfFalse(implements, t, "not a Loggable")
}

func TestImplementsLogger(t *testing.T) {
	var fileLogger interface{} = &FileLogger{}
	//_, implementsLogger := fileLogger.(Logger)
	//errorIfFalse(implementsLogger, t, "FileLogger does not implement Logger")

	var multiFileLogger interface{} = &MultiFileLogger{}
	//_, implementsLogger = multiFileLogger.(Logger)
	//errorIfFalse(implementsLogger, t, "MultiFileLogger does not implement Logger")

	var rollingFileLogger interface{} = &SizeBasedRollingFileLogger{}
	//_, implementsLogger = rollingFileLogger.(Logger)
	//errorIfFalse(implementsLogger, t, "SizeBasedRollingFileLogger does not implement Logger")

	_, implementsLogger := fileLogger.(Logger)
	errorIfFalse(implementsLogger, t, "FileLogger does not implement Logger")

	_, implementsLogger = multiFileLogger.(Logger)
	errorIfFalse(implementsLogger, t, "MultiFileLogger does not implement Logger")

	_, implementsLogger = rollingFileLogger.(Logger)
	errorIfFalse(implementsLogger, t, "SizeBasedRollingFileLogger does not implement Logger")
}

func errorIfFalse(val bool, t *testing.T, failMessage string) {
	if !val {
		t.Error(failMessage)
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
		getStackTrace(2, defaultStackTraceDepth)
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
