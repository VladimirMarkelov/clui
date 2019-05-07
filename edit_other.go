// +build !darwin

package clui

import (
	"github.com/atotto/clipboard"
	xs "github.com/huandu/xstrings"
	term "github.com/nsf/termbox-go"
	мКнст "./пакКонстанты"
)

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

	onChange   func(мКнст.Event)
	onKeyPress func(term.Key, rune) bool
}

// CreateEditField creates a new EditField control
// view - is a View that manages the control
// parent - is container that keeps the control. The same View can be a view and a parent at the same time.
// width - is minimal width of the control.
// text - text to edit.
// scale - the way of scaling the control when the parent is resized. Use DoNotScale constant if the
//  control should keep its original size.
func CreateEditField(parent Control, width int, text string, scale int) *EditField {
	e := new(EditField)
	e.BaseControl = NewBaseControl()
	e.onChange = nil
	e.SetTitle(text)
	e.SetEnabled(true)

	if width == мКнст.AutoSize {
		width = xs.Len(text) + 1
	}

	e.SetSize(width, 1)
	e.cursorPos = xs.Len(text)
	e.offset = 0
	e.parent = parent
	e.readonly = false
	e.SetScale(scale)

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
func (e *EditField) ProcessEvent(event мКнст.Event) bool {
	if !e.Active() || !e.Enabled() {
		return false
	}

	if event.Type == мКнст.EventActivate && event.X == 0 {
		term.HideCursor()
	}

	if event.Type == мКнст.EventKey && event.Key != term.KeyTab {
		if e.onKeyPress != nil {
			res := e.onKeyPress(event.Key, event.Ch)
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
		case term.KeyCtrlC:
			if !e.showStars {
				clipboard.WriteAll(e.Title())
			}
			return true
		case term.KeyCtrlV:
			if !e.readonly {
				s, _ := clipboard.ReadAll()
				e.SetTitle(s)
				e.end()
			}
			return true
		default:
			if event.Ch != 0 {
				e.insertRune(event.Ch)
				return true
			}
		}
		return false
	}

	return false
}
