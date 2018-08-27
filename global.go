package sherlog

import "time"

const (
	defaultStackTraceDepth    = 64
	defaultStackTraceLineLen  = 96
	defaultStackTraceNumBytes = defaultStackTraceLineLen * defaultStackTraceDepth
	timeFmt                   = "2006-01-02 15:04:05" // yyyy-mm-dd hh:mm:ss
	timeFileNameFmt           = "_2006-01-02"
)

var (
	/*
	All timestamps and time-based logging will use this location.
	Defaults to UTC. Can set to a different zone from the IANA Time Zone database.
	For example, if you want pacific time, you can do:
		sherlog.SherlogLocation, _ = time.LoadLocation("America/Los_Angeles")
	Wikipedia has a good list of IANA time zones: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
	 */
	SherlogLocation = time.UTC
)