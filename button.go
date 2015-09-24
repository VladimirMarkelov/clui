package clui

import (
	"github.com/VladimirMarkelov/termbox-go"
	"time"
)

/* Push button control
onClick function is called every time a user clicks the Button. Sent event has only one valid field sender ID(Ctrl field of Event structure).
Button can be clicked with mouse or using space on keyboard when the button is active.
*/
type Button struct {
	posX, posY    int
	width, height int
	title         string
	anchor        Anchor
	id            WinId
	enabled       bool
	align         Align
	active        bool
	pressed       bool
	parent        Window
	visible       bool
	tabStop       bool
	textColor     Color
	backColor     Color
	scale         int

	minW, minH int

	onClick func(Event)
}

func NewButton(parent Window, id WinId, x, y, width, height int, title string, props Props) *Button {
	b := new(Button)
	b.SetEnabled(true)
	b.SetPos(x, y)
	if height < 4 {
		height = 4
	}
	if width < 6 {
		width = 6
	}
	b.SetSize(width, height)
	b.anchor = props.Anchors
	b.pressed = false
	b.active = false
	b.title = title
	b.parent = parent
	b.visible = true
	b.tabStop = true
	b.id = id
	b.minW, b.minH = 6, 4

	return b
}

func (b *Button) SetText(title string) {
	b.title = title
}

func (b *Button) GetText() string {
	return b.title
}

func (b *Button) GetId() WinId {
	return b.id
}

func (b *Button) GetSize() (int, int) {
	return b.width, b.height
}

func (b *Button) GetConstraints() (int, int) {
	return b.minW, b.minH
}

func (b *Button) SetConstraints(minW, minH int) {
	if minW >= 6 {
		b.minW = minW
	}
	if minH >= 4 {
		b.minH = minH
	}
}

func (b *Button) SetSize(width, height int) {
	width, height = ApplyConstraints(b, width, height)

	b.width = width
	b.height = height
}

func (b *Button) GetPos() (int, int) {
	return b.posX, b.posY
}

func (b *Button) SetPos(x, y int) {
	b.posX = x
	b.posY = y
}

func (b *Button) Redraw(canvas Canvas) {
	x, y := b.GetPos()
	w, h := b.GetSize()

	tm := canvas.Theme()

	fg, bg := b.textColor, b.backColor
	shadow := ColorDefault
	if fg == ColorDefault {
		if b.enabled {
			fg = tm.GetSysColor(ColorControlText)
		} else {
			fg = tm.GetSysColor(ColorGrayText)
		}
	}
	if bg == ColorDefault {
		if b.active {
			bg = tm.GetSysColor(ColorControlActiveBack)
		} else {
			bg = tm.GetSysColor(ColorControlBack)
		}
	}
	if shadow == ColorDefault {
		shadow = tm.GetSysColor(ColorControlShadow)
	}

	dy := int((h - 1) / 2)
	if !b.pressed {
		canvas.ClearRect(x+1, y+1, w-1, h-1, shadow)
		canvas.ClearRect(x, y, w-1, h-1, bg)
		canvas.DrawAlignedText(x, y+dy, w-1, b.title, fg, bg, AlignCenter)
	} else {
		canvas.ClearRect(x+1, y+1, w-1, h-1, bg)
		canvas.DrawAlignedText(x+1, y+1+dy, w-1, b.title, fg, bg, AlignCenter)
	}
}

func (b *Button) GetEnabled() bool {
	return b.enabled
}

func (b *Button) SetEnabled(enabled bool) {
	b.enabled = enabled
}

func (b *Button) SetAlign(align Align) {
	// nothing
}

func (b *Button) GetAlign() Align {
	return b.align
}

func (b *Button) SetAnchors(anchor Anchor) {
	b.anchor = anchor
}

func (b *Button) GetAnchors() Anchor {
	return b.anchor
}

func (b *Button) GetActive() bool {
	return b.active
}

func (b *Button) SetActive(active bool) {
	if b.active != active {
		b.active = active
	}
}

func (b *Button) GetTabStop() bool {
	return b.tabStop
}

func (b *Button) SetTabStop(tab bool) {
	b.tabStop = tab
}

func (b *Button) ProcessEvent(event Event) bool {
	if (!b.active && event.Type == EventKey) || !b.enabled || b.pressed {
		return false
	}

	if (event.Type == EventKey && event.Key == termbox.KeySpace) || event.Type == EventMouseClick || event.Type == EventMouse {
		b.pressed = true
		timer := time.NewTimer(time.Millisecond * 150)
		go func() {
			<-timer.C
			b.pressed = false
			// generate ButtonClickEvent
			if b.parent != nil {
				if b.onClick != nil {
					ev := Event{Ctrl: b.id}
					b.onClick(ev)
				}

				ev := InternalEvent{act: EventRedraw, sender: b.id}
				b.parent.SendEvent(ev)
			}
		}()
		return true
	}

	return false
}

func (b *Button) SetVisible(visible bool) {
	b.visible = visible
}

func (b *Button) GetVisible() bool {
	return b.visible
}

func (b *Button) OnClick(fn func(Event)) {
	b.onClick = fn
}

func (b *Button) GetColors() (Color, Color) {
	return b.textColor, b.backColor
}

func (b *Button) SetTextColor(clr Color) {
	b.textColor = clr
}

func (b *Button) SetBackColor(clr Color) {
	b.backColor = clr
}

func (b *Button) HideChildren() {
	// nothing to do
}

func (b *Button) GetScale() int {
	return b.scale
}

func (b *Button) SetScale(scale int) {
	b.scale = scale
}
