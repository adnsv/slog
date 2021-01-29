package slog

import "time"

// FilterOpts configures which levels to output
type FilterOpts struct {
	Trace bool
	Debug bool
}

// Filter creates an output sink with level filtering
func Filter(opts *FilterOpts, target Sink) Sink {
	return func(ts time.Time, lvl Level, prefix string, msg []byte) {
		allow := true
		switch lvl {
		case TraceLevel:
			allow = opts.Trace
		case DebugLevel:
			allow = opts.Debug
		}
		if allow && target != nil {
			target(ts, lvl, prefix, msg)
		}
	}
}
