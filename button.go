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
	BaseControl
	shadowColor term.Attribute
	bgActive    term.Attribute
	pressed     bool
	onClick     func(Event)
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
func CreateButton(parent Control, width, height int, title string, scale int) *Button {
	b := new(Button)

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
	b.SetScale(scale)

	if parent != nil {
		parent.AddChild(b)
	}

	return b
}

// Repaint draws the control on its View surface
func (b *Button) Draw() {
	PushAttributes()
	defer PopAttributes()

	x, y := b.Pos()
	w, h := b.Size()

	fg, bg := b.fg, b.bg
	shadow := RealColor(b.shadowColor, ColorButtonShadow)
	if !b.Enabled() {
		fg, bg = RealColor(fg, ColorButtonDisabledText), RealColor(bg, ColorButtonDisabledBack)
	} else if b.Active() {
		fg, bg = RealColor(b.fgActive, ColorButtonActiveText), RealColor(b.bgActive, ColorButtonActiveBack)
	} else {
		fg, bg = RealColor(fg, ColorButtonText), RealColor(bg, ColorButtonBack)
	}

	dy := int((h - 1) / 2)
	SetTextColor(fg)
	shift, text := AlignColorizedText(b.title, w-1, b.align)
	if !b.pressed {
		SetBackColor(shadow)
		FillRect(x+1, y+1, w-1, h-1, ' ')
		SetBackColor(bg)
		FillRect(x, y, w-1, h-1, ' ')
		DrawText(x+shift, y+dy, text)
	} else {
		SetBackColor(bg)
		FillRect(x+1, y+1, w-1, h-1, ' ')
		DrawText(x+1+shift, y+1+dy, b.title)
	}
}

/*
ProcessEvent processes all events come from the control parent. If a control
processes an event it should return true. If the method returns false it means
that the control do not want or cannot process the event and the caller sends
the event to the control parent
*/
func (b *Button) ProcessEvent(event Event) bool {
	if !b.Enabled() {
		return false
	}

	if event.Type == EventKey {
		if event.Key == term.KeySpace && !b.pressed {
			b.pressed = true
			ev := Event{Type: EventRedraw}
			go func() {
				b.Draw()
				PutEvent(ev)
				time.Sleep(100 * time.Millisecond)
				b.pressed = false
				b.Draw()
				PutEvent(ev)
			}()
			if b.onClick != nil {
				b.onClick(event)
			}
			return true
		} else if event.Key == term.KeyEsc && b.pressed {
			b.pressed = false
			ReleaseEvents()
			return true
		}
	} else if event.Type == EventMouse {
		if event.Key == term.MouseLeft {
			b.pressed = true
			GrabEvents(b)
			return true
		} else if event.Key == term.MouseRelease && b.pressed {
			ReleaseEvents()
			if event.X >= b.x && event.Y >= b.y && event.X < b.x+b.width && event.Y < b.y+b.height {
				if b.onClick != nil {
					b.onClick(event)
				}
			}
			b.pressed = false
			return true
		}
	}

	return false
}

// OnClick sets the callback that is called when one clicks button
// with mouse or pressing space on keyboard while the button is active
func (b *Button) OnClick(fn func(Event)) {
	b.onClick = fn
}
