package slog

import (
	"io"
	"time"
)

// PlainSink is a simple pass-through message formatter that outputs
// content as-is
func PlainSink(target io.Writer) Sink {
	return func(ts time.Time, lvl Level, domains []string, msg []byte) {
		target.Write(msg)
	}
}

// DecoratedSink provides adds timestamp, level, and domain to the output
func DecoratedSink(target io.Writer, decor DecorProvider) Sink {

	prefixer := func(ts time.Time, domains []string) []byte {
		if currLevel != stopped && decor != nil {
			return decor(ts, currLevel, domains, true)
		}
		return nil
	}

	postfixer := func(ts time.Time, domains []string) []byte {
		if currLevel != stopped && decor != nil {
			return decor(ts, currLevel, domains, false)
		}
		return nil
	}

	ld := lineDecorator(
		// prefix
		prefixer,
		//postfix
		postfixer,
		// target writer
		func(p []byte) {
			target.Write(p)
		})

	return func(ts time.Time, lvl Level, domains []string, msg []byte) {
		if len(msg) > 0 {
			ld(ts, domains, msg)
		}
	}
}

func lineDecorator(
	getPrefix, getPostfix func(ts time.Time, domains []string) []byte,
	target func(p []byte)) func(ts time.Time, domains []string, p []byte) {

	if target == nil {
		return func(time.Time, []string, []byte) {}
	}

	last := byte('\n')

	return func(ts time.Time, domains []string, p []byte) {

		for len(p) > 0 {
			n := len(p)
			if last == '\n' {
				// write prefix
				target(getPrefix(ts, domains))
			} else if last == '\r' && p[0] == '\n' {
				// write eol
				rn := [2]byte{'\r', '\n'}
				target(rn[:])
				last = '\n'
				p = p[1:]
				continue
			}

			i := 0
			for i < n {
				last = p[i]
				i++
				if last == '\n' {
					break
				}
			}

			if last == '\n' {
				t := i - 1

				if t > 0 && p[i-1] == '\r' {
					t--
				}
				// write content before (\r)?\n
				if t > 0 {
					target(p[:t])
				}
				// sink postfix
				target(getPostfix(ts, domains))
				// sink eol sequence
				target(p[t:i])
				p = p[i:]
			} else if last == '\r' {
				// write content before \r
				if i > 0 {
					target(p[:i-1])
				}
				// sink postfix
				target(getPostfix(ts, domains))
				p = p[i:]
			} else {
				target(p[:i])
				p = p[i:]
			}
		}
	}
}
