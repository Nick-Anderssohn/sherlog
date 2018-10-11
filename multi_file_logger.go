package sherlog

import "time"

/*
LeveledLoggable is a Loggable that also has a log level attached to it.
*/
type LeveledLoggable interface {
	Loggable
	GetLevel() Level
}

/*
MultiFileLogger logs to multiple files based off of level. If you provide the same path for multiple
log levels, then they will be logged to that same file with no problems.

Is thread safe :)
*/
type MultiFileLogger struct {
	loggers       map[Level]Logger
	defaultLogger *FileLogger // If a Loggable without a log level is provided, this is the logger that will be used
}

/*
NewMultiFileLoggerRollOnDuration returns a new MultiFileLogger. Logs will roll every duration.
*/
func NewMultiFileLoggerRollOnDuration(paths map[Level]string, defaultLogPath string, duration time.Duration) (*MultiFileLogger, error) {
	loggers, err := createRollingFileLoggersCustomDuration(paths, duration)
	if err != nil {
		return nil, err
	}
	defaultLogger, err := NewFileLogger(defaultLogPath)
	if err != nil {
		return nil, err
	}
	return &MultiFileLogger{
		loggers:       loggers,
		defaultLogger: defaultLogger,
	}, nil
}

/*
NewMultiFileLoggerRoleNightly returns a new MultiFileLogger. Logs will roll daily (at midnight).
*/
func NewMultiFileLoggerRoleNightly(paths map[Level]string, defaultLogPath string) (*MultiFileLogger, error) {
	loggers, err := createNightlyRollingFileLogger(paths)
	if err != nil {
		return nil, err
	}
	defaultLogger, err := NewFileLogger(defaultLogPath)
	if err != nil {
		return nil, err
	}
	return &MultiFileLogger{
		loggers:       loggers,
		defaultLogger: defaultLogger,
	}, nil
}

/*
NewMultiFileLoggerWithSizeBaseRollingLogs returns a new MultiFileLogger. Logs will roll when they maxLogMessagesPerLogFile
*/
func NewMultiFileLoggerWithSizeBaseRollingLogs(paths map[Level]string, defaultLogPath string, maxLogMessagesPerLogFile int) (*MultiFileLogger, error) {
	loggers, err := createSizedBasedRollingFileLoggers(paths, maxLogMessagesPerLogFile)
	if err != nil {
		return nil, err
	}
	defaultLogger, err := NewFileLogger(defaultLogPath)
	if err != nil {
		return nil, err
	}
	return &MultiFileLogger{
		loggers:       loggers,
		defaultLogger: defaultLogger,
	}, nil
}

/*
NewMultiFileLogger returns a new file logger. Will log to different files based off of the paths given for each
log level. If you want some log levels to be logged to the same file, just pass in the same path
for those levels. defaultLogPath is the file to log to if a Loggable is provided that does not have a level.
*/
func NewMultiFileLogger(paths map[Level]string, defaultLogPath string) (*MultiFileLogger, error) {
	loggers, err := createFileLoggers(paths)
	if err != nil {
		return nil, err
	}
	defaultLogger, err := NewFileLogger(defaultLogPath)
	if err != nil {
		return nil, err
	}
	return &MultiFileLogger{
		loggers:       loggers,
		defaultLogger: defaultLogger,
	}, nil
}

// Creates loggers for the various levels. Any levels that share the same path will use the same logger.
func createRobustLoggers(paths map[Level]string, loggerConstructor func(path string) (Logger, error)) (loggers map[Level]Logger, err error) {
	loggers = map[Level]Logger{}
	cachedLoggers := map[string]Logger{}

	for logLevel, path := range paths {
		// Use existing logger if one exists for the path
		logger := cachedLoggers[path]
		if logger == nil {
			logger, err = loggerConstructor(path)
			if err != nil {
				return
			}
			cachedLoggers[path] = logger
		}
		loggers[logLevel] = logger
	}

	return
}

// *************** These functions leverage the createRobustLoggers function to instantiate the needed loggers *************

func createRollingFileLoggersCustomDuration(paths map[Level]string, duration time.Duration) (map[Level]Logger, error) {
	constructLogger := func(loggerPath string) (Logger, error) {
		return NewCustomRollingFileLogger(loggerPath, duration)
	}

	return createRobustLoggers(paths, constructLogger)
}

func createNightlyRollingFileLogger(paths map[Level]string) (map[Level]Logger, error) {
	constructLogger := func(loggerPath string) (Logger, error) {
		return NewNightlyRollingFileLogger(loggerPath)
	}
	return createRobustLoggers(paths, constructLogger)
}

func createSizedBasedRollingFileLoggers(paths map[Level]string, maxLogMessagesPerLogFile int) (map[Level]Logger, error) {
	constructLogger := func(loggerPath string) (Logger, error) {
		return NewRollingFileLoggerWithSizeLimit(loggerPath, maxLogMessagesPerLogFile)
	}
	return createRobustLoggers(paths, constructLogger)
}

func createFileLoggers(paths map[Level]string) (map[Level]Logger, error) {
	constructLogger := func(loggerPath string) (Logger, error) {
		return NewFileLogger(loggerPath)
	}
	return createRobustLoggers(paths, constructLogger)
}

// *************************************************************************************************************************

/*
Log logs the error.
If not a sherlog error, will just be logged with a timestamp and message.

Is thread safe :)
*/
func (mfl *MultiFileLogger) Log(errToLog error) error {
	if leveledLoggable, isLeveled := errToLog.(LeveledLoggable); isLeveled {
		logger := mfl.loggers[leveledLoggable.GetLevel()]
		if logger != nil {
			return logger.Log(errToLog)
		}
	}
	return mfl.defaultLogger.Log(errToLog)
}

/*
LogNoStack logs the error without the stack trace.

Is thread safe :)
*/
func (mfl *MultiFileLogger) LogNoStack(errToLog error) error {
	if leveledLoggable, isLeveled := errToLog.(LeveledLoggable); isLeveled {
		logger := mfl.loggers[leveledLoggable.GetLevel()]
		if logger != nil {
			return logger.LogNoStack(errToLog)
		}
	}
	return mfl.defaultLogger.LogNoStack(errToLog)
}

/*
LogJson logs the error as a json blob.
If not a sherlog error, will just include message.

Is thread safe :)
*/
func (mfl *MultiFileLogger) LogJson(errToLog error) error {
	if leveledLoggable, isLeveled := errToLog.(LeveledLoggable); isLeveled {
		logger := mfl.loggers[leveledLoggable.GetLevel()]
		if logger != nil {
			return logger.LogJson(errToLog)
		}
	}
	return mfl.defaultLogger.LogJson(errToLog)
}

/*
Close closes all loggers.
*/
func (mfl *MultiFileLogger) Close() {
	for _, logger := range mfl.loggers {
		logger.Close()
	}
	mfl.defaultLogger.Close()
}

/*
ErrorIsLoggable checks if an error is loggable by MultiFileLogger
*/
func (mfl *MultiFileLogger) ErrorIsLoggable(err error) bool {
	_, isLoggable := err.(Loggable)
	return isLoggable
}

/*
Critical turns values into a *LeveledException with level CRITICAL and then calls the logger's
Log function.
*/
func (mfl *MultiFileLogger) Critical(values ...interface{}) error {
	return mfl.Log(graduateOrConcatAndCreate(EnumCritical, values...))
}

/*
Error turns values into a *LeveledException with level ERROR and then calls the logger's
Log function.
*/
func (mfl *MultiFileLogger) Error(values ...interface{}) error {
	return mfl.Log(graduateOrConcatAndCreate(EnumError, values...))
}

/*
OpsError turns values into a *LeveledException with level OPS_ERROR and then calls the logger's
Log function.
*/
func (mfl *MultiFileLogger) OpsError(values ...interface{}) error {
	return mfl.Log(graduateOrConcatAndCreate(EnumOpsError, values...))
}

/*
Warn turns values into a *LeveledException with level WARNING and then calls the logger's
Log function.
*/
func (mfl *MultiFileLogger) Warn(values ...interface{}) error {
	return mfl.Log(graduateOrConcatAndCreate(EnumWarning, values...))
}

/*
Info turns values into a *LeveledException with level INFO and then calls the logger's
Log function.
*/
func (mfl *MultiFileLogger) Info(values ...interface{}) error {
	return mfl.Log(graduateOrConcatAndCreate(EnumInfo, values...))
}

/*
Debug turns values into a *LeveledException with level DEBUG and then calls the logger's
Log function.
*/
func (mfl *MultiFileLogger) Debug(values ...interface{}) error {
	return mfl.Log(graduateOrConcatAndCreate(EnumDebug, values...))
}
