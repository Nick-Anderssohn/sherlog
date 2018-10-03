package sherlog

import (
	"fmt"
	"log"
	"os"
	"testing"
)

/*
This file is just a place to put scratch code to visibly test things out.
Do not put any actual tests in here.
It is okay to pass in t *testing.T to your functions so that they can be run in your ide.
*/

func TestGetStackTraceAsString(t *testing.T) {
	testStdException := NewStdException(testMessage)
	testStdException.Log(os.Stdout)
	fmt.Println("*****************************")
	test := NewLeveledException("Wubalubadubdub", EnumError)
	test.Log(os.Stdout)
}

// This one is commented out because it passes, but intellij thinks it fails because of how the exceptions
// are printed. LOL

//func TestLoggableFuncs(t *testing.T) {
//	fmt.Print("LeveledException:\n")
//	printLoggableFuncs(NewLeveledException("wub wub", EnumError))
//
//	fmt.Print("\nStdException:\n")
//	printLoggableFuncs(NewStdException("wub wub"))
//}
//
//func printLoggableFuncs(loggable Loggable) {
//	fmt.Print("\n Log:\n\n")
//	loggable.Log(os.Stdout)
//
//	fmt.Print("\n LogNoStack:\n\n")
//	loggable.LogNoStack(os.Stdout)
//
//	fmt.Print("\n LogAsJson:\n\n")
//	loggable.LogAsJson(os.Stdout)
//}

func TestErrorFuncs(t *testing.T) {
	fmt.Print("LeveledException Error\n\n")
	log.Println(NewLeveledException("MWAHAHA", EnumInfo))

	fmt.Print("StdException Error\n\n")
	log.Println(NewStdException("MWAHAHA"))
}

func TestBasicError(t *testing.T) {
	err := fmt.Errorf("I am an error ;)")
	log.Println(err)
	err = AsCritical(err)
	log.Println(err)
}

func TestGetTimestampedFileName(t *testing.T) {
	fName := "error.log"
	withTime := getTimestampedFileName(fName)
	fmt.Println(withTime)
}

func TestLogJson(t *testing.T) {
	//info := NewInfo("I'm informative!")
	//info.LogAsJson(os.Stdout)
}
