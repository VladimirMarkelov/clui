package clui

import (
	term "github.com/nsf/termbox-go"
)

func InitLibrary() bool {
	initThemeManager()
	initComposer()
	initMainLoop()
	return initCanvas()
}

// Close closes console management and makes a console cursor visible
func DeinitLibrary() {
	term.SetCursor(3, 3)
	term.Close()
}
