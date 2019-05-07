package пакИнтерфейсы

/*
	Модуль предоставляет базовый виджет для всех виджетов
*/
import (
	term "github.com/nsf/termbox-go"
)

//ИВиджет -- базовая абстракция для всех виджетов
type ИВиджет interface {
	Title() string
	SetTitle(title string)
	Size() (widht int, height int)
	SetSize(width, height int)
	Pos() (x int, y int)
	SetPos(x, y int)
	Constraints() (minw int, minh int)
	SetConstraints(minw, minh int)
	Active() bool
	SetActive(active bool)
	TabStop() bool
	SetTabStop(tabstop bool)
	Enabled() bool
	SetEnabled(enabled bool)
	Visible() bool
	SetVisible(enabled bool)
	Parent() ИВиджет
	SetParent(parent ИВиджет)
	Modal() bool
	SetModal(modal bool)
	Paddings() (px int, py int)
	SetPaddings(px, py int)
	Gaps() (dx int, dy int)
	SetGaps(dx, dy int)
	Pack() PackType
	SetPack(pack PackType)
	Scale() int
	SetScale(scale int)
	Align() Align
	SetAlign(align Align)
	TextColor() term.Attribute
	SetTextColor(clr term.Attribute)
	BackColor() term.Attribute
	SetBackColor(clr term.Attribute)
	ActiveColors() (term.Attribute, term.Attribute)
	SetActiveBackColor(term.Attribute)
	SetActiveTextColor(term.Attribute)
	AddChild(control ИВиджет)
	Children() []ИВиджет
	ChildExists(control ИВиджет) bool
	MinimalSize() (w int, h int)
	ChildrenScale() int
	ResizeChildren()
	PlaceChildren()
	Draw()
	DrawChildren()
	HitTest(x, y int) HitResult
	ProcessEvent(ev ИСобытие) bool
	RefID() int64
	RemoveChild(control ИВиджет)
	Destroy()
	Style() string
	SetStyle(style string)
	Clipped() bool
	SetClipped(clipped bool)
	Clipper() (int, int, int, int)

}
