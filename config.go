package slog

import (
	"io"
	"os"

	"github.com/adnsv/slog/ansi"
)

type OutputFormat int

const (
	AutoColorOutput = OutputFormat(iota) // use ANSI colors only if VT is available
	PlainOutput                          // use plain bracketed format
	ColorOutput                          // emit ansi colors
)

type Options struct {
	Timestamp TSFormat
	Output    io.Writer
	Format    OutputFormat
	Filter    FilterOpts
}

var termRestore = func() {}

func Close() {
	termRestore()
	termRestore = func() {}
}

func Configure(opt Options) {
	Close()

	output := opt.Output
	if output == nil {
		output = os.Stdout
	}

	var decorator DecorProvider
	if opt.Format == AutoColorOutput {
		f, isfile := output.(*os.File)
		if isfile {
			ok, cleanup := ansi.EnableVT(f)
			if ok {
				decorator = ColoredDecorator(opt.Timestamp)
				termRestore = cleanup
			}
		}
		if decorator == nil {
			decorator = BracketedDecorator(opt.Timestamp)
		}
	} else if opt.Format == ColorOutput {
		decorator = ColoredDecorator(opt.Timestamp)
	} else {
		decorator = BracketedDecorator(opt.Timestamp)
	}

	RootSink = Filter(&opt.Filter, DecoratedSink(output, decorator))
}

func (opt *Options) Apply(ss ...string) {
	for _, s := range ss {
		switch s {
		case "notime":
			opt.Timestamp &= ^(TSTime | TSMicroseconds | TSUTC)
		case "time":
			opt.Timestamp |= TSTime
		case "microsecond":
			opt.Timestamp |= TSTime | TSMicroseconds
		case "utc":
			opt.Timestamp |= TSTime | TSUTC
		case "date":
			opt.Timestamp |= TSDate
		case "nodate":
			opt.Timestamp &= ^TSDate
		case "debug":
			opt.Filter.Debug = true
		case "trace":
			opt.Filter.Trace = true
		case "stdout":
			opt.Output = os.Stdout
		case "stderr":
			opt.Output = os.Stderr
		}
	}
}
