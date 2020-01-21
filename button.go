package clui

import (
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
	"sync/atomic"
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
	pressed     int32
	shadowType  ButtonShadow
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
	b.BaseControl = NewBaseControl()

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
	shadow := RealColor(b.shadowColor, b.Style(), ColorButtonShadow)
	if b.disabled {
		fg, bg = RealColor(fg, b.Style(), ColorButtonDisabledText), RealColor(bg, b.Style(), ColorButtonDisabledBack)
	} else if b.Active() {
		fg, bg = RealColor(b.fgActive, b.Style(), ColorButtonActiveText), RealColor(b.bgActive, b.Style(), ColorButtonActiveBack)
	} else {
		fg, bg = RealColor(fg, b.Style(), ColorButtonText), RealColor(bg, b.Style(), ColorButtonBack)
	}

	dy := int((h - 1) / 2)
	SetTextColor(fg)
	shift, text := AlignColorizedText(b.title, w-1, b.align)
	if b.isPressed() == 0 {
		switch b.shadowType {
		case ShadowFull:
			SetBackColor(shadow)
			FillRect(x+1, y+h-1, w-1, 1, ' ')
			FillRect(x+w-1, y+1, 1, h-1, ' ')
		case ShadowHalf:
			parts := []rune(SysObject(ObjButton))
			var bottomCh, rightCh rune
			if len(parts) < 2 {
				bottomCh, rightCh = '▀', '█'
			} else {
				bottomCh, rightCh = parts[0], parts[1]
			}
			SetTextColor(shadow)
			FillRect(x+1, y+h-1, w-1, 1, bottomCh)
			FillRect(x+w-1, y+1, 1, h-2, rightCh)
		}
		SetTextColor(fg)
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
func (b *Button) ProcessEvent(event Event) bool {
	if !b.Enabled() {
		return false
	}

	if event.Type == EventKey {
		if event.Key == term.KeySpace && b.isPressed() == 0 {
			b.setPressed(1)
			ev := Event{Type: EventRedraw}

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
	} else if event.Type == EventMouse {
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
func (b *Button) OnClick(fn func(Event)) {
	b.onClick = fn
}

// ShadowType returns type of a show the button drops
func (b *Button) ShadowType() ButtonShadow {
	return b.shadowType
}

// SetShadowType changes the shadow the button drops
func (b *Button) SetShadowType(sh ButtonShadow) {
	b.mtx.Lock()
	b.shadowType = sh
	b.mtx.Unlock()
}
