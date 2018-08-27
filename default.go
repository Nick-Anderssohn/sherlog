package sherlog

type LevelEnum int

/*
Default log level enums that implement Level.
These are my recommended log levels, but you can create different ones simply by implementing
the Level interface if you would like.
*/
const (
	/**
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
Returns the integer value of the LevelEnum.
*/
func (le LevelEnum) GetLevelId() int {
	return int(le)
}

/*
Returns the text representation of the LevelEnum.
For example, EnumError returns ERROR.
*/
func (le LevelEnum) GetLabel() string {
	return levelLabels[le]
}

/**
Graduates a normal error to a LeveledException with error level CRITICAL.
If err is already a LevelWrapper, then it's level will be changed to CRITICAL without
overriding the stack trace.
*/
func AsCritical(err error) error {
	return errorToLeveledError(err, EnumCritical, 6)
}

/**
Graduates a normal error to a LeveledException with error level ERROR.
If err is already a LevelWrapper, then it's level will be changed to ERROR without
overriding the stack trace.
*/
func AsError(err error) error {
	return errorToLeveledError(err, EnumError, 6)
}

/**
Graduates a normal error to a LeveledException with error level OPS_ERROR.
If err is already a LevelWrapper, then it's level will be changed to OPS_ERROR without
overriding the stack trace.
*/
func AsOpsError(err error) error {
	return errorToLeveledError(err, EnumOpsError, 6)
}

/**
Graduates a normal error to a LeveledException with error level WARNING.
If err is already a LevelWrapper, then it's level will be changed to WARNING without
overriding the stack trace.
*/
func AsWarning(err error) error {
	return errorToLeveledError(err, EnumWarning, 6)
}

/**
Graduates a normal error to a LeveledException with error level INFO.
If err is already a LevelWrapper, then it's level will be changed to INFO without
overriding the stack trace.
*/
func AsInfo(err error) error {
	return errorToLeveledError(err, EnumInfo, 6)
}

/**
Graduates a normal error to a LeveledException with error level DEBUG.
If err is already a LevelWrapper, then it's level will be changed to DEBUG without
overriding the stack trace.
*/
func AsDebug(err error) error {
	return errorToLeveledError(err, EnumDebug, 6)
}

/**
Graduates a normal error to a LeveledException with the specified level.
If err is already a LevelWrapper, then it's level will be changed without creating
a new stack trace.
*/
func errorToLeveledError(err error, level Level, skip int) error {
	if isLevelWrapper(err) {
		err.(LevelWrapper).SetLevel(level)
		return err
	}
	return newLeveledException(err.Error(), level, defaultStackTraceDepth, skip)
}

func NewCritical(message string) error {
	return newLeveledException(message, EnumCritical, defaultStackTraceDepth, 5)
}

func NewError(message string) error {
	return newLeveledException(message, EnumError, defaultStackTraceDepth, 5)
}

func NewOpsError(message string) error {
	return newLeveledException(message, EnumOpsError, defaultStackTraceDepth, 5)
}

func NewWarning(message string) error {
	return newLeveledException(message, EnumWarning, defaultStackTraceDepth, 5)
}

func NewInfo(message string) error {
	return newLeveledException(message, EnumInfo, defaultStackTraceDepth, 5)
}

func NewDebug(message string) error {
	return newLeveledException(message, EnumDebug, defaultStackTraceDepth, 5)
}

/*
Creates a MultiFileLogger setup to use the default Levels that this package provides.
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
