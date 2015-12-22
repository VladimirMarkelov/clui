package clui

import (
	"testing"
)

func TestEllipsize(t *testing.T) {
	cases := []struct {
		in, want string
		max      int
	}{
		{"abcdefgh", "abcdefgh", -1},
		{"abcdefgh", "a...gh", 6},
		{"abcdefgh", "ab...gh", 7},
		{"abcdefgh", "abcdefgh", 10},
		{"abcdefgh", "abcd", 4},
	}

	for _, c := range cases {
		got := Ellipsize(c.in, c.max)
		if got != c.want {
			t.Errorf("Ellipsize (%v to %v) == <%v>, want <%v>", c.in, c.max, got, c.want)
		}
	}
}

func TestCutText(t *testing.T) {
	cases := []struct {
		in, want string
		max      int
	}{
		{"abcdefgh", "abcdefgh", -1},
		{"abcdefgh", "abcd", 4},
		{"abcdefgh", "abcde", 5},
		{"abcdefgh", "abcdefgh", 10},
	}

	for _, c := range cases {
		got := CutText(c.in, c.max)
		if got != c.want {
			t.Errorf("CutText (%v of %v) == <%v>, want <%v>", c.in, c.max, got, c.want)
		}
	}
}

func TestAlignText(t *testing.T) {
	cases := []struct {
		in, want   string
		align      Align
		max, shift int
	}{
		{"abcdefgh", "abcde", AlignLeft, 5, 0},
		{"abcdefgh", "defgh", AlignRight, 5, 0},
		{"abcdefgh", "bcdef", AlignCenter, 5, 0},
		{"abcdefgh", "abcdefgh", AlignLeft, 10, 0},
		{"abcdefgh", "abcdefgh", AlignRight, 10, 2},
		{"abcdefgh", "abcdefgh", AlignCenter, 10, 1},
		{"abcdefg", "abcdefg", AlignCenter, 10, 2},
	}

	for _, c := range cases {
		sh, got := AlignText(c.in, c.max, c.align)
		if got != c.want && sh != c.shift {
			t.Errorf("AlignText (%v of %v to %v) == <%v : %v>, want <%v : %v>", c.in, c.max, c.align, got, sh, c.want, c.shift)
		}
	}
}

func TestAlignColorizedText(t *testing.T) {
	cases := []struct {
		in, want   string
		align      Align
		max, shift int
	}{
		// uncolored cases
		{"abcdefgh", "abcde", AlignLeft, 5, 0},
		{"abcdefgh", "defgh", AlignRight, 5, 0},
		{"abcdefgh", "bcdef", AlignCenter, 5, 0},
		{"abcdefgh", "abcdefgh", AlignLeft, 10, 0},
		{"abcdefgh", "abcdefgh", AlignRight, 10, 2},
		{"abcdefgh", "abcdefgh", AlignCenter, 10, 1},
		{"abcdefg", "abcdefg", AlignCenter, 10, 2},
		// colored cases
		{"abc<t:green>defg", "abc<t:green>defg", AlignCenter, 10, 2},
		{"abc<t:green>defgh", "abc<t:green>defgh", AlignRight, 10, 2},
		{"abc<t:green>defgh", "abc<t:green>defgh", AlignCenter, 10, 1},
		{"<b:blue>ab<b:cyan>cdefgh", "<b:blue>ab<b:cyan>cde", AlignLeft, 5, 0},
		{"<b:blue>abcdefgh", "<b:cyan>defgh", AlignRight, 5, 0},
		{"<b:blue>abcdefgh", "<b:blue>b<b:cyan>cdef", AlignCenter, 5, 0},
		{"abc<t:green>defg", "ab", AlignLeft, 2, 0},
	}

	for _, c := range cases {
		sh, got := AlignColorizedText(c.in, c.max, c.align)
		if got != c.want && sh != c.shift {
			t.Errorf("AlignColorizedText (%v of %v to %v) == <%v : %v>, want <%v : %v>", c.in, c.max, c.align, got, sh, c.want, c.shift)
		}
	}
}

func TestSliceColorized(t *testing.T) {
	cases := []struct {
		in, want   string
		start, end int
	}{
		// uncolored cases
		{"abcdefgh", "abcde", 0, 5},
		{"abcdefgh", "defgh", 3, 9},
		{"abcdefgh", "bcdef", 1, 6},
		{"abcdefgh", "abcdefgh", 0, -1},
		{"abcdefgh", "abcde", -4, 5},
		// colored cases
		{"ab<t:blue>cde<t:green>fgh", "ab<t:blue>cde", 0, 5},
		{"ab<t:blue>cde<t:green>fgh", "<t:blue>de<t:green>fgh", 3, 9},
		{"ab<t:blue>cde<t:green>fgh", "b<t:blue>cde<t:green>f", 1, 6},
		{"ab<t:blue>cde<t:green>fgh", "<t:green>gh", 6, -1},
	}

	for _, c := range cases {
		got := SliceColorized(c.in, c.start, c.end)
		if got != c.want {
			t.Errorf("SliceColorized (%v from %v to %v) == <%v>, want <%v>", c.in, c.start, c.end, got, c.want)
		}
	}
}

func TestUnColorizeText(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"<f>abcd", "<f>abcd"},
		{"ab<f:>cd", "abcd"},
		{"ab<t:green>cd<b:blue>ef", "abcdef"},
		{"<f:black>", ""},
	}

	for _, c := range cases {
		got := UnColorizeText(c.in)
		if got != c.want {
			t.Errorf("UnColorize (%v) == <%v>, want <%v>", c.in, got, c.want)
		}
	}
}
