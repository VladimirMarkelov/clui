package clui

import (
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
	"time"
)

/*
Button is a simpe push button control. Every time a user clicks a Button, it
emits OnClick event. Event has only one valid field Sender.
Button can be clicked with mouse or using space on keyboard while the Button is active.
*/
type Button struct {
	ControlBase
	pressed     bool
	shadowColor term.Attribute

	onClick func(Event)
}

/*
NewButton creates a new Button.
view - is a View that manages the control
parent - is container that keeps the control. The same View can be a view and a parent at the same time.
width and heigth - are minimal size of the control.
title - button title.
scale - the way of scaling the control when the parent is resized. Use DoNotScale constant if the
control should keep its original size.
*/
func NewButton(view View, parent Control, width, height int, title string, scale int) *Button {
	b := new(Button)

	b.view = view
	b.parent = parent
	b.align = AlignCenter

	if height == AutoSize {
		height = 4
	}
	if width == AutoSize {
		width = xs.Len(title) + 2 + 1
	}

	if height < 4 {
		height = 4
	}
	if width < 6 {
		width = 6
	}

	b.SetTitle(title)
	b.SetSize(width, height)
	b.SetConstraints(width, height)

	if parent != nil {
		parent.AddChild(b, scale)
	}

	return b
}

// Repaint draws the control on its View surface
func (b *Button) Repaint() {
	x, y := b.Pos()
	w, h := b.Size()
	canvas := b.view.Canvas()
	tm := b.view.Screen().Theme()

	fg, bg := b.fg, b.bg
	shadow := RealColor(tm, b.shadowColor, ColorButtonShadow)
	if !b.Enabled() {
		fg, bg = RealColor(tm, fg, ColorButtonDisabledText), RealColor(tm, bg, ColorButtonDisabledBack)
	} else if b.Active() {
		fg, bg = RealColor(tm, b.fgActive, ColorButtonActiveText), RealColor(tm, b.bgActive, ColorButtonActiveBack)
	} else {
		fg, bg = RealColor(tm, fg, ColorButtonText), RealColor(tm, bg, ColorButtonBack)
	}

	dy := int((h - 1) / 2)
	shift, text := AlignText(b.title, w-1, b.align)
	if !b.pressed {
		canvas.FillRect(x+1, y+1, w-1, h-1, term.Cell{Ch: ' ', Bg: shadow})
		canvas.FillRect(x, y, w-1, h-1, term.Cell{Ch: ' ', Bg: bg})
		canvas.PutText(x+shift, y+dy, text, fg, bg)
	} else {
		canvas.FillRect(x+1, y+1, w-1, h-1, term.Cell{Ch: ' ', Bg: bg})
		canvas.PutText(x+1+shift, y+1+dy, b.title, fg, bg)
	}
}

/*
ProcessEvent processes all events come from the control parent. If a control
processes an event it should return true. If the method returns false it means
that the control do not want or cannot process the event and the caller sends
the event to the control parent
*/
func (b *Button) ProcessEvent(event Event) bool {
	if (!b.active && event.Type == EventKey) || !b.Enabled() || b.pressed {
		return false
	}

	if (event.Type == EventKey && event.Key == term.KeySpace) || event.Type == EventMouse {
		b.pressed = true
		timer := time.NewTimer(time.Millisecond * 150)
		go func() {
			<-timer.C
			b.pressed = false
			// generate ButtonClickEvent
			if b.parent != nil {
				if b.onClick != nil {
					ev := Event{Sender: b}
					b.onClick(ev)
				}

				ev := Event{Type: EventRedraw, Sender: b}
				b.view.Screen().PutEvent(ev)
			}
		}()
		return true
	}

	return false
}

// OnClick sets the callback that is called when one clicks button
// with mouse or pressing space on keyboard while the button is active
func (b *Button) OnClick(fn func(Event)) {
	b.onClick = fn
}
