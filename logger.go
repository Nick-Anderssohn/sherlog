package logging

import (
	"io"
	"sync"
	"os"
)

type Notifiable interface {
	Notify(message interface{})
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
type FileLogger struct {
	logFilePath               string
	mutex                     sync.Mutex
	file                      *os.File
	successfullyLoggedHandler Notifiable
}

func NewFileLogger(logFilePath string )(*FileLogger, error) {
	return NewFileLoggerWithNotifiable(logFilePath, NilNotifiable{})
}

func NewFileLoggerWithNotifiable(logFilePath string, successfullyLoggedNotifiable Notifiable) (*FileLogger, error) {
	file, err := os.OpenFile(logFilePath, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &FileLogger{
		logFilePath: logFilePath,
		file: file,
		successfullyLoggedHandler: successfullyLoggedNotifiable,
	}, nil
}

func (l *FileLogger) LogCompactFmt(loggable Loggable) error {
	return l.log(loggable.LogCompactFmt)
}

func (l *FileLogger) LogJson(loggable Loggable) error {
	return l.log(loggable.LogAsJson)
}

func (l *FileLogger) log(logFunc logFunction) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	logMessage, err := logFunc(l.file)
	if err != nil {
		return err
	}
	go l.successfullyLoggedHandler.Notify(logMessage)
	return nil
}

func (l *FileLogger) Close() {
	l.file.Close()
}