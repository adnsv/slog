package slog

import "time"

type FilterOpts struct {
	Trace bool
	Debug bool
}

func Filter(opts FilterOpts, target Sink) Sink {
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
