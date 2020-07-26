package ansi

import "strconv"

// An index into a standard 256-color ANSI palette
type Index uint8

const Black = Index(0)
const White = Index(15)

func RGB(r, g, b byte) Index {
	rr := (uint(r) * 3) >> 7
	gg := (uint(g) * 3) >> 7
	bb := (uint(b) * 3) >> 7
	return Index(rr*36 + gg*6 + bb + 16)
}

func Gray(l byte) Index {
	v := (uint32(l)*25 + 128) >> 8
	if v == 0 {
		return Black
	} else if v == 25 {
		return White
	}
	v += 231
	return Index(v)
}

const ResetSeq = "\x1b[0m"

func FgSeq(v Index) string {
	return "\x1b[38;5;" + strconv.Itoa(int(v)) + "m"
}
func BgSeq(v Index) string {
	return "\x1b[48;5;" + strconv.Itoa(int(v)) + "m"
}
