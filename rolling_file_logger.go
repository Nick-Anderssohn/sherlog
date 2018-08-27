package sherlog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

/*
RollingFileLogger is a logger that will automatically start a new log file after a certain amount of time
*/
type RollingFileLogger struct {
	FileLogger
	baseFilePath string
	running      bool
}

/*
NewNightlyRollingFileLogger is a logger that rolls at midnight.
*/
func NewNightlyRollingFileLogger(logFilePath string) (*RollingFileLogger, error) {
	fileLogger, err := NewFileLogger(getTimestampedFileName(logFilePath))
	if err != nil {
		return nil, err
	}
	rollingFileLogger := &RollingFileLogger{
		FileLogger:   *fileLogger,
		baseFilePath: logFilePath,
	}
	go rollingFileLogger.rollNightly()
	return rollingFileLogger, nil
}

/*
NewCustomRollingFileLogger is a logger that rolls every duration. Starts timer upon instantiation
*/
func NewCustomRollingFileLogger(logFilePath string, duration time.Duration) (*RollingFileLogger, error) {
	fileLogger, err := NewFileLogger(getTimestampedFileName(logFilePath))
	if err != nil {
		return nil, err
	}
	rollingFileLogger := &RollingFileLogger{
		FileLogger:   *fileLogger,
		baseFilePath: logFilePath,
	}
	go rollingFileLogger.rollEvery(duration)
	return rollingFileLogger, nil
}

/*
Close closes the file writer.
*/
func (rfl *RollingFileLogger) Close() {
	rfl.running = false
	rfl.FileLogger.Close()
}

func (rfl *RollingFileLogger) rollEvery(duration time.Duration) {
	rfl.running = true
	for rfl.running {
		rfl.rollIn(duration)
	}
}

func (rfl *RollingFileLogger) rollNightly() {
	rfl.running = true
	for rfl.running {
		rfl.rollIn(getDurationUntilTomorrowAtMidnight())
	}
}

func (rfl *RollingFileLogger) rollIn(duration time.Duration) {
	time.Sleep(duration)
	rfl.roll()
}

func (rfl *RollingFileLogger) roll() error {
	rfl.mutex.Lock()
	defer rfl.mutex.Unlock()
	rfl.file.Close()
	rfl.logFilePath = getTimestampedFileName(rfl.baseFilePath)
	newFile, err := openFile(rfl.logFilePath)
	rfl.file = newFile
	return err
}

func getTimestampedFileName(fileName string) string {
	now := time.Now().In(Location)
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

func getDurationUntilTomorrowAtMidnight() time.Duration {
	now := time.Now().In(Location)
	tomorrow := now.AddDate(0, 0, 1)
	tomorrow = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 1, Location) // Tomorrow at midnight
	return tomorrow.Sub(now)
}
