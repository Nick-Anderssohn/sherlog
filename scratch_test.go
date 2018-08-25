package sherlock

import (
	"testing"
	"fmt"
)

/*
This file is just a place to put scratch code to visibly test things out.
Do not put any actual tests in here.
It is okay to pass in t *testing.T to your functions so that they can be run in your ide.
*/

func TestGetStackTraceAsString_Scratch(t *testing.T) {
	testStdException := NewStdException(testMessage)
	fmt.Println("This is a test StdException:")
	fmt.Println(testStdException.GetStackTraceAsString())
}

