# sherlog [![CodeFactor](https://www.codefactor.io/repository/github/nick-anderssohn/sherlog/badge)](https://www.codefactor.io/repository/github/nick-anderssohn/sherlog) [![GoDoc](https://godoc.org/github.com/Nick-Anderssohn/sherlog?status.svg)](https://godoc.org/github.com/Nick-Anderssohn/sherlog) [![Build Status](https://travis-ci.org/Nick-Anderssohn/sherlog.svg?branch=master)](https://travis-ci.org/Nick-Anderssohn/sherlog)

Current version: 1.5.1 See releases [here](https://github.com/Nick-Anderssohn/sherlog/releases)

A simple logging library designed to be robust and performant (TELL ME THAT'S NOT A WORD. I DARE YOU). It's complete
with exceptions (no throw-catch obviously :disappointed:) and automatic stack traces. Creating an exception
with a stack trace is about five to six times faster than in Java (Benchmarks at bottom of readme). I am guessing this is simply due to the nature of Java...
It is likely there are more things that it needs to allocate when creating objects.

## Quick Start
Requires go 1.10 or higher!
Go get the package:
```
go get github.com/Nick-Anderssohn/sherlog
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
    return sherlog.AsError(err)
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

### Using sherlog.Logger
The sherlog.Logger interface is simple:
```
type Logger interface {
	Log(loggable Loggable) error
	Close()
}
```
sherlog offers a couple implementations of this. If you just want one log file, use `sherlog.FileLogger` or `sherlog.RollingFileLogger`. If you want different log files
for different log levels/combinations of log levels, use `sherlog.MultiFileLogger`. Check out the GoDocs for those. I recommend having one singleton
instance that you use for logging over your entire project. In the future, I will make an implementation of `sherlog.Logger` that automatically
sends logs to a gui, so it will be easy to add to your project.

#### Why no Package-Level Log Function?
Because I want people to write flexible code. Let's say I did provide a `Log` function that isn't tied to any struct, and it just calls
some hidden-away singleton instance of a logger that I have in the sherlog package. Let's say some dude named Cameron decided to use this 
function instead of creating an instance of `sherlog.FileLogger`. Cameron is happy because all of his stuff gets logged real ez, and he doesn't
have to do a lot of work. However, one day sherlog offers a new instance of Logger called `sherlog.GuiLogger` that automatically sends logs
to a database that can be used by a gui. If Cameron had already created a singleton instance of a Logger and used that, he could just instantiate
it with a `sherlog.GuiLogger` instead of `sherlog.FileLogger` to switch over and then call it a day. But, because he doesn't have an instance of Logger, he can't do that :( Interfaces
offer flexibility, so that's why I stick to the pattern of creating a logger, instead of just offering a function.

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

#### Distinguish between ERROR and OPS_ERROR
It is especially important to distinguish between ERROR and OPS_ERROR. The goal is to limit things marked with ERROR to only be messages that 100% indicate a bug in the code. OPS_ERROR may indicate an ops issue like something being offline...If there is a case where an error could indicate either a bug in the code or an ops error, then choose ops error.

## Benchmarks
All benchmarks ran on my desktop with Intel® Core™ i7-6700 CPU @ 3.40GHz × 8 running Ubuntu 18.04.1 LTS. The benchmark
for getting a stack trace in go was:
```
goos: linux
goarch: amd64
pkg: sherlog
BenchmarkGetStackTrace-8   	 1000000	      1580 ns/op
PASS
```
The benchmark for getting a stack trace in Java was:
```
Benchmark                                Mode  Cnt     Score     Error  Units
StackTraceBenchmark.BenchmarkStackTrace  avgt   30  9197.803 ± 134.794  ns/op
```
The code for the java benchmark is [here](https://github.com/Nick-Anderssohn/java-stack-trace-benchmark)

## GoDocs
See the auto-generated documentation [here](https://godoc.org/github.com/Nick-Anderssohn/sherlog)
