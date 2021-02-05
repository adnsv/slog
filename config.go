package slog

import (
	"io"
	"os"
	"strings"

	"github.com/adnsv/slog/ansi"
)

// OutputFormat specifies if the output should be ANSI-colored
type OutputFormat int

// Accepted OutputFormat values
const (
	AutoColorOutput = OutputFormat(iota) // use ANSI colors only if VT is available
	PlainOutput                          // use plain bracketed format
	ColorOutput                          // emit ANSI colors
)

// Options specifies formatting style and output target. Pass this as a
// parameter to the Configure function.
type Options struct {
	Timestamp TSFormat
	Output    io.Writer
	Format    OutputFormat
	Filter    FilterOpts
}

var termRestore = func() {}

// Close restores the terminal to its original space
// This may be required if the terminal was switched into
// colored output mode
func Close() {
	termRestore()
	termRestore = func() {}
}

// Configure configures slog output target and formatting according to the
// specified options
func Configure(opt *Options) {
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

var ConfigurationTokens = []string{
	"notime",
	"time",
	"microsecond",
	"utc",
	"date",
	"nodate",
	"debug",
	"trace",
	"plain",
	"color",
	"stdout",
	"stderr",
}

// Apply applies a set of configuration tokens to slog options
// This is useful when configuring slog from command line parameters
func (opt *Options) Apply(ss ...string) {
	for _, s := range ss {
		switch strings.TrimSpace(s) {
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
		case "plain":
			opt.Format = PlainOutput
		case "color":
			opt.Format = ColorOutput
		case "stdout":
			opt.Output = os.Stdout
		case "stderr":
			opt.Output = os.Stderr
		}
	}
}
