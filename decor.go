package slog

import (
	"fmt"
	"strings"
	"time"

	"github.com/adnsv/slog/ansi"
)

// DecorProvider is an abstract callback that implements formatting of a messages
type DecorProvider = func(ts time.Time, lvl Level, domains []string, front bool) []byte

// BracketedDecorator implements a simple `timestamp [level:domain] ` prefix for message formatting
func BracketedDecorator(tsf TSFormat) func(ts time.Time, lvl Level, domains []string, front bool) []byte {

	l := map[Level]string{}
	ld := map[Level]string{}

	tss := ""
	if tsf&(TSDate|TSTime) != 0 {
		tss = "%s "
	}

	for lv, ln := range LevelNames {
		l[lv] = tss + "[" + ln + "] "
		ld[lv] = tss + "[" + ln + ":" + "%s] "
	}

	ds := ":" // domain separator

	if tsf&(TSDate|TSTime) != 0 {
		// with timestamp
		return func(ts time.Time, lvl Level, domains []string, front bool) []byte {
			if !front {
				return nil
			}
			if len(domains) > 0 {
				if f := ld[lvl]; len(f) > 0 {
					return []byte(fmt.Sprintf(f, FormatTimestamp(ts, tsf), strings.Join(domains, ds)))
				}
			} else if f := l[lvl]; len(f) > 0 {
				return []byte(fmt.Sprintf(f, FormatTimestamp(ts, tsf)))
			}
			return nil
		}
	}
	// without timestamp
	return func(ts time.Time, lvl Level, domains []string, front bool) []byte {
		if !front {
			return nil
		}
		if len(domains) > 0 {
			if f := ld[lvl]; len(f) > 0 {
				return []byte(fmt.Sprintf(f, strings.Join(domains, ds)))
			}
		} else if f := l[lvl]; len(f) > 0 {
			return []byte(f)
		}
		return nil
	}
}

// ColoredDecorator implements message formatting with ansi terminal colors
func ColoredDecorator(tsf TSFormat) func(ts time.Time, lvl Level, domains []string, front bool) []byte {

	l := map[Level]string{}
	ld := map[Level]string{}

	dss := ansi.FgSeq(cpdomain.fg) + ansi.BgSeq(cpdomain.bg) + " %s "
	tss := ""
	if tsf&(TSDate|TSTime) != 0 {
		tss = ansi.FgSeq(cptime.fg) + ansi.BgSeq(cptime.bg) + " %s "
	}

	ds := ":" // domain separator

	for lv, cp := range cplevel {
		l[lv] = tss + ansi.FgSeq(cp.fg) + ansi.BgSeq(cp.bg) + " " + LevelNames[lv] + " " + ansi.ResetSeq + " "
		ld[lv] = tss + ansi.FgSeq(cp.fg) + ansi.BgSeq(cp.bg) + " " + LevelNames[lv] + " " + dss + ansi.ResetSeq + " "
	}

	if tsf&(TSDate|TSTime) != 0 {
		// with timestamp
		return func(ts time.Time, lvl Level, domains []string, front bool) []byte {
			if !front {
				return nil
			}
			if len(domains) > 0 {
				if f := ld[lvl]; len(f) > 0 {
					return []byte(fmt.Sprintf(f, FormatTimestamp(ts, tsf), strings.Join(domains, ds)))
				}
			} else if f := l[lvl]; len(f) > 0 {
				return []byte(fmt.Sprintf(f, FormatTimestamp(ts, tsf)))
			}
			return nil
		}
	}
	// without timestamp
	return func(ts time.Time, lvl Level, domains []string, front bool) []byte {
		if !front {
			return nil
		}
		if len(domains) > 0 {
			if f := ld[lvl]; len(f) > 0 {
				return []byte(fmt.Sprintf(f, strings.Join(domains, ds)))
			}
		} else if f := l[lvl]; len(f) > 0 {
			return []byte(f)
		}
		return nil
	}
}

type cpair struct{ fg, bg ansi.Index }

var cptime = cpair{fg: ansi.Black, bg: ansi.Gray(128)}
var cpdomain = cpair{fg: ansi.White, bg: ansi.Gray(64)}

var cplevel = map[Level]cpair{
	TraceLevel: {fg: ansi.Gray(96), bg: ansi.RGB(32, 32, 48)},
	DebugLevel: {fg: ansi.Gray(96), bg: ansi.RGB(72, 32, 64)},
	InfoLevel:  {fg: ansi.Gray(240), bg: ansi.RGB(64, 64, 96)},
	WarnLevel:  {fg: ansi.Black, bg: ansi.RGB(128, 128, 64)},
	ErrorLevel: {fg: ansi.White, bg: ansi.RGB(128, 64, 64)},
	FatalLevel: {fg: ansi.Black, bg: ansi.RGB(240, 64, 64)},
}
