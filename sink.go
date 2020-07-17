package slog

import (
	"io"
	"time"
)

// PlainSink is a simple pass-through message formatter that outputs
// content as-is
func PlainSink(target io.Writer) Sink {
	currlvl := ContinueLevel
	return func(ts time.Time, lvl Level, prefix string, msg []byte) {
		if lvl == stopLevel {
			if currlvl != ContinueLevel {
				target.Write([]byte{'\n'})
			}
			currlvl = ContinueLevel
		} else if lvl != ContinueLevel {
			if currlvl != ContinueLevel {
				target.Write([]byte{'\n'})
			}
			currlvl = lvl
		}
		if len(msg) > 0 {
			target.Write(msg)
		}
	}
}

// DecoratedSink provides an opportunity to add timestamp, level, domain to the output
func DecoratedSink(target io.Writer, decor DecorProvider) Sink {
	currlvl := ContinueLevel
	currdomain := ""
	currts := time.Now()

	prefixer := func(ts time.Time, domain string) []byte {
		if currlvl != ContinueLevel && decor != nil {
			return decor(ts, currlvl, domain, true)
		}
		return nil
	}

	postfixer := func(ts time.Time, domain string) []byte {
		if currlvl != ContinueLevel && decor != nil {
			return decor(ts, currlvl, domain, false)
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

	return func(ts time.Time, lvl Level, domain string, msg []byte) {
		if lvl == stopLevel {
			if currlvl != ContinueLevel {
				ld(currts, currdomain, []byte{'\n'})
			}
			currlvl = ContinueLevel
		} else if lvl != ContinueLevel {
			if currlvl != ContinueLevel {
				ld(currts, currdomain, []byte{'\n'})
			}
			currlvl = lvl
		}
		currts = ts
		currdomain = domain
		if len(msg) > 0 {
			ld(ts, currdomain, msg)
		}
	}
}

func lineDecorator(getPrefix, getPostfix func(ts time.Time, domain string) []byte, target func(p []byte)) func(ts time.Time, domain string, p []byte) {

	if target == nil {
		return func(time.Time, string, []byte) {}
	}

	last := byte('\n')

	return func(ts time.Time, domain string, p []byte) {

		for len(p) > 0 {
			n := len(p)
			if last == '\n' {
				// write prefix
				target(getPrefix(ts, domain))
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
				target(getPostfix(ts, domain))
				// sink eol sequence
				target(p[t:i])
				p = p[i:]
			} else if last == '\r' {
				// write content before \r
				if i > 0 {
					target(p[:i-1])
				}
				// sink postfix
				target(getPostfix(ts, domain))
				p = p[i:]
			} else {
				target(p[:i])
				p = p[i:]
			}
		}
	}
}
