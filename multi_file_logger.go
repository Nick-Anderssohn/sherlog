package sherlog

type MultiLevelLogger interface {
	LogNoStack(loggable LeveledLoggable) error
	LogJson(loggable LeveledLoggable) error
}

type LeveledLoggable interface {
	Loggable
	GetLevel() Level
}

type FileLoggerParams struct {
	LogLevel                  Level
	LogPath                   string
}

/*
Logs to multiple files based off of level. If you provide the same path for multiple
log levels, then they will be logged to that same file with no problems.

Is thread safe :)
 */
type MultiFileLogger struct {
	loggers map[Level]*FileLogger
	defaultLogger *FileLogger // If a Loggable without a log level is provided, this is the logger that will be used
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
		loggers: loggers,
		defaultLogger: defaultLogger,
	}, nil
}

func createFileLoggersFromParams(paths map[Level]string) (map[Level]*FileLogger, error) {
	loggers := map[Level]*FileLogger{}

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

Is thread safe :)
 */
func (mfl *MultiFileLogger) Log(loggable Loggable) error {
	leveledLoggable, isLeveled := loggable.(LeveledLoggable)
	if isLeveled {
		return mfl.loggers[leveledLoggable.GetLevel()].Log(loggable)
	}
	return mfl.defaultLogger.Log(loggable)
}

/*
Logs the error without the stack trace.

Is thread safe :)
 */
func (mfl *MultiFileLogger) LogNoStack(loggable Loggable) error {
	leveledLoggable, isLeveled := loggable.(LeveledLoggable)
	if isLeveled {
		return mfl.loggers[leveledLoggable.GetLevel()].LogNoStack(loggable)
	}
	return mfl.defaultLogger.LogNoStack(loggable)
}

/*
Logs the error as a json blob.

Is thread safe :)
 */
func (mfl *MultiFileLogger) LogJson(loggable Loggable) error {
	leveledLoggable, isLeveled := loggable.(LeveledLoggable)
	if isLeveled {
		return mfl.loggers[leveledLoggable.GetLevel()].LogJson(loggable)
	}
	return mfl.defaultLogger.LogJson(loggable)
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