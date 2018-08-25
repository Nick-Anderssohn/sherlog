package sherlock

type LevelEnum int

const (
	Critical LevelEnum = iota
	Error
	Warning
	Info
	Debug
)

var levelLabels = map[LevelEnum]string {
	Critical: "CRITICAL",
	Error:    "ERROR",
	Warning:  "WARNING",
	Info:     "INFO",
	Debug:    "DEBUG",
}

func (le LevelEnum) GetLevelId() int {
	return int(le)
}

func (le LevelEnum) GetLabel() string {
	return levelLabels[le]
}

// A MultiFileLogger setup to be instantiated with levels included in this package (see leveled_exception.go)
type StandardMultiFileLogger struct {
	MultiFileLogger
}

func NewStandardMultiFileLogger(paths map[LevelEnum]string) (*StandardMultiFileLogger, error) {
	logger, err := NewMultiFileLogger(createDefaultMultiFileLoggerParams(paths))
	if err != nil {
		return nil, err
	}
	return &StandardMultiFileLogger{
		MultiFileLogger: *logger,
	}, nil
}

func createDefaultMultiFileLoggerParams(paths map[LevelEnum]string) (params []*FileLoggerParams) {
	for level, path := range paths {
		params = append(params, &FileLoggerParams{
			LogLevel: level,
			LogPath: path,
		})
	}
	return
}