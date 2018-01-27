// +build darwin

package clui

import (
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
	"time"
)

const charInvervalMs = 20

/*
EditField is a single-line text edit contol. Edit field consumes some keyboard
events when it is active: all printable charaters; Delete, BackSpace, Home,
End, left and right arrows; Ctrl+R to clear EditField.
Edit text can be limited. By default a user can enter text of any length.
Use SetMaxWidth to limit the maximum text length. If the text is longer than
maximun then the text is automatically truncated.
EditField calls onChage in case of its text is changed. Event field Msg contains the new text
*/
type EditField struct {
	BaseControl
	// cursor position in edit text
	cursorPos int
	// the number of the first displayed text character - it is used in case of text is longer than edit width
	offset    int
	readonly  bool
	maxWidth  int
	showStars bool

	onChange   func(Event)
	onKeyPress func(term.Key) bool

	lastEvent time.Time
}

// NewEditField creates a new EditField control
// view - is a View that manages the control
// parent - is container that keeps the control. The same View can be a view and a parent at the same time.
// width - is minimal width of the control.
// text - text to edit.
// scale - the way of scaling the control when the parent is resized. Use DoNotScale constant if the
//  control should keep its original size.
func CreateEditField(parent Control, width int, text string, scale int) *EditField {
	e := new(EditField)
	e.onChange = nil
	e.SetTitle(text)
	e.SetEnabled(true)

	if width == AutoSize {
		width = xs.Len(text) + 1
	}

	e.SetSize(width, 1)
	e.cursorPos = xs.Len(text)
	e.offset = 0
	e.parent = parent
	e.readonly = false
	e.SetScale(scale)
	e.lastEvent = time.Now()

	e.SetConstraints(width, 1)

	e.end()

	if parent != nil {
		parent.AddChild(e)
	}

	return e
}

/*
ProcessEvent processes all events come from the control parent. If a control
processes an event it should return true. If the method returns false it means
that the control do not want or cannot process the event and the caller sends
the event to the control parent
*/
func (e *EditField) ProcessEvent(event Event) bool {
	if !e.Active() || !e.Enabled() {
		return false
	}

	if event.Type == EventActivate && event.X == 0 {
		term.HideCursor()
	}

	if event.Type == EventMouse && event.Key == term.MouseLeft {
		e.lastEvent = time.Now()
	}

	if event.Type == EventKey && event.Key != term.KeyTab {
		if e.onKeyPress != nil {
			res := e.onKeyPress(event.Key)
			if res {
				return true
			}
		}

		switch event.Key {
		case term.KeyEnter:
			return false
		case term.KeySpace:
			e.insertRune(' ')
			return true
		case term.KeyBackspace, term.KeyBackspace2:
			e.backspace()
			return true
		case term.KeyDelete:
			e.del()
			return true
		case term.KeyArrowLeft:
			e.charLeft()
			return true
		case term.KeyHome:
			e.home()
			return true
		case term.KeyEnd:
			e.end()
			return true
		case term.KeyCtrlR:
			if !e.readonly {
				e.Clear()
			}
			return true
		case term.KeyArrowRight:
			e.charRight()
			return true
		default:
			if event.Ch != 0 {
				elapsed := time.Now().Sub(e.lastEvent)
				if elapsed > time.Duration(charInvervalMs)*time.Millisecond {
					e.insertRune(event.Ch)
					e.lastEvent = time.Now()
				}
				return true
			}
		}
		return false
	}

	return false
}
