package sherlock

import (
	"testing"
	"os"
	"fmt"
)

/*
This file is just a place to put scratch code to visibly test things out.
Do not put any actual tests in here.
It is okay to pass in t *testing.T to your functions so that they can be run in your ide.
*/

func TestGetStackTraceAsString_Scratch(t *testing.T) {
	testStdException := NewStdException(testMessage)
	testStdException.Log(os.Stdout)
	fmt.Println("*****************************")
	test := NewLeveledException("Wubalubadubdub", EnumError)
	test.Log(os.Stdout)
}