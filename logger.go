package sherlog

import (
	"io"
	"sync"
	"os"
)

type logFunction func(writer io.Writer) (error)

type Loggable interface {
	Log(writer io.Writer) error
	LogNoStack(writer io.Writer) error
	LogAsJson(writer io.Writer) error
}

type Logger interface {
	Log(loggable Loggable) error
	Close()
}

// Logs exceptions to a single file path
// Writes are not buffered. Opens and closes per exception written
type FileLogger struct {
	logFilePath               string
	mutex                     sync.Mutex
	file                      *os.File
}

/*
Create a new logger that will write to logFilePath. Will append to the file if it already exists. Will
create it if it doesn't.
 */
func NewFileLogger(logFilePath string ) (*FileLogger, error) {
	file, err := os.OpenFile(logFilePath, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &FileLogger{
		logFilePath: logFilePath,
		file: file,
	}, nil
}

/*
Calls loggable's Log function. Is thread safe :)
 */
func (l *FileLogger) Log(loggable Loggable) error {
	return l.log(loggable.Log)
}

/*
Calls loggable's LogNoStack function. Is thread safe :)
 */
func (l *FileLogger) LogNoStack(loggable Loggable) error {
	return l.log(loggable.LogNoStack)
}

/*
Calls loggable's LogJson function. Is thread safe :)
 */
func (l *FileLogger) LogJson(loggable Loggable) error {
	return l.log(loggable.LogAsJson)
}

/*
Closes the file writer.
*/
func (l *FileLogger) Close() {
	l.file.Close()
}

/*
Checks if an error is loggable by MultiFileLogger

Is thread safe :)
 */
func (l *FileLogger) ErrorIsLoggable(err error) bool {
	_, isLoggable := err.(Loggable)
	return isLoggable
}

func (l *FileLogger) log(logFunc logFunction) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	err := logFunc(l.file)
	if err != nil {
		return err
	}
	err = l.file.Sync() // To improve perf, may want to move this to just run every minute or so
	if err != nil {
		return err
	}
	return nil
}

