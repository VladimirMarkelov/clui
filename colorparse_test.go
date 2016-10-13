package clui

import (
	"testing"
)

func TestParserEmpty(t *testing.T) {
	prs := NewColorParser("", ColorBlack, ColorWhite)

	elem := prs.NextElement()

	if elem.Type != ElemEndOfText {
		t.Errorf("Empty string must return end of text (%v vs %v)",
			ElemEndOfText, elem.Type)
	}
}

func TestParserColors(t *testing.T) {
	prs := NewColorParser("a<b:green>c<t:red>d<b:>e<t:>fg\nf",
		ColorBlack, ColorWhite)
	elems := []TextElement{
		{ElemPrintable, 'a', ColorBlack, ColorWhite},
		{ElemBackColor, ' ', ColorBlack, ColorGreen},
		{ElemPrintable, 'c', ColorBlack, ColorGreen},
		{ElemTextColor, 'c', ColorRed, ColorGreen},
		{ElemPrintable, 'd', ColorRed, ColorGreen},
		{ElemBackColor, 'd', ColorRed, ColorWhite},
		{ElemPrintable, 'e', ColorRed, ColorWhite},
		{ElemTextColor, 'e', ColorBlack, ColorWhite},
		{ElemPrintable, 'f', ColorBlack, ColorWhite},
		{ElemPrintable, 'g', ColorBlack, ColorWhite},
		{ElemLineBreak, 'g', ColorBlack, ColorWhite},
		{ElemPrintable, 'f', ColorBlack, ColorWhite},
	}

	idx := 0
	el := prs.NextElement()

	for el.Type != ElemEndOfText {
		if idx >= len(elems) {
			t.Errorf("Size mismatch: string must have only %v items", len(elems))
		}

		if el.Type != elems[idx].Type ||
			(el.Type == ElemPrintable && (el.Ch != elems[idx].Ch || el.Fg != elems[idx].Fg || el.Bg != elems[idx].Bg)) ||
			(el.Type == ElemTextColor && el.Fg != elems[idx].Fg) ||
			(el.Type == ElemBackColor && el.Bg != elems[idx].Bg) {
			t.Errorf("Elements mismatch at %v: {%v, %v, %v, %v} = {%v, %v, %v, %v}",
				idx, el.Type, elems[idx].Type,
				el.Ch, elems[idx].Ch,
				el.Fg, elems[idx].Fg,
				el.Bg, elems[idx].Bg)
		}

		el = prs.NextElement()
		idx += 1
	}
}
