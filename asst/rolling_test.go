package asst

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestWindowOverflow(t *testing.T) {
	r := &Rolling{}
	// generate input longer than maxCols
	input := strings.Repeat("x", maxCols+5)
	r.Append(input)
	w := r.Window()
	if utf8.RuneCountInString(w) != maxCols {
		t.Errorf(
			"expected window rune count %d, got %d",
			maxCols, utf8.RuneCountInString(w),
		)
	}
	if !strings.HasPrefix(w, "…") {
		t.Errorf("expected ellipsis prefix, got %q", w)
	}
}

func TestStringEscapesNewline(t *testing.T) {
	r := &Rolling{}
	r.Append("hello\nworld")
	s := r.String()
	// should escape newline as \n
	if !strings.Contains(s, `\n`) {
		t.Errorf("expected escaped newline in %q", s)
	}
	// should not contain literal newline
	if strings.Contains(s, "\n") {
		t.Errorf("expected no literal newline in %q", s)
	}
}
