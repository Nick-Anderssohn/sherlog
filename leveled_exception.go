package logging

import (
	"io"
	"fmt"
	"encoding/json"
)

type Level interface {
	GetLevelId() int
	GetLabel() string
}

// ************************* LeveledException **************************

// An exception with a level such as ERROR or WARNING
type LeveledException struct {
	StdException
	level Level
}

func (le *LeveledException) GetLevel() Level {
	return le.level
}

func NewLeveledException(message string, level Level) *LeveledException {
	return NewLeveledExceptionWithStackTraceSize(message, level, defaultStackTraceSize)
}

func NewLeveledExceptionWithStackTraceSize(message string, level Level, stackTraceSize int) *LeveledException {
	return &LeveledException{
		StdException: *NewStdExceptionWithStackTraceSize(message, stackTraceSize),
		level:        level,
	}
}

// Writes "timestamp - level - message" to writer.
// Returns returns the logged message or an error if there is one.
func (le *LeveledException) LogCompactFmt(writer io.Writer) ([]byte, error) {
	logMessage := []byte(fmt.Sprintf("%s - %s - %s", le.timestamp.Format(timeFmt), le.level.GetLabel(), le.message))
	_, err := writer.Write(logMessage)
	if err != nil {
		return nil, err
	}
	return logMessage, nil
}

// Packages up the exception's info into json and writes it to writer.
// Returns returns the logged message or an error if there is one.
func (le *LeveledException) LogAsJson(writer io.Writer) (jsonBytes []byte, err error) {
	jsonBytes, err = json.Marshal(map[string]interface{}{
		"Time":       le.timestamp.Format(timeFmt),
		"LevelId":    le.level.GetLevelId(),
		"Level":      le.level.GetLabel(),
		"Message":    le.message,
		"StackTrace": le.stackTrace,
	})

	if err != nil {
		return
	}

	_, err = writer.Write(jsonBytes)
	if err != nil {
		return nil, err
	}

	return
}