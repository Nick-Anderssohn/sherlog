package logging

type MultiLevelLogger interface {
	LogCompactFmt(loggable LeveledLoggable) error
	LogJson(loggable LeveledLoggable) error
}

type LeveledLoggable interface {
	Loggable
	getLevel() Level
}

type FileLoggerParams struct {
	LogLevel                  Level
	LogPath                   string
	SuccessfullyLoggedHandler Notifiable
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
		if params.SuccessfullyLoggedHandler == nil {
			params.SuccessfullyLoggedHandler = NilNotifiable{}
		}
		logger, err := NewFileLoggerWithNotifiable(params.LogPath, params.SuccessfullyLoggedHandler)
		if err != nil {
			return nil, err
		}
		loggers[params.LogLevel] = logger
	}

	return loggers, nil
}

func (mfl *MultiFileLogger) LogCompactFmt(loggable LeveledLoggable) error {
	return mfl.loggers[loggable.getLevel()].LogCompactFmt(loggable)
}

func (mfl *MultiFileLogger) LogJson(loggable LeveledLoggable) error {
	return mfl.loggers[loggable.getLevel()].LogJson(loggable)
}