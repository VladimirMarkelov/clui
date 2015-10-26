package clui

import (
	term "github.com/nsf/termbox-go"
)

const (
	ElemPrintable = iota
	ElemBackColor
	ElemTextColor
	ElemLineBreak
	ElemEndOfText
)

type TextElementType int

type TextElement struct {
	Type TextElementType
	Ch   rune
	Fg   term.Attribute
	Bg   term.Attribute
}

type ColorParser struct {
	text     []rune
	index    int
	defBack  term.Attribute
	defText  term.Attribute
	currBack term.Attribute
	currText term.Attribute
}

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
		step  int = StepType
		done  bool
	)

	for {
		if newIdx >= length {
			ok = false
			break
		}

		switch step {
		case StepType:
			c := p.text[newIdx]
			if c == 't' || c == 'f' {
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
					attr = ColorDefault
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
	// logger.Printf("PARSED: %v, %v, %v (at %v)", attr, atype, ok, p.index)
	if !ok {
		p.index++
		return TextElement{Type: ElemPrintable, Ch: p.text[p.index-1], Fg: p.currText, Bg: p.currBack}
	}

	if atype == ElemBackColor {
		if attr == ColorDefault {
			p.currBack = p.defBack
		} else {
			p.currBack = attr
		}
	} else if atype == ElemTextColor {
		if attr == ColorDefault {
			p.currText = p.defText
		} else {
			p.currText = attr
		}
	}
	return TextElement{Type: atype, Fg: p.currText, Bg: p.currBack}
}
