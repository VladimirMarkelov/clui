package clui

import (
	term "github.com/nsf/termbox-go"
)

/*
Theme support for controls.
The current implementation is limited but later the manager will be
able to load a requested theme on demand and use deep inheritance.
Theme 'default' exists always - it is predefinded and always complete.
User-defined themes may omit any theme section, all omitted items
are loaded from parent theme. The only required property that a user-
defined theme must have is a theme name.
*/
type ThemeManager struct {
	// available theme list
	themes map[string]theme
	// name of the current theme
	current string
}

/*
A theme structure. It keeps all colors, characters for the theme.
Parent property determines a theme name that is used if a requested
theme object is not declared in the current one. If no parent is
defined then the library uses default built-in theme.
*/
type theme struct {
	colors  map[ColorId]term.Attribute
	objects map[ObjId]string
	parent  string
}

const defaultTheme = "default"

// Theme color constants
// const (
// )

func NewThemeManager() *ThemeManager {
	sm := new(ThemeManager)
	sm.current = defaultTheme
	sm.themes = make(map[string]theme, 0)

	defTheme := theme{parent: ""}
	defTheme.colors = make(map[ColorId]term.Attribute, 0)
	defTheme.objects = make(map[ObjId]string, 0)

	defTheme.objects[ObjSingleBorder] = "─│┌┐└┘"
	defTheme.objects[ObjDoubleBorder] = "═║╔╗╚╝"
	defTheme.objects[ObjEdit] = "←→V"
	defTheme.objects[ObjScrollBar] = "|O^v"
	defTheme.objects[ObjViewButtons] = "^↓○[]"
	defTheme.objects[ObjCheckBox] = "[] X?"
	defTheme.objects[ObjRadio] = "() *"
	defTheme.objects[ObjProgressBar] = "░▒"

	defTheme.colors[ColorText] = ColorWhite
	defTheme.colors[ColorBack] = ColorBlack
	defTheme.colors[ColorViewBack] = ColorBlackBold
	defTheme.colors[ColorViewText] = ColorWhite

	sm.themes[defaultTheme] = defTheme

	return sm
}

func (s *ThemeManager) SysColor(color ColorId) term.Attribute {
	sch, ok := s.themes[s.current]
	if !ok {
		sch = s.themes[defaultTheme]
	}

	clr, okclr := sch.colors[color]
	if !okclr && sch.parent != "" {
		sch = s.themes[sch.parent]
		clr, okclr = sch.colors[color]
		if !okclr {
			clr = term.ColorDefault
		}
	}

	return clr
}

func (s *ThemeManager) SysObject(object ObjId) string {
	sch, ok := s.themes[s.current]
	if !ok {
		sch = s.themes[defaultTheme]
	}

	obj, okobj := sch.objects[object]
	if !okobj && sch.parent != "" {
		sch = s.themes[sch.parent]
		obj = sch.objects[object]
	}

	return obj
}

func (s *ThemeManager) GetThemeList() []string {
	str := make([]string, len(s.themes))
	for k := range s.themes {
		str = append(str, k)
	}

	return str
}

func (s *ThemeManager) GetCurrentTheme() string {
	return s.current
}

func (s *ThemeManager) SetCurrentTheme(name string) bool {
	if _, ok := s.themes[name]; ok {
		s.current = name
		return true
	}
	return false
}
