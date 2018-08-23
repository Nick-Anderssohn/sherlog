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

type logFunction func(writer io.Writer) ([]byte, error)

type Loggable interface {
	LogCompactFmt(writer io.Writer) ([]byte, error)
	LogAsJson(writer io.Writer) ([]byte, error)
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
	notifiable Notifiable
}

func NewLogger(logFilePath string, failureHandler FailureHandler) (*Logger, error) {
	return NewLoggerWithNotifiable(logFilePath, failureHandler, NilNotifiable{})
}

func NewLoggerWithNotifiable(logFilePath string, failureHandler FailureHandler, notifiable Notifiable) (*Logger, error) {
	file, err := os.OpenFile(logFilePath, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{
		logFilePath: logFilePath,
		failureHandler: failureHandler,
		file: file,
		notifiable: notifiable,
	}, nil
}

func (l *Logger) LogCompactFmt(loggable Loggable) {
	l.log(loggable.LogCompactFmt)
}

func (l *Logger) LogJson(loggable Loggable) {
	l.log(loggable.LogAsJson)
}

func (l *Logger) log(logFunc logFunction) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	logMessage, err := logFunc(l.file)
	if err != nil {
		l.failureHandler.HandleFail(err)
	} else {
		go l.notifiable.Notify(logMessage)
	}
}

func (l *Logger) Close() {
	l.file.Close()
}

// TODO: create Logger that handles multiple loggers for different classifications