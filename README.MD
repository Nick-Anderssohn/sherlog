# Sherlock [![GoDoc](https://godoc.org/github.com/Nick-Anderssohn/sherlock?status.svg)](https://godoc.org/github.com/Nick-Anderssohn/sherlock)
A simple logging library designed to be very performant (TELL ME THAT'S NOT A WORD. I DARE YOU). It's complete
with exceptions (no throw-catch obviously :disappointed:) and automatic stack traces. Creating an exception
with a stack trace is about 5x faster than in Java. This is probably mostly due to Go being a faster language
than Java. Although, Java may also instantiate extra stuff while creating the stack trace and exception. Benchmarks
are farther down in the readme.

##GoDocs
See the auto-generated documentation [here](https://godoc.org/github.com/Nick-Anderssohn/sherlock)
