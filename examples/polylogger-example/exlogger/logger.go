package exlogger

import (
	"sherlog"
	"sherlog/examples/exception-returner"
)

// I recommend you create your own logger package in your project to hold the singleton instance
// of a sherlog logger

var Logger sherlog.Logger

// I want to initialize the Logger as soon as the package is created, so put in a init func.
func init() {
	// I want to log everything to one file and log things to separate files.
	// I want all files to be rolled.
	// So I will use a PolyLogger with a RollingFileLogger for the one big file, and a MultiFileLogger
	// for the files that are separated based off of log level.
	// If I cannot instantiate one of the loggers, I don't even want my program to launch so I will panic.

	messagesForFilesToHoldBeforeRolled := 500

	allLogMessagesLogger, err := sherlog.NewRollingFileLogger("all_log_messages.log", messagesForFilesToHoldBeforeRolled)
	if err != nil {
		panic(err)
	}

	// Declare the paths for the multi file logger

	// I want criticals and errors to share a file.
	critAndErrorFilePath := "criticals_and_errors.log"

	// I want ops errors to have their own file
	opsErrorFilePath := "opsErrors.log"

	// I want warnings to have their own file
	warningFilePath := "warnings.log"

	// I want info and debug to share a file.
	infoAndDebugPath := "info_and_debug.log"

	// I want my own file for the custom log level I created called WEIRD_LOG_LEVEL
	customLogLevelPath := "custom.log"

	// I want errors that don't have a log level (aka I forgot to use a sherlog func on it) to go here
	unknownLevelPath := "unknown.log"

	paths := map[sherlog.Level]string {
		// provide same path for crit and error since I want them to share a file
		sherlog.EnumCritical: critAndErrorFilePath,
		sherlog.EnumError: critAndErrorFilePath,
		sherlog.EnumOpsError: opsErrorFilePath,
		sherlog.EnumWarning: warningFilePath,
		// provide same path for info and debug since I want them to share a file
		sherlog.EnumInfo: infoAndDebugPath,
		sherlog.EnumDebug: infoAndDebugPath,
		// Can include log levels not from sherlog as long as they implement sherlog.Level
		exception_returner.WeirdLogLevel: customLogLevelPath,
	}

	multiFileLogger, err := sherlog.NewMultiFileLoggerWithRollingLogs(paths, unknownLevelPath, messagesForFilesToHoldBeforeRolled)
	if err != nil {
		// couldn't create the logger, don't run the program
		panic(err)
	}

	// Now finally instantiate the singleton logger that the program will use
	Logger = sherlog.NewPolyLogger([]sherlog.Logger{allLogMessagesLogger, multiFileLogger})
}