package clui

import (
	"testing"
	мКнст "./пакКонстанты"
)

func TestParserEmpty(t *testing.T) {
	prs := NewColorParser("", мКнст.ColorBlack, мКнст.ColorWhite)

	elem := prs.NextElement()

	if elem.Type != ElemEndOfText {
		t.Errorf("Empty string must return end of text (%v vs %v)",
			ElemEndOfText, elem.Type)
	}
}

func TestParserColors(t *testing.T) {
	prs := NewColorParser("a<b:green>c<t:red>d<b:>e<t:>fg\nf",
		мКнст.ColorBlack, мКнст.ColorWhite)
	elems := []TextElement{
		{ElemPrintable, 'a', мКнст.ColorBlack, мКнст.ColorWhite},
		{ElemBackColor, ' ', мКнст.ColorBlack, мКнст.ColorGreen},
		{ElemPrintable, 'c', мКнст.ColorBlack, мКнст.ColorGreen},
		{ElemTextColor, 'c', мКнст.ColorRed, мКнст.ColorGreen},
		{ElemPrintable, 'd', мКнст.ColorRed, мКнст.ColorGreen},
		{ElemBackColor, 'd', мКнст.ColorRed, мКнст.ColorWhite},
		{ElemPrintable, 'e', мКнст.ColorRed, мКнст.ColorWhite},
		{ElemTextColor, 'e', мКнст.ColorBlack, мКнст.ColorWhite},
		{ElemPrintable, 'f', мКнст.ColorBlack, мКнст.ColorWhite},
		{ElemPrintable, 'g', мКнст.ColorBlack, мКнст.ColorWhite},
		{ElemLineBreak, 'g', мКнст.ColorBlack, мКнст.ColorWhite},
		{ElemPrintable, 'f', мКнст.ColorBlack, мКнст.ColorWhite},
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
		idx++
	}
}
