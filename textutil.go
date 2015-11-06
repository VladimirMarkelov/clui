package clui

import (
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
	"regexp"
)

// Ellipsize truncates text to maxWidth by replacing a
// substring in the middle with ellipsis and keeping
// the beginning and ending of the string untouched.
// If maxWidth is less than 5 then no ellipsis is
// added, the text is just truncated from the right.
func Ellipsize(str string, maxWidth int) string {
	ln := xs.Len(str)
	if ln <= maxWidth {
		return str
	}

	if maxWidth < 5 {
		return xs.Slice(str, 0, maxWidth)
	}

	left := int((maxWidth - 3) / 2)
	right := maxWidth - left - 3
	return xs.Slice(str, 0, left) + "..." + xs.Slice(str, ln-right, -1)
}

// CutText makes a text no longer than maxWidth
func CutText(str string, maxWidth int) string {
	ln := xs.Len(str)
	if ln <= maxWidth {
		return str
	}

	return xs.Slice(str, 0, maxWidth)
}

// AlignText calculates the initial position of the text
// output depending on str length and available width.
// The str is truncated in case of its lenght greater than
// width. Function returns shift that should be added to
// original label position before output instead of padding
// the string with spaces. The reason is to make possible
// to draw a label aligned but with transparent beginning
// and ending. If you do not need transparency you can
// add spaces manually using the returned shift value
func AlignText(str string, width int, align Align) (shift int, out string) {
	length := xs.Len(str)

	if length >= width {
		return 0, CutText(str, width)
	}

	if align == AlignRight {
		return width - length, str
	} else if align == AlignCenter {
		return (width - length) / 2, str
	}

	return 0, str
}

// AlignColorizedText does the same as AlignText does but
// it preserves the color of the letters byt adding correct
// colro tags to the line beginning.
// Note: function is ineffective and a bit slow - do not use
// it everywhere
func AlignColorizedText(str string, width int, align Align) (int, string) {
	rawText := UnColorizeText(str)
	length := xs.Len(rawText)

	if length <= width {
		shift, _ := AlignText(rawText, width, align)
		return shift, str
	}

	skip := 0
	if align == AlignRight {
		skip = length - width
	} else if align == AlignCenter {
		skip = (length - width) / 2
	}

	fgChanged, bgChanged := false, false
	curr := 0
	parser := NewColorParser(str, term.ColorBlack, term.ColorBlack)
	out := ""
	for curr <= skip+width {
		elem := parser.NextElement()

		if elem.Type == ElemEndOfText {
			break
		}

		if elem.Type == ElemPrintable {
			curr++
			if curr == skip+1 {
				if fgChanged {
					out += "<t:" + ColorToString(elem.Fg) + ">"
				}
				if bgChanged {
					out += "<b:" + ColorToString(elem.Bg) + ">"
				}
				out += string(elem.Ch)
			} else if curr > skip+1 {
				out += string(elem.Ch)
			}
		} else if elem.Type == ElemTextColor {
			fgChanged = true
			if curr > skip+1 {
				out += "<t:" + ColorToString(elem.Fg) + ">"
			}
		} else if elem.Type == ElemBackColor {
			bgChanged = true
			if curr > skip+1 {
				out += "<b:" + ColorToString(elem.Bg) + ">"
			}
		}
	}

	return 0, out
}

// SliceColorized returns a slice of text with correct color
// tags. start and end are real printable rune indices
func SliceColorized(str string, start, end int) string {
	if str == "" {
		return str
	}

	fgChanged, bgChanged := false, false
	curr := 0
	parser := NewColorParser(str, term.ColorBlack, term.ColorBlack)
	var out string
	for {
		if curr >= end {
			break
		}
		elem := parser.NextElement()
		if elem.Type == ElemEndOfText {
			break
		}

		switch elem.Type {
		case ElemTextColor:
			fgChanged = true
			if out != "" {
				out += "<t:" + ColorToString(elem.Fg) + ">"
			}
		case ElemBackColor:
			bgChanged = true
			if out != "" {
				out += "<b:" + ColorToString(elem.Bg) + ">"
			}
		case ElemPrintable:
			if curr == start {
				if fgChanged {
					out += "<t:" + ColorToString(elem.Fg) + ">"
				}
				if bgChanged {
					out += "<b:" + ColorToString(elem.Bg) + ">"
				}
			}
			if curr >= start {
				out += string(elem.Ch)
			}
			curr++
		}
	}

	return out
}

// UnColorizeText removes all color-related tags from the
// string. Tags to remove: <(f|t|b):.*>
func UnColorizeText(str string) string {
	r1 := regexp.MustCompile("<f:[^>]*>")
	r2 := regexp.MustCompile("<t:[^>]*>")
	r3 := regexp.MustCompile("<b:[^>]*>")

	str = r1.ReplaceAllString(str, "")
	str = r2.ReplaceAllString(str, "")
	return r3.ReplaceAllString(str, "")
}
