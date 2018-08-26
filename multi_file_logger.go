package sherlock

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

// Logs to multiple files based off of level
type MultiFileLogger struct {
	loggers map[Level]*FileLogger
}

func NewMultiFileLogger(fileLoggerParams []*FileLoggerParams) (*MultiFileLogger, error) {
	loggers, err := createFileLoggersFromParams(fileLoggerParams)
	if err != nil {
		return nil, err
	}

	return &MultiFileLogger{
		loggers: loggers,
	}, nil
}

func createFileLoggersFromParams(fileLoggerParams []*FileLoggerParams) (map[Level]*FileLogger, error) {
	loggers := map[Level]*FileLogger{}

	for _, params := range fileLoggerParams {
		logger, err := NewFileLogger(params.LogPath)
		if err != nil {
			return nil, err
		}
		loggers[params.LogLevel] = logger
	}

	return loggers, nil
}

func (mfl *MultiFileLogger) LogNoStack(loggable LeveledLoggable) error {
	return mfl.loggers[loggable.GetLevel()].LogNoStack(loggable)
}

func (mfl *MultiFileLogger) LogJson(loggable LeveledLoggable) error {
	return mfl.loggers[loggable.GetLevel()].LogJson(loggable)
}

func (mfl *MultiFileLogger) ErrorIsLoggable(err error) bool {
	_, isLoggable := err.(LeveledLoggable)
	return isLoggable
}