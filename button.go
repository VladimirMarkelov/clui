package clui

import (
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
	"sync/atomic"
	"time"
	мКнст "./пакКонстанты"
	мИнт "./пакИнтерфейсы"
)

/*
Button is a simpe push button control. Every time a user clicks a Button, it
emits OnClick event. Event has only one valid field Sender.
Button can be clicked with mouse or using space on keyboard while the Button is active.
*/
type Button struct {
	*BaseControl
	shadowColor term.Attribute
	bgActive    term.Attribute
	pressed     int32
	onClick     func(мИнт.ИСобытие)
}

/*
CreateButton creates a new Button.
view - is a View that manages the control
parent - is container that keeps the control. The same View can be a view and a parent at the same time.
width and heigth - are minimal size of the control.
title - button title.
scale - the way of scaling the control when the parent is resized. Use DoNotScale constant if the
control should keep its original size.
*/
func CreateButton(parent мИнт.ИВиджет, width, height int, title string, scale int) *Button {
	b := new(Button)
	b.BaseControl = NewBaseControl()

	b.parent = parent
	b.align = мКнст.AlignCenter

	if height == мКнст.AutoSize {
		height = 4
	}
	if width == мКнст.AutoSize {
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

//Draw Repaint draws the control on its View surface
func (b *Button) Draw() {
	if b.hidden {
		return
	}

	b.mtx.RLock()
	defer b.mtx.RUnlock()
	PushAttributes()
	defer PopAttributes()

	x, y := b.Pos()
	w, h := b.Size()

	fg, bg := b.fg, b.bg
	shadow := RealColor(b.shadowColor, b.Style(), мКнст.ColorButtonShadow)
	if b.disabled {
		fg, bg = RealColor(fg, b.Style(), мКнст.ColorButtonDisabledText), RealColor(bg, b.Style(), мКнст.ColorButtonDisabledBack)
	} else if b.Active() {
		fg, bg = RealColor(b.fgActive, b.Style(), мКнст.ColorButtonActiveText), RealColor(b.bgActive, b.Style(), мКнст.ColorButtonActiveBack)
	} else {
		fg, bg = RealColor(fg, b.Style(), мКнст.ColorButtonText), RealColor(bg, b.Style(), мКнст.ColorButtonBack)
	}

	dy := int((h - 1) / 2)
	SetTextColor(fg)
	shift, text := AlignColorizedText(b.title, w-1, b.align)
	if b.isPressed() == 0 {
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

func (b *Button) isPressed() int32 {
	return atomic.LoadInt32(&b.pressed)
}

func (b *Button) setPressed(pressed int32) {
	atomic.StoreInt32(&b.pressed, pressed)
}

/*
ProcessEvent processes all events come from the control parent. If a control
processes an event it should return true. If the method returns false it means
that the control do not want or cannot process the event and the caller sends
the event to the control parent
*/
func (b *Button) ProcessEvent(event мИнт.ИСобытие) bool {
	if !b.Enabled() {
		return false
	}

	if event.Type == мКнст.EventKey {
		if event.Key == term.KeySpace && b.isPressed() == 0 {
			b.setPressed(1)
			ev := мКнст.Event{Type: мКнст.EventRedraw}

			go func() {
				PutEvent(ev)
				time.Sleep(100 * time.Millisecond)
				b.setPressed(0)
				PutEvent(ev)
			}()

			if b.onClick != nil {
				b.onClick(event)
			}
			return true
		} else if event.Key == term.KeyEsc && b.isPressed() != 0 {
			b.setPressed(0)
			ReleaseEvents()
			return true
		}
	} else if event.Type == мКнст.EventMouse {
		if event.Key == term.MouseLeft {
			b.setPressed(1)
			GrabEvents(b)
			return true
		} else if event.Key == term.MouseRelease && b.isPressed() != 0 {
			ReleaseEvents()
			if event.X >= b.x && event.Y >= b.y && event.X < b.x+b.width && event.Y < b.y+b.height {
				if b.onClick != nil {
					b.onClick(event)
				}
			}
			b.setPressed(0)
			return true
		}
	}

	return false
}

// OnClick sets the callback that is called when one clicks button
// with mouse or pressing space on keyboard while the button is active
func (b *Button) OnClick(fn func(мИнт.ИСобытие)) {
	b.onClick = fn
}
