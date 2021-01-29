package slog

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Sink is a core logging callback
type Sink = func(timestamp time.Time, lvl Level, domain string, msg []byte)

// RootSink is the currently active Sink there all the messages will go
var RootSink Sink = DecoratedSink(os.Stderr, BracketedDecorator(TSTime|TSMicroseconds))

var (
	currTime   = time.Time{}
	currLevel  = stopped
	currDomain = ""
	eol        = []byte{'\n'}
)

// CurrentLevel returns the current logging level
func CurrentLevel() Level {
	return currLevel
}

// CurrentDomain returns the current logging domain
func CurrentDomain() string {
	return currDomain
}

// Stop finalizes the last logging entry and disables further logging.
// Normally, this means writing a pending EOL into the output sink.
// This function must be called at the end of the application
func Stop() {
	if currLevel != stopped && RootSink != nil {
		RootSink(currTime, currLevel, currDomain, eol)
	}
	currLevel = stopped
}

// StartLevel starts a new line that targets the specified logging level and domain
func StartLevel(lvl Level, domain string) {
	Stop()
	currTime = time.Now()
	currLevel = lvl
	currDomain = domain
}

// WantLevel starts a new level/domain line if it differs from the current
func WantLevel(lvl Level, domain string) {
	if currLevel != lvl || currDomain != domain {
		Stop()
		currTime = time.Now()
		currLevel = lvl
		currDomain = domain
	}
}

// Append adds []byte content to current level
func Append(msg []byte) {
	if currLevel != stopped && RootSink != nil {
		RootSink(currTime, currLevel, currDomain, msg)
	}
}

// AppendStr adds string content to current level
func AppendStr(s string) {
	Append([]byte(s))
}

// Print adds stringers to current level
func Print(a ...interface{}) {
	AppendStr(fmt.Sprint(a...))
}

// Printf adds formatted output to current level
func Printf(format string, a ...interface{}) {
	AppendStr(fmt.Sprintf(format, a...))
}

// common logging functions

// Trace starts TraceLevel and adds stringers to it
func Trace(a ...interface{}) {
	StartLevel(TraceLevel, "")
	Print(a...)
}

// Tracef starts TraceLevel and adds formatted output to it
func Tracef(format string, a ...interface{}) {
	StartLevel(TraceLevel, "")
	Printf(format, a...)
}

// Debug starts DebugLevel and adds stringers to it
func Debug(a ...interface{}) {
	StartLevel(DebugLevel, "")
	Print(a...)
}

// Debugf starts DebugLevel and adds formatted content to it
func Debugf(format string, a ...interface{}) {
	StartLevel(DebugLevel, "")
	Printf(format, a...)
}

// Info starts InfoLevel and adds stringers to it
func Info(a ...interface{}) {
	StartLevel(InfoLevel, "")
	Print(a...)
}

// Infof starts InfoLevel and adds formatted content to it
func Infof(format string, a ...interface{}) {
	StartLevel(InfoLevel, "")
	Printf(format, a...)
}

// Warn starts WarnLevel and adds stringers to it
func Warn(a ...interface{}) {
	StartLevel(WarnLevel, "")
	Print(a...)
}

// Warnf starts WarnLevel and adds formatted content to it
func Warnf(format string, a ...interface{}) {
	StartLevel(WarnLevel, "")
	Printf(format, a...)
}

// Error starts ErrorLevel and adds stringers to it
func Error(a ...interface{}) {
	StartLevel(ErrorLevel, "")
	Print(a...)
}

// Errorf starts ErrorLevel and adds formatted content to it
func Errorf(format string, a ...interface{}) {
	StartLevel(ErrorLevel, "")
	Printf(format, a...)
}

// Fatal starts FatalLevel and adds stringers to it
func Fatal(a ...interface{}) {
	StartLevel(FatalLevel, "")
	Print(a...)
}

// Fatalf starts FatalLevel and adds formatted content to it
func Fatalf(format string, a ...interface{}) {
	StartLevel(FatalLevel, "")
	Printf(format, a...)
}

// logging with domain

func DomainTrace(domain string, a ...interface{}) {
	StartLevel(TraceLevel, domain)
	Print(a...)
}

func DomainTracef(domain string, format string, a ...interface{}) {
	StartLevel(TraceLevel, domain)
	Printf(format, a...)
}

func DomainDebug(domain string, a ...interface{}) {
	StartLevel(DebugLevel, domain)
	Print(a...)
}

func DomainDebugf(domain string, format string, a ...interface{}) {
	StartLevel(DebugLevel, domain)
	Printf(format, a...)
}

func DomainInfo(domain string, a ...interface{}) {
	StartLevel(InfoLevel, domain)
	Print(a...)
}

func DomainInfof(domain string, format string, a ...interface{}) {
	StartLevel(InfoLevel, domain)
	Printf(format, a...)
}

func DomainWarn(domain string, a ...interface{}) {
	StartLevel(WarnLevel, domain)
	Print(a...)
}

func DomainWarnf(domain string, format string, a ...interface{}) {
	StartLevel(WarnLevel, domain)
	Printf(format, a...)
}

func DomainError(domain string, a ...interface{}) {
	StartLevel(ErrorLevel, domain)
	Print(a...)
}

func DomainErrorf(domain string, format string, a ...interface{}) {
	StartLevel(ErrorLevel, domain)
	Printf(format, a...)
}

func DomainFatal(domain string, a ...interface{}) {
	StartLevel(FatalLevel, domain)
	Print(a...)
}

func DomainFatalf(domain string, format string, a ...interface{}) {
	StartLevel(FatalLevel, domain)
	Printf(format, a...)
}

// stream adapters

type funcWriterAdapter func(p []byte)

func (fw funcWriterAdapter) Write(p []byte) (n int, err error) {
	n = len(p)
	if fw != nil {
		fw(p)
	}
	return
}

// LevelWriter produces an io.Writer that can be used for streaming
// into logs
func LevelWriter(lvl Level, domain string) io.Writer {
	StartLevel(lvl, domain)
	return funcWriterAdapter(Append)
}

// Writer produces an io.Writer that appends to the current level
func Writer() io.Writer {
	return funcWriterAdapter(Append)
}

func TraceWriter(domain string) io.Writer {
	StartLevel(TraceLevel, domain)
	return Writer()
}

func DebugWriter(domain string) io.Writer {
	StartLevel(DebugLevel, domain)
	return Writer()
}

func InfoWriter(domain string) io.Writer {
	StartLevel(InfoLevel, domain)
	return Writer()
}

func WarnWriter(domain string) io.Writer {
	StartLevel(WarnLevel, domain)
	return Writer()
}

func ErrorWriter(domain string) io.Writer {
	StartLevel(ErrorLevel, domain)
	return Writer()
}

func FatalWriter(domain string) io.Writer {
	StartLevel(FatalLevel, domain)
	return Writer()
}
