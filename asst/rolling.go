package asst

import (
	"fmt"
)

const (
	maxCols    = 30
	eraseLine  = "\r\x1b[K"
	startColor = "\x1b[90m"
	endColor   = "\x1b[0m"
)

type Rolling struct {
	Buf []rune
}

func (r *Rolling) Append(s string) {
	r.Buf = append(r.Buf, []rune(s)...)
}

func (r *Rolling) Window() string {
	if l := len(r.Buf); l > maxCols {
		return "…" + string(r.Buf[l-maxCols+1:])
	}
	return string(r.Buf)
}

func (r *Rolling) String() string {
	return fmt.Sprintf("%s%s%q%s", eraseLine, startColor, r.Window(), endColor)
}
