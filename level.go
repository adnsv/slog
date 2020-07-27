package slog

// Level indicates a logging message priority
type Level int

// Supported Level values
const (
	TraceLevel = Level(iota + 1)
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// internals
const (
	stopped = Level(0) // no level
	
)

// LevelNames is a map of default level names
var LevelNames = map[Level]string{
	TraceLevel: "TRCE",
	DebugLevel: "DBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERR!",
	FatalLevel: "FATAL!",
}
