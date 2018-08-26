# Sherlock [![GoDoc](https://godoc.org/github.com/Nick-Anderssohn/sherlock?status.svg)](https://godoc.org/github.com/Nick-Anderssohn/sherlock)
A simple logging library designed to be robust and performant (TELL ME THAT'S NOT A WORD. I DARE YOU). It's complete
with exceptions (no throw-catch obviously :disappointed:) and automatic stack traces. Creating an exception
with a stack trace is about 5x faster than in Java. I am guessing Java is slower because it allocates more memory for
things while it creates the stack trace.

## Quick Start
Go get the package:
```
go get github.com/Nick-Anderssohn/sherlock
```
### Adding to an Existing Project
You most likely are grabbing errors like this:
```
err := someFunc()
if err != nil {
    return err
}
```
And then you are probably logging the error somewhere higher up in your code:
```
log.Println(err)
```
If so, your error is being logged like this:
```
2018/08/25 19:28:37 your error message
```
To get levels and stack traces added, all you have to do is use the `AsCritical`, `AsError`, `AsOpsError`, `AsWarning`, `AsInfo`, and `AsDebug`
functions. These functions graduate an error to a `LeveledException`, which creates a stack trace upon creation. For example, to give 
`err` the log level `ERROR` and give it a stack trace, all you have to do is change the code to this:
```
err := someFunc()
if err != nil {
    return sherlock.AsError(err)
}
```
There is no need to change the return type of your function, `AsError` and the others have a return type of `error`. You don't have to change
whatever logging code you are already using:
```
log.Println(err)
```
Now, the error will automatically be logged like this:
```
2018/08/25 19:28:37  - ERROR - your error message:
	something.WhateverFuncCalledAsError(/home/nick/go/src/something/something.go:57)
	something.Blah(/home/nick/go/src/something/something.go:777)
	something.Wonkers(/home/nick/go/src/something/something.go:2361)
```
That's it! That is all you need to do to get stack traces and levels added to your existing project :slightly_smiling_face:
In the future, I will add a gui to display, sort, search, and do whatever with your logs. So I recommend you also switch to my logger
which is documented below.
### Using sherlock.Logger
The sherlock.Logger interface is simple:
```
type Logger interface {
	Log(loggable Loggable) error
	Close()
}
```
sherlock offers a couple implementations of this. If you just want one log file, use `sherlock.FileLogger`. If you want different log files
for different log levels/combinations of log levels, use `sherlock.MultiFileLogger`. Check out the GoDocs for those. I recommend having one singleton
instance that you use for logging over your entire project. In the future, I will make an implementation of `sherlock.Logger` that automatically
sends logs to a gui, so it will be easy to add to your project.

## Log Levels
It is completely optional to use the default log levels I created. You can create your own by implementing this interface:
```
/**
An interface used to specify the log level on an exception/error.
LevelId is meant to be something along the lines of an enum, so
that we don't have to switch based off of the string value of the
log level. Label is the string representation.
 */
type Level interface {
	GetLevelId() int
	GetLabel() string
}
```
Of course, then you won't be able to use functions like `AsError`. You will have to manually create the error like this:
```
NewLeveledException("your error message", YourCustomeLogLevel)
```
### Why use my log levels
However, I do believe the log levels I chose create a great logging system. Here is what the different levels are intended for:

#### CRITICAL
CRITICAL is the intended log level for panics that are caught in the recover function.

#### ERROR
ERROR is the intended log level for something that should never ever happen and for sure
means there is a bug in your code.

#### OPS_ERROR
OPS_ERROR is the intended log level for an error that is known to be possible due to an operations issue.
For example, failing to query a database should be an OpsError because it lets you know that your database is
offline. It doesn't mean there is a bug in your code, but it is still something that needs to be fixed asap.

#### WARNING
WARNING is the intended log level for something that doesn't cause failure, but may be frowned
upon or bad anyways. For example, use of a deprecated endpoint may be logged as a warning. Or, a slow query could be a 
warning. A warning should say, "Hey, this is bad, but not a failure."

#### INFO
INFO is the intended log level for something that you want logged purely to collect information or metrics.

#### DEBUG
DEBUG is for any debug messages you want logged. Ideally, you are not logging these in production.

#### Distinguish between ERROR and WARNING
It is especially important to pay attention to the difference between an `ERROR` and a `WARNING`. Warnings can be equally
important as an error. For example, it is absolutely critical to get your database back online if it goes down! Leaving the
`ERROR` level to be exclusively used to indicate there is 100% a bug in the code will keep your logs cleaner and allow your
developers to track down bugs in the code. Don't be tempted to log a warning as an error just because it seems _important_.
Once again, warnings and errors are equally important, but you want to keep bug-identifying log messages separate from those
that indicate something could be wrong operationally.

## GoDocs
See the auto-generated documentation [here](https://godoc.org/github.com/Nick-Anderssohn/sherlock)
