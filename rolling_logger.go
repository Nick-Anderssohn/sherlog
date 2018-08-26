package sherlog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RollingFileLogger struct {
	FileLogger
	baseFilePath  string
	countToRollOn int
	curCount      int
}

func NewRollingFileLogger(logFilePath string, numMessagesPerFile int) (*RollingFileLogger, error) {
	if numMessagesPerFile <= 0 {
		return nil, NewLeveledException("log files must have room for at least 1 message.", EnumError)
	}
	fileLogger, err := NewFileLogger(logFilePath)
	if err != nil {
		return nil, err
	}
	return &RollingFileLogger{
		FileLogger:    *fileLogger,
		baseFilePath:  logFilePath,
		countToRollOn: numMessagesPerFile,
	}, nil
}

/*
Calls loggable's Log function. Is thread safe :)
*/
func (rfl *RollingFileLogger) Log(loggable Loggable) error {
	err := rfl.log(loggable.Log)
	if err != nil {
		return err
	}
	return rfl.incAndRollIfNecessary()
}

/*
Calls loggable's LogNoStack function. Is thread safe :)
*/
func (rfl *RollingFileLogger) LogNoStack(loggable Loggable) error {
	err := rfl.log(loggable.LogNoStack)
	if err != nil {
		return err
	}
	return rfl.incAndRollIfNecessary()
}

/*
Calls loggable's LogJson function. Is thread safe :)
*/
func (rfl *RollingFileLogger) LogJson(loggable Loggable) error {
	err := rfl.log(loggable.LogAsJson)
	if err != nil {
		return err
	}
	return rfl.incAndRollIfNecessary()
}

func (rfl *RollingFileLogger) incAndRollIfNecessary() error {
	rfl.curCount++
	if rfl.curCount >= rfl.countToRollOn {
		rfl.Close()
		rfl.logFilePath = getTimestampedFileName(rfl.baseFilePath)
		newFile, err := openFile(rfl.logFilePath)
		rfl.file = newFile
		return err
	}
	return nil
}

func getTimestampedFileName(fileName string) string {
	now := time.Now().UTC()
	ext := filepath.Ext(fileName)
	fileName = fileName[:len(fileName)-len(ext)] + now.Format(timeFileNameFmt) + ext
	return incFileNameUntilNotExists(fileName)
}

func incFileNameUntilNotExists(fileName string) string {
	for fileExists(fileName) {
		fileName = incFileName(fileName)
	}
	return fileName
}

func fileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

// Assumes that a file that has ")" right before the extension needs the number inside incremented
func incFileName(fileName string) string {
	ext := filepath.Ext(fileName)
	fileName = fileName[:len(fileName)-len(ext)]
	var fileVersion int
	if fileName[len(fileName)-1] == ')' {
		firstParenIndex := strings.Index(fileName, "(")
		fmt.Fscanf(strings.NewReader(fileName), fileName[:firstParenIndex]+"(%d)", &fileVersion)
		fileName = fileName[:firstParenIndex]
	}
	fileVersion++

	return fmt.Sprintf(fileName+"(%d)"+ext, fileVersion)
}
