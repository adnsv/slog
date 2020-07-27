package slog

type Progresser interface {
	Position(pos int64)
	Done(msg string)
}

func IncrementalProgressIndicator(maxPos int64) Progresser {
	rb := &runeBar{maxPos: maxPos}
	return rb
}

type runeBar struct {
	maxPos     int64
	showingPos int64 // [0..maxPos]
	showingStr string
}

const countdownStr = "9:.:8:.:7:.:6:.:5:.:4:.:3:.:2:.:1:.:0"

func (b *runeBar) Position(pos int64) {
	if b.showingPos == pos {
		return
	}

	b.showingPos = pos
	var str string

	if b.maxPos <= 0 {
		str = "..."
	} else {
		str = "["
		nmax := len(countdownStr)
		n := int(b.showingPos * int64(nmax) / b.maxPos)
		if n < 0 {
			n = 0
		} else if n > nmax {
			n = nmax
		}
		str += countdownStr[:n]
		if b.showingPos >= b.maxPos {
			str += "]"
		}
	}

	if len(str) > len(b.showingStr) {
		Print(str[len(b.showingStr):])
		b.showingStr = str
	}
}

func (b *runeBar) Done(msg string) {
	if len(b.showingStr) > 0 && len(msg) > 0 {
		Print(" ")
	}
	Print(msg)
}
