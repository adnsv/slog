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

// Stop should be called at the end of the application to emit the last EOL
func Stop() {
	if RootSink != nil {
		RootSink(time.Now(), stopLevel, "", nil)
	}
}

func levelMsg(lvl Level, domain string, msg string) {
	if RootSink != nil {
		RootSink(time.Now(), lvl, domain, []byte(msg))
	}
}

// common logging functions

func Trace(a ...interface{}) {
	levelMsg(TraceLevel, "", fmt.Sprint(a...))
}

func Tracef(format string, a ...interface{}) {
	levelMsg(TraceLevel, "", fmt.Sprintf(format, a...))
}

func Debug(a ...interface{}) {
	levelMsg(DebugLevel, "", fmt.Sprint(a...))
}

func Debugf(format string, a ...interface{}) {
	levelMsg(DebugLevel, "", fmt.Sprintf(format, a...))
}

func Info(a ...interface{}) {
	levelMsg(InfoLevel, "", fmt.Sprint(a...))
}

func Infof(format string, a ...interface{}) {
	levelMsg(InfoLevel, "", fmt.Sprintf(format, a...))
}

func Warn(a ...interface{}) {
	levelMsg(WarnLevel, "", fmt.Sprint(a...))
}

func Warnf(format string, a ...interface{}) {
	levelMsg(WarnLevel, "", fmt.Sprintf(format, a...))
}

func Error(a ...interface{}) {
	levelMsg(ErrorLevel, "", fmt.Sprint(a...))
}

func Errorf(format string, a ...interface{}) {
	levelMsg(ErrorLevel, "", fmt.Sprintf(format, a...))
}

func Fatal(a ...interface{}) {
	levelMsg(FatalLevel, "", fmt.Sprint(a...))
	Stop()
}

func Fatalf(format string, a ...interface{}) {
	levelMsg(FatalLevel, "", fmt.Sprintf(format, a...))
	Stop()
}

// logging with domain

func DomainTrace(domain string, a ...interface{}) {
	levelMsg(TraceLevel, domain, fmt.Sprint(a...))
}

func DomainTracef(domain string, format string, a ...interface{}) {
	levelMsg(TraceLevel, domain, fmt.Sprintf(format, a...))
}

func DomainDebug(domain string, a ...interface{}) {
	levelMsg(DebugLevel, domain, fmt.Sprint(a...))
}

func DomainDebugf(domain string, format string, a ...interface{}) {
	levelMsg(DebugLevel, domain, fmt.Sprintf(format, a...))
}

func DomainInfo(domain string, a ...interface{}) {
	levelMsg(InfoLevel, domain, fmt.Sprint(a...))
}

func DomainInfof(domain string, format string, a ...interface{}) {
	levelMsg(InfoLevel, domain, fmt.Sprintf(format, a...))
}

func DomainWarn(domain string, a ...interface{}) {
	levelMsg(WarnLevel, domain, fmt.Sprint(a...))
}

func DomainWarnf(domain string, format string, a ...interface{}) {
	levelMsg(WarnLevel, domain, fmt.Sprintf(format, a...))
}

func DomainError(domain string, a ...interface{}) {
	levelMsg(ErrorLevel, domain, fmt.Sprint(a...))
}

func DomainErrorf(domain string, format string, a ...interface{}) {
	levelMsg(ErrorLevel, domain, fmt.Sprintf(format, a...))
}

func DomainFatal(domain string, a ...interface{}) {
	levelMsg(FatalLevel, domain, fmt.Sprint(a...))
	Stop()
}

func DomainFatalf(domain string, format string, a ...interface{}) {
	levelMsg(FatalLevel, domain, fmt.Sprintf(format, a...))
	Stop()
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

// levelWriter produces an io.Writer that can be used for streaming
// into logs
func levelWriter(lvl Level, domain string) io.Writer {
	if RootSink != nil {
		RootSink(time.Now(), lvl, domain, nil)
	}

	f := func(p []byte) {
		if RootSink != nil {
			RootSink(time.Now(), ContinueLevel, domain, p)
		}
	}
	return funcWriterAdapter(f)
}

func TraceWriter(domain string) io.Writer {
	return levelWriter(TraceLevel, domain)
}

func DebugWriter(domain string) io.Writer {
	return levelWriter(DebugLevel, domain)
}

func InfoWriter(domain string) io.Writer {
	return levelWriter(InfoLevel, domain)
}

func WarnWriter(domain string) io.Writer {
	return levelWriter(WarnLevel, domain)
}

func ErrorWriter(domain string) io.Writer {
	return levelWriter(ErrorLevel, domain)
}

func FatalWriter(domain string) io.Writer {
	return levelWriter(FatalLevel, domain)
}
