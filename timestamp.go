package slog

import "time"

// TSFormat control conversion of timestamp to a string
type TSFormat uint

// Accepted TSFormat flags
const (
	TSNone         = TSFormat(0)
	TSDate         = TSFormat(1 << iota) // date in 2009/01/23
	TSTime                               // time 01:23:23
	TSMicroseconds                       // add microseconds to time
	TSUTC                                // use UTC time instead of local
)

// FormatTimestamp implements conversion of time to a string using the specified
// formatting flags
func FormatTimestamp(t time.Time, fmt TSFormat) string {

	if fmt&(TSDate|TSTime) == 0 {
		return ""
	}

	if fmt&TSUTC != 0 {
		t = t.UTC()
	}

	buf := make([]byte, 0, 64)

	wc := func(c byte) {
		buf = append(buf, c)
	}

	wn := func(i int, w int) {
		var b [32]byte
		bp := len(b) - 1
		for i >= 10 || w > 1 {
			w--
			q := i / 10
			b[bp] = byte('0' + i - q*10)
			bp--
			i = q
		}
		// i < 10
		b[bp] = byte('0' + i)
		buf = append(buf, b[bp:]...)
	}

	if fmt&TSDate != 0 {
		y, m, d := t.Date()
		wn(y, 4)
		wc('/')
		wn(int(m), 2)
		wc('/')
		wn(d, 2)
		if fmt&TSTime != 0 {
			wc(' ')
		}
	}
	if fmt&TSTime != 0 {
		h, m, s := t.Clock()
		wn(h, 2)
		wc(':')
		wn(m, 2)
		wc(':')
		wn(s, 2)
		if fmt&TSMicroseconds != 0 {
			wc('.')
			wn(t.Nanosecond()/1000, 6)
		}
	}

	return string(buf)
}
