package sherlog

import (
	"io"
	"os"
	"sync"
	"time"
	"encoding/json"
)

type logFunction func(writer io.Writer) error

/*
Implement for something to be loggable by a Logger's Log function
 */
type Loggable interface {
	error
	Log(writer io.Writer) error
}

/*
Implement for something to be loggable by a RobustLogger's LogNoStack function
 */
type LoggableWithNoStackOption interface {
	Loggable
	LogNoStack(writer io.Writer) error
}

/*
Implement for something to be loggable by a RobustLogger's LogJson function
 */
type JsonLoggable interface {
	error
	LogAsJson(writer io.Writer) error
}

/*
An interface representing an incredibly basic logger.
*/
type Logger interface {
	Log(errToLog error) error
	Close()
}

/*
An interface representing a Logger that can call all of a Loggable's log functions.
*/
type RobustLogger interface {
	Logger
	LogNoStack(errToLog error) error
	LogJson(errToLog error) error
}

/*
Logs exceptions to a single file path.
Writes are not buffered. Opens and closes per exception written.
*/
type FileLogger struct {
	logFilePath string
	mutex       sync.Mutex
	file        *os.File
}

/*
Create a new logger that will write to logFilePath. Will append to the file if it already exists. Will
create it if it doesn't.
*/
func NewFileLogger(logFilePath string) (*FileLogger, error) {
	file, err := openFile(logFilePath)
	if err != nil {
		return nil, AsError(err)
	}

	return &FileLogger{
		logFilePath: logFilePath,
		file:        file,
	}, nil
}

func openFile(fileName string) (*os.File, error) {
	return os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

/*
Calls loggable's Log function. Is thread safe :)
Non-sherlog errors get logged with only timestamp and message
*/
func (l *FileLogger) Log(errToLog error) error {
	if loggable, isLoggable := errToLog.(Loggable); isLoggable {
		return l.log(loggable.Log)
	}
	return l.logNonSherlogError(errToLog)
}

/*
Calls loggable's LogNoStack function. Is thread safe :)
Non-sherlog errors get logged with only timestamp and message
*/
func (l *FileLogger) LogNoStack(errToLog error) error {
	if loggable, isLoggable := errToLog.(LoggableWithNoStackOption); isLoggable {
		return l.log(loggable.LogNoStack)
	}
	return l.logNonSherlogError(errToLog)
}

/*
Calls loggable's LogJson function. Is thread safe :)
Non-sherlog errors get logged with only timestamp and message
*/
func (l *FileLogger) LogJson(errToLog error) error {
	if loggable, isLoggable := errToLog.(JsonLoggable); isLoggable {
		return l.log(loggable.LogAsJson)
	}

	// Else, manually extract info...
	jsonBytes, err := json.Marshal(map[string]interface{}{
		"Time": time.Now().In(SherlogLocation).Format(timeFmt), // Use log time instead of time of creation since we don't have one....
		"Message": errToLog.Error(),
	})
	if err != nil {
		return err
	}

	l.mutex.Lock()
	_, err = l.file.Write(jsonBytes)
	l.mutex.Unlock()
	return err
}

/*
Closes the file writer.
*/
func (l *FileLogger) Close() {
	l.file.Close()
}

func (l *FileLogger) log(logFunc logFunction) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	err := logFunc(l.file)
	if err != nil {
		return err
	}
	l.file.Write([]byte("\n\n"))
	err = l.file.Sync() // To improve perf, may want to move this to just run every minute or so
	if err != nil {
		return err
	}
	return nil
}

func (l *FileLogger) logNonSherlogError(errToLog error) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	now := time.Now().In(SherlogLocation).Format(timeFmt) // Use log time instead of time of creation since we don't have one....

	_, err := l.file.Write([]byte(now))
	if err != nil {
		return err
	}

	_, err = l.file.Write([]byte(" - "))
	if err != nil {
		return err
	}

	_, err = l.file.Write([]byte(errToLog.Error()))
	return err
}