package clui

import (
	term "github.com/nsf/termbox-go"
)
//InitLibrary --
func InitLibrary() bool {
	initThemeManager()
	initComposer()
	initMainLoop()
	return initCanvas()
}

//DeinitLibrary Close closes console management and makes a console cursor visible
func DeinitLibrary() {
	term.SetCursor(3, 3)
	term.Close()
}
