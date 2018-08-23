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
// On failure, it will pass an error to failureHandler.
func (ce *ClassifiedException) LogCompactFmt(writer io.Writer, failureHandler FailureHandler) {
	_, err := writer.Write([]byte(fmt.Sprintf("%s - %s - %s", ce.timestamp.Format(timeFmt), ce.classification.GetLabel(), ce.message)))

	if err != nil {
		failureHandler.HandleFail(err)
	}
}

// Packages up the exception's info into json and writes it to writer.
// On failure, it will pass an error to failureHandler.
func (ce *ClassifiedException) LogAsJson(writer io.Writer, failureHandler FailureHandler) {
	jsonBytes, err := json.Marshal(map[string]interface{}{
		"Time": ce.timestamp.Format(timeFmt),
		"ClassificationId": ce.classification.GetClassificationId(),
		"Classification": ce.classification.GetLabel(),
		"Message": ce.message,
		"StackTrace": ce.stackTrace,
	})

	if err != nil {
		failureHandler.HandleFail(err)
		return
	}

	_, err = writer.Write(jsonBytes)

	if err != nil {
		failureHandler.HandleFail(err)
	}
}