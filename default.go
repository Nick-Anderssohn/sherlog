package sherlog

import (
	"fmt"
)

/*
LevelEnum is the default enum sherlog offers that implements Level
*/
type LevelEnum int

/*
Default log level enums that implement Level.
These are my recommended log levels, but you can create different ones simply by implementing
the Level interface if you would like.
*/
const (
	/*
		CRITICAL is the intended log level for panics that are caught in the recover function.
	*/
	EnumCritical LevelEnum = iota

	/*
		ERROR is the intended log level for something that should never ever happen and for sure
		means there is a bug in your code.
	*/
	EnumError

	/*
		OPS_ERROR is the intended log level for an error that is known to be possible due to an operations issue.
		For example, failing to query a database should be an OpsError because it lets you know that your database is
		offline. It doesn't mean there is a bug in your code, but it is still something that needs to be fixed asap.
	*/
	EnumOpsError

	/*
		WARNING is the intended log level for something that doesn't cause failure, but may be frowned
		upon anyways. For example, use of a deprecated endpoint may be logged as a warning. A warning should say,
		"Hey, we don't want to be doing this. It works, but it is bad."
	*/
	EnumWarning

	/*
		INFO is the intended log level for something that you want logged purely to collect information or metrics.
	*/
	EnumInfo

	/*
		DEBUG is for any debug messages you want logged. Ideally, you are not logging these in production.
	*/
	EnumDebug
)

var levelLabels = map[LevelEnum]string{
	EnumCritical: "CRITICAL",
	EnumError:    "ERROR",
	EnumOpsError: "OPS_ERROR",
	EnumWarning:  "WARNING",
	EnumInfo:     "INFO",
	EnumDebug:    "DEBUG",
}

/*
GetLevelId returns the integer value of the LevelEnum.
*/
func (le LevelEnum) GetLevelId() int {
	return int(le)
}

/*
GetLabel returns the text representation of the LevelEnum.
For example, EnumError returns ERROR.
*/
func (le LevelEnum) GetLabel() string {
	return levelLabels[le]
}

/*
AsCritical graduates a normal error to a LeveledException with error level CRITICAL.
If err is already a LevelWrapper, then it's level will be changed to CRITICAL without
overriding the stack trace. As of 1.7.0: if multiple values are passed in, then they will
be concatenated before returning the error.
*/
func AsCritical(values ...interface{}) *LeveledException {
	return graduateOrConcatAndCreate(EnumCritical, values...)
}

/*
AsError graduates a normal error to a LeveledException with error level ERROR.
If err is already a LevelWrapper, then it's level will be changed to ERROR without
overriding the stack trace. As of 1.7.0: if multiple values are passed in, then they will
be concatenated before returning the error.
*/
func AsError(values ...interface{}) *LeveledException {
	return graduateOrConcatAndCreate(EnumError, values...)
}

/*
AsOpsError graduates a normal error to a LeveledException with error level OPS_ERROR.
If err is already a LevelWrapper, then it's level will be changed to OPS_ERROR without
overriding the stack trace. As of 1.7.0: if multiple values are passed in, then they will
be concatenated before returning the error.
*/
func AsOpsError(values ...interface{}) *LeveledException {
	return graduateOrConcatAndCreate(EnumOpsError, values...)
}

/*
AsWarning graduates a normal error to a LeveledException with error level WARNING.
If err is already a LevelWrapper, then it's level will be changed to WARNING without
overriding the stack trace. As of 1.7.0: if multiple values are passed in, then they will
be concatenated before returning the error.
*/
func AsWarning(values ...interface{}) *LeveledException {
	return graduateOrConcatAndCreate(EnumWarning, values...)
}

/*
AsInfo graduates a normal error to a LeveledException with error level INFO.
If err is already a LevelWrapper, then it's level will be changed to INFO without
overriding the stack trace. As of 1.7.0: if multiple values are passed in, then they will
be concatenated before returning the error.
*/
func AsInfo(values ...interface{}) *LeveledException {
	return graduateOrConcatAndCreate(EnumInfo, values...)
}

/*
AsDebug graduates a normal error to a LeveledException with error level DEBUG.
If err is already a LevelWrapper, then it's level will be changed to DEBUG without
overriding the stack trace. As of 1.7.0: if multiple values are passed in, then they will
be concatenated before returning the error.
*/
func AsDebug(values ...interface{}) *LeveledException {
	return graduateOrConcatAndCreate(EnumDebug, values...)
}

/*
graduateOrConcatAndCreate was added in 1.7.0 to extend the behavior of the AsFoo functions
so that they can accept multiple arguments.
*/
func graduateOrConcatAndCreate(level Level, values ...interface{}) *LeveledException {
	// If values simply contains one err, maintain behavior from 1.6.2
	if len(values) == 1 {
		err, ok := values[0].(error)
		if ok {
			return errorToLeveledError(err, level, 7)
		}
		if err == nil {
			return nil
		}
	}

	// ^1.7.0 will concatenate values into an error
	return errorToLeveledError(fmt.Errorf(fmt.Sprint(values...)), level, 7)
}

/*
errorToLeveledError graduates a normal error to a LeveledException with the specified level.
If err is already a *LeveledException, then it's level will be changed without creating
a new stack trace.
*/
func errorToLeveledError(err error, level Level, skip int) *LeveledException {
	if err == nil {
		return nil
	}
	leveledException, ok := err.(*LeveledException)
	if ok {
		leveledException.SetLevel(level)
		return leveledException
	}
	return newLeveledException(err.Error(), level, defaultStackTraceDepth, skip)
}

/*
NewCritical returns a new LeveledException with the level set to CRITICAL.
*/
func NewCritical(message string) *LeveledException {
	return newLeveledException(message, EnumCritical, defaultStackTraceDepth, 5)
}

/*
NewError returns a new LeveledException with the level set to ERROR.
*/
func NewError(message string) *LeveledException {
	return newLeveledException(message, EnumError, defaultStackTraceDepth, 5)
}

/*
NewOpsError returns a new LeveledException with the level set to OPS_ERROR.
*/
func NewOpsError(message string) *LeveledException {
	return newLeveledException(message, EnumOpsError, defaultStackTraceDepth, 5)
}

/*
NewWarning returns a new LeveledException with the level set to WARNING.
*/
func NewWarning(message string) *LeveledException {
	return newLeveledException(message, EnumWarning, defaultStackTraceDepth, 5)
}

/*
NewInfo returns a new LeveledException with the level set to INFO.
*/
func NewInfo(message string) *LeveledException {
	return newLeveledException(message, EnumInfo, defaultStackTraceDepth, 5)
}

/*
NewDebug returns a new LeveledException with the level set to DEBUG.
*/
func NewDebug(message string) *LeveledException {
	return newLeveledException(message, EnumDebug, defaultStackTraceDepth, 5)
}

/*
CreateDefaultMultiFileLogger creates a MultiFileLogger setup to use the default Levels that this package provides.
*/
func CreateDefaultMultiFileLogger(criticalPath, errorPath, warningPath, infoPath, debugPath, defaultPath string) (*MultiFileLogger, error) {
	paths := map[Level]string{
		EnumCritical: criticalPath,
		EnumError:    errorPath,
		EnumWarning:  warningPath,
		EnumInfo:     infoPath,
		EnumDebug:    debugPath,
	}

	return NewMultiFileLogger(paths, defaultPath)
}
