package clui

import (
	term "github.com/nsf/termbox-go"
	мИнт "./пакИнтерфейсы"
)

// TextElementType type of the parsed element of the string
type TextElementType int

// TextElementType values
const (
	// ElemPrintable - the item is a rune
	ElemPrintable = iota
	// ElemBackColor - the item sets new background color
	ElemBackColor
	// ElemTextColor - the item sets new text color
	ElemTextColor
	// ElemLineBreak - line break
	ElemLineBreak
	// ElemEndOfText - the string parsing has complited
	ElemEndOfText
)

// TextElement is currently parsed text element
type TextElement struct {
	// Type is an element type
	Type TextElementType
	// Ch is a parsed rune, it is filled only if Type is ElemPrintable
	Ch rune
	// Fg is a text color for the rune
	Fg term.Attribute
	// Bg is a background color for the rune
	Bg term.Attribute
}

// ColorParser is a string parser to process a text with color tags
// inside the string
type ColorParser struct {
	text     []rune
	index    int
	defBack  term.Attribute
	defText  term.Attribute
	currBack term.Attribute
	currText term.Attribute
}

// NewColorParser creates a new string parser.
// str is a string to parse.
// defText is a default text color.
// defBack is a default background color.
// Default colors are applied in case of reset color tag
func NewColorParser(str string, defText, defBack term.Attribute) *ColorParser {
	p := new(ColorParser)
	p.text = []rune(str)
	p.defBack, p.defText = defBack, defText
	p.currBack, p.currText = defBack, defText
	return p
}

func (p *ColorParser) parseColor() (term.Attribute, TextElementType, bool) {
	newIdx := p.index + 1
	length := len(p.text)
	ok := true

	const (
		StepType = iota
		StepColon
		StepValue
	)

	var (
		cText string
		attr  term.Attribute
		t     TextElementType
		done  bool
	)
	step := StepType

	for {
		if newIdx >= length {
			ok = false
			break
		}

		switch step {
		case StepType:
			c := p.text[newIdx]
			if c == 't' || c == 'f' || c == 'c' {
				t = ElemTextColor
			} else if c == 'b' {
				t = ElemBackColor
			} else {
				ok = false
				break
			}
			step = StepColon
			newIdx++
		case StepColon:
			c := p.text[newIdx]
			if c != ':' {
				ok = false
				break
			}
			newIdx++
			step = StepValue
		case StepValue:
			c := p.text[newIdx]
			if c == '>' {
				p.index = newIdx + 1
				if cText == "" {
					attr = мИнт.ColorDefault
				} else {
					attr = StringToColor(cText)
				}
				done = true
				break
			} else {
				if c != ' ' || cText != "" {
					cText += string(c)
				}
				newIdx++
			}
		}

		if done || !ok {
			break
		}
	}

	return attr, t, ok
}

// NextElement parses and returns the next string element
func (p *ColorParser) NextElement() TextElement {
	if p.index >= len(p.text) {
		return TextElement{Type: ElemEndOfText}
	}

	if p.text[p.index] == '\n' {
		p.index++
		return TextElement{Type: ElemLineBreak}
	}

	if p.text[p.index] != '<' {
		p.index++
		return TextElement{Type: ElemPrintable, Ch: p.text[p.index-1], Fg: p.currText, Bg: p.currBack}
	}

	attr, atype, ok := p.parseColor()
	if !ok {
		p.index++
		return TextElement{Type: ElemPrintable, Ch: p.text[p.index-1], Fg: p.currText, Bg: p.currBack}
	}

	if atype == ElemBackColor {
		if attr == мИнт.ColorDefault {
			p.currBack = p.defBack
		} else {
			p.currBack = attr
		}
	} else if atype == ElemTextColor {
		if attr == мИнт.ColorDefault {
			p.currText = p.defText
		} else {
			p.currText = attr
		}
	}
	return TextElement{Type: atype, Fg: p.currText, Bg: p.currBack}
}
