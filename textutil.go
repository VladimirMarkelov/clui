package clui

import (
	//"fmt"
	xs "github.com/huandu/xstrings"
	"regexp"
)

// Truncates text to maxWidth by replacing long
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

// Make a text no longer than maxWidth
func CutText(str string, maxWidth int) string {
	ln := xs.Len(str)
	if ln <= maxWidth {
		return str
	}

	return xs.Slice(str, 0, maxWidth)
}

func AlignText(str string, width int, align Align) (int, string) {
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

func UnColorizeText(str string) string {
	r1 := regexp.MustCompile("<f:[^>]*>")
	r2 := regexp.MustCompile("<t:[^>]*>")
	r3 := regexp.MustCompile("<b:[^>]*>")

	str = r1.ReplaceAllString(str, "")
	str = r2.ReplaceAllString(str, "")
	return r3.ReplaceAllString(str, "")
}
