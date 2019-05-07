package пакСобытия

import (
	мИнт "../пакИнтерфейсы"
	term "github.com/nsf/termbox-go"
)

// Event is structure used by Views and controls to communicate with Composer
// and vice versa
type Event struct {
	// Event type - the first events are mapped to termbox Event and then a few
	// own events added to the end
	_Type мИнт.EventType
	// Mod - is a key modifier. Only Alt modifier is supported
	_Mod term.Modifier
	// Msg is a text part of the event. Used by few events: e.g, ListBox click
	// sends a value of clicked item
	_Msg string
	// X and Y are multi-purpose fields: mouse coordinated for click event,
	// X is used to indicate on/off for events like Activate
	// Y is used for vertical-based events like ListBox item selection - id of the item
	_X, _Y int
	// Err is error got from termbox library
	_Err error
	// Key is a pressed key
	_Key term.Key
	// Ch is a printable representation of pressed key combinaton
	_Ch rune
	// For resize event - new terminal size
	_Width  int
	_Height int
	_Target мИнт.ИВиджет
}

//СобытиеНов - -возвращает ссылку на новый Event
func СобытиеНов(ev *term.Event) (e *Event) {
	e = &Event{}
	e._Type = мИнт.EventType(ev.Type)
	e._Ch = ev.Ch
	e._Key = ev.Key
	e._Err = ev.Err
	e._X = ev.MouseX
	e._Y = ev.MouseY
	e._Mod = ev.Mod
	e._Width = ev.Width
	e._Height = ev.Height
	return e
}

//TypeSet -- устанавилвает тип события
func (сам *Event) TypeSet(пТип мИнт.EventType) {
	сам._Type = пТип
}

//Ch -- возвращает руну чего-то там
func (сам *Event) Ch() rune {
	return сам._Ch
}

//Err -- возвращает ошибку события
func (сам *Event) Err() error {
	return сам._Err
}

//Height -- возвращает высоту события
func (сам *Event) Height() int {
	return сам._Height
}

//Key -- возвращает клавишу события
func (сам *Event) Key() term.Key {
	return сам._Key
}

//KeySet -- устанавлиает клавишу события
func (сам *Event) KeySet(пKey term.Key) {
	сам._Key = пKey
}

//Mod -- возвращает модификатор события
func (сам *Event) Mod() term.Modifier {
	return сам._Mod
}

//Target -- возвращает цель события
func (сам *Event) Target() мИнт.ИВиджет {
	return сам._Target
}

//TargetSet -- устанавлиает цель события
func (сам *Event) TargetSet(пЦель мИнт.ИВиджет) {
	сам._Target = пЦель
}

//Msg -- возвращает сообщение события
func (сам *Event) Msg() string {
	return сам._Msg
}

//MsgSet -- устанавливает сообщение события
func (сам *Event) MsgSet(пСбщ string) {
	сам._Msg = пСбщ
}

//Type -- возвращает тип события
func (сам *Event) Type() мИнт.EventType {
	return сам._Type
}

//Width -- возвращает ширину события
func (сам *Event) Width() int {
	return сам._Width
}

//X -- возвращает X события
func (сам *Event) X() int {
	return сам._X
}

//SetX -- устанавливает X события
func (сам *Event) SetX(пХ int) {
	сам._X = пХ
}

//Y -- возвращает Y события
func (сам *Event) Y() int {
	return сам._Y
}

//SetY -- устанавливает X события
func (сам *Event) SetY(пY int) {
	сам._Y = пY
}
