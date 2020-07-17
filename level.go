package slog

// Level indicates a logging message priority
type Level int

// Supported Level values
const (
	stopLevel     = Level(-1)
	ContinueLevel = Level(iota) // keep appending to last level
	TraceLevel
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
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
