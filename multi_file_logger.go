package sherlog

type LeveledLoggable interface {
	Loggable
	GetLevel() Level
}

/*
Logs to multiple files based off of level. If you provide the same path for multiple
log levels, then they will be logged to that same file with no problems.

Is thread safe :)
*/
type MultiFileLogger struct {
	loggers       map[Level]RobustLogger
	defaultLogger *FileLogger // If a Loggable without a log level is provided, this is the logger that will be used
}

func NewMultiFileLoggerWithRollingLogs(paths map[Level]string, defaultLogPath string, maxLogMessagesPerLogFile int) (*MultiFileLogger, error) {
	loggers, err := createRollingFileLoggersFromParams(paths, maxLogMessagesPerLogFile)
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

func createRollingFileLoggersFromParams(paths map[Level]string, maxLogMessagesPerLogFile int) (map[Level]RobustLogger, error) {
	loggers := map[Level]RobustLogger{}

	for logLevel, path := range paths {
		logger, err := NewRollingFileLogger(path, maxLogMessagesPerLogFile)
		if err != nil {
			return nil, err
		}
		loggers[logLevel] = logger
	}

	return loggers, nil
}

/*
Returns a new file logger. Will log to different files based off of the paths given for each
log level. If you want some log levels to be logged to the same file, just pass in the same path
for those levels. defaultLogPath is the file to log to if a Loggable is provided that does not have a level.
*/
func NewMultiFileLogger(paths map[Level]string, defaultLogPath string) (*MultiFileLogger, error) {
	loggers, err := createFileLoggersFromParams(paths)
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

func createFileLoggersFromParams(paths map[Level]string) (map[Level]RobustLogger, error) {
	loggers := map[Level]RobustLogger{}

	for logLevel, path := range paths {
		logger, err := NewFileLogger(path)
		if err != nil {
			return nil, err
		}
		loggers[logLevel] = logger
	}

	return loggers, nil
}

/*
Logs the error.
If not a sherlog error, will just be logged with a timestamp and message.

Is thread safe :)
*/
func (mfl *MultiFileLogger) Log(errToLog error) error {
	if leveledLoggable, isLeveled := errToLog.(LeveledLoggable); isLeveled {
		return mfl.loggers[leveledLoggable.GetLevel()].Log(errToLog)
	}
	return mfl.defaultLogger.Log(errToLog)
}

/*
Logs the error without the stack trace.

Is thread safe :)
*/
func (mfl *MultiFileLogger) LogNoStack(errToLog error) error {
	if leveledLoggable, isLeveled := errToLog.(LeveledLoggable); isLeveled {
		return mfl.loggers[leveledLoggable.GetLevel()].LogNoStack(errToLog)
	}
	return mfl.defaultLogger.LogNoStack(errToLog)
}

/*
Logs the error as a json blob.
If not a sherlog error, will just include message.

Is thread safe :)
*/
func (mfl *MultiFileLogger) LogJson(errToLog error) error {
	if leveledLoggable, isLeveled := errToLog.(LeveledLoggable); isLeveled {
		return mfl.loggers[leveledLoggable.GetLevel()].LogJson(errToLog)
	}
	return mfl.defaultLogger.LogJson(errToLog)
}

/*
Closes all loggers.
*/
func (mfl *MultiFileLogger) Close() {
	for _, logger := range mfl.loggers {
		logger.Close()
	}
	mfl.defaultLogger.Close()
}

/*
Checks if an error is loggable by MultiFileLogger

Is thread safe :)
*/
func (mfl *MultiFileLogger) ErrorIsLoggable(err error) bool {
	_, isLoggable := err.(Loggable)
	return isLoggable
}
