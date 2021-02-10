package slog

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Sink is a core logging callback
type Sink = func(timestamp time.Time, lvl Level, domains []string, msg []byte)

// RootSink is the currently active Sink there all the messages will go
var RootSink Sink = DecoratedSink(os.Stderr, BracketedDecorator(TSTime|TSMicroseconds))

var (
	currTime   = time.Time{}
	currLevel  = stopped
	currDomain = ""
	eol        = []byte{'\n'}
	domains    []string
)

// CurrentLevel returns the current logging level
func CurrentLevel() Level {
	return currLevel
}

// CurrentDomain returns the current logging domain
func CurrentDomain() string {
	return currDomain
}

// PushDomain pushes a domain into the domain chain
func PushDomain(s string) {
	domains = append(domains, s)
}

// PopDomain removes the last item from the domain chain
func PopDomain() {
	n := len(domains)
	if n > 0 {
		domains = domains[:n-1]
	}
}

func domainChain() []string {
	if currDomain == "" {
		return domains
	}
	return append(domains, currDomain)
}

// Flush finalizes the last logging entry.
// Normally, this means writing a pending EOL into the output sink.
// This function must be called at the end of the application
func Flush() {
	if currLevel != stopped && RootSink != nil {
		RootSink(currTime, currLevel, domainChain(), eol)
	}
	currLevel = stopped
}

// Stop is replaced with Flush
// Deprecated: use Flush() instead
func Stop() {
	Flush()
}

// StartLevel starts a new line that targets the specified logging level and domain
func StartLevel(lvl Level, domain string) {
	Flush()
	currTime = time.Now()
	currLevel = lvl
	currDomain = domain
}

// WantLevel starts a new level/domain line if it differs from the current
func WantLevel(lvl Level, domain string) {
	if currLevel != lvl || currDomain != domain {
		Flush()
		currTime = time.Now()
		currLevel = lvl
		currDomain = domain
	}
}

// Append adds []byte content to current level
func Append(msg []byte) {
	if currLevel != stopped && RootSink != nil {
		currTime = time.Now()
		RootSink(currTime, currLevel, domainChain(), msg)
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

type appendWriter struct {
}

func (w *appendWriter) Write(p []byte) (n int, err error) {
	Append(p)
	return len(p), nil
}

type teeWriter struct {
	output io.Writer
}

func (t *teeWriter) Write(p []byte) (n int, err error) {
	Append(p)
	if t.output != nil {
		return t.output.Write(p)
	}
	return len(p), nil
}

// Writer produces an io.Writer that can be used for streaming
// into logs
func Writer(lvl Level, domain string) io.Writer {
	StartLevel(lvl, domain)
	return &appendWriter{}
}

// TeeWriter produces an io.Writer that appends to the current level
// and also writes to another output
func TeeWriter(lvl Level, domain string, output io.Writer) io.Writer {
	StartLevel(lvl, domain)
	return &teeWriter{output: output}
}
