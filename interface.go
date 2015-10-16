package clui

import (
	term "github.com/nsf/termbox-go"
	"log"
)

type Screen interface {
	Theme() Theme
	PutEvent(Event)

	Logger() *log.Logger
}

type Canvas interface {
	SetSize(int, int)
	Size() (int, int)
	PutSymbol(int, int, term.Cell) bool
	PutText(int, int, string, term.Attribute, term.Attribute)
	Symbol(int, int) (term.Cell, bool)
	Clear(term.Attribute)
	FillRect(int, int, int, int, term.Cell)
	DrawFrame(int, int, int, int, term.Attribute, term.Attribute, string)
	SetCursorPos(int, int)
}

type Theme interface {
	SysObject(ObjId) string
	SysColor(ColorId) term.Attribute
}

type View interface {
	Title() string
	SetTitle(string)
	Draw(Canvas)
	Repaint()
	Constraints() (int, int)
	Size() (int, int)
	SetSize(int, int)
	Pos() (int, int)
	SetPos(int, int)
	Canvas() Canvas
	Active() bool
	SetActive(bool)
	ProcessEvent(Event) bool
	ActivateControl(Control)
	RegisterControl(Control)
	Screen() Screen
	Parent() Control
	HitTest(int, int) HitResult

	Paddings() (int, int, int, int)
	AddChild(Control, int)
	SetPack(PackType)
	Pack() PackType
	Children() []Control
	ChildExists(Control) bool
	Scale() int
	SetScale(int)
	TabStop() bool
	Colors() (term.Attribute, term.Attribute)

	Logger() *log.Logger
}

type Control interface {
	Title() string
	SetTitle(string)
	Pos() (int, int)
	SetPos(int, int)
	Size() (int, int)
	SetSize(int, int)
	Scale() int
	SetScale(int)
	Constraints() (int, int)
	Paddings() (int, int, int, int)
	Repaint()
	AddChild(Control, int)
	SetPack(PackType)
	Pack() PackType
	Children() []Control
	Active() bool
	SetActive(bool)
	ProcessEvent(Event) bool
	TabStop() bool
	Parent() Control
	Colors() (term.Attribute, term.Attribute)

	Logger() *log.Logger
}
