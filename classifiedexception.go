package logging

import (
	"io"
	"fmt"
	"encoding/json"
)

type Classification interface {
	GetClassificationId() int
	GetLabel() string
}

// *************** Default Classifications **************************

// User of this package can make his own classifications if he wants, but I provide these
// because they are common classifications

type BasicClassification int

const (
	ClassificationError BasicClassification = iota
	ClassificationWarning
	ClassificationInfo
	ClassificationDebug
	ClassificationCritical
)

var classificationLabels = map[BasicClassification]string {
	ClassificationError: "ERROR",
	ClassificationWarning: "WARNING",
	ClassificationInfo: "INFO",
	ClassificationDebug: "DEBUG",
	ClassificationCritical: "CRITICAL",
}

func (bc BasicClassification) GetClassificationId() int {
	return int(bc)
}

func (bc BasicClassification) GetLabel() string {
	return classificationLabels[bc]
}

// ************************* ClassifiedException **************************

// An exception with a classification such as ERROR or WARNING
type ClassifiedException struct {
	StdException
	classification Classification
}

func NewClassifiedException(message string, classification Classification) *ClassifiedException {
	return NewClassifiedExceptionWithStackTrace(message, classification, defaultStackTraceSize)
}

func NewClassifiedExceptionWithStackTrace(message string, classification Classification, stackTraceSize int) *ClassifiedException {
	return &ClassifiedException{
		StdException: *NewStdExceptionWithStackTraceSize(message, stackTraceSize),
		classification: classification,
	}
}

// Writes "timestamp - classification - message" to writer.
// Returns returns the logged message or an error if there is one.
func (ce *ClassifiedException) LogCompactFmt(writer io.Writer) ([]byte, error) {
	logMessage := []byte(fmt.Sprintf("%s - %s - %s", ce.timestamp.Format(timeFmt), ce.classification.GetLabel(), ce.message))
	_, err := writer.Write(logMessage)
	if err != nil {
		return nil, err
	}
	return logMessage, nil
}

// Packages up the exception's info into json and writes it to writer.
// Returns returns the logged message or an error if there is one.
func (ce *ClassifiedException) LogAsJson(writer io.Writer) (jsonBytes []byte, err error) {
	jsonBytes, err = json.Marshal(map[string]interface{}{
		"Time": ce.timestamp.Format(timeFmt),
		"ClassificationId": ce.classification.GetClassificationId(),
		"Classification": ce.classification.GetLabel(),
		"Message": ce.message,
		"StackTrace": ce.stackTrace,
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