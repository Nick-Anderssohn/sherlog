package logging

import (
	"io"
	"sync"
	"os"
)

type Notifiable interface {
	Notify(message interface{})
}

type FailureHandler interface {
	HandleFail(err error)
}

type Loggable interface {
	LogCompactFmt(writer io.Writer, failureHandler FailureHandler)
	LogAsJson(writer io.Writer, failureHandler FailureHandler)
}

// A Notifiable that ignores the notification
type NilNotifiable struct{}
func (NilNotifiable) Notify(message interface{}) {}

// Logs exceptions to a single file path
// Writes are not buffered. Opens and closes per exception written
type Logger struct {
	logFilePath string
	failureHandler FailureHandler
	mutex sync.Mutex
	file *os.File
}

func NewLogger(logFilePath string, failureHandler FailureHandler) (*Logger, error) {
	file, err := os.OpenFile(logFilePath, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{
		logFilePath: logFilePath,
		failureHandler: failureHandler,
		file: file,
	}, nil
}

func (l *Logger) LogCompactFmt(loggable Loggable) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	loggable.LogCompactFmt(l.file, l.failureHandler)
}

func (l *Logger) LogJson(loggable Loggable) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	loggable.LogAsJson(l.file, l.failureHandler)
}

// TODO: create Logger that handles multiple loggers for different classifications