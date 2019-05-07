package clui

import (
	"testing"
	мИнт "./пакИнтерфейсы"
)

func TestParserEmpty(t *testing.T) {
	prs := NewColorParser("", мИнт.ColorBlack, мИнт.ColorWhite)

	elem := prs.NextElement()

	if elem.Type != ElemEndOfText {
		t.Errorf("Empty string must return end of text (%v vs %v)",
			ElemEndOfText, elem.Type)
	}
}

func TestParserColors(t *testing.T) {
	prs := NewColorParser("a<b:green>c<t:red>d<b:>e<t:>fg\nf",
		мИнт.ColorBlack, мИнт.ColorWhite)
	elems := []TextElement{
		{ElemPrintable, 'a', мИнт.ColorBlack, мИнт.ColorWhite},
		{ElemBackColor, ' ', мИнт.ColorBlack, мИнт.ColorGreen},
		{ElemPrintable, 'c', мИнт.ColorBlack, мИнт.ColorGreen},
		{ElemTextColor, 'c', мИнт.ColorRed, мИнт.ColorGreen},
		{ElemPrintable, 'd', мИнт.ColorRed, мИнт.ColorGreen},
		{ElemBackColor, 'd', мИнт.ColorRed, мИнт.ColorWhite},
		{ElemPrintable, 'e', мИнт.ColorRed, мИнт.ColorWhite},
		{ElemTextColor, 'e', мИнт.ColorBlack, мИнт.ColorWhite},
		{ElemPrintable, 'f', мИнт.ColorBlack, мИнт.ColorWhite},
		{ElemPrintable, 'g', мИнт.ColorBlack, мИнт.ColorWhite},
		{ElemLineBreak, 'g', мИнт.ColorBlack, мИнт.ColorWhite},
		{ElemPrintable, 'f', мИнт.ColorBlack, мИнт.ColorWhite},
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
