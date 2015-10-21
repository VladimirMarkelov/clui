package clui

import (
	"bufio"
	term "github.com/nsf/termbox-go"
	"io/ioutil"
	"os"
	"strings"
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
	current   string
	themePath string
	version   string
}

type ThemeInfo struct {
	parent  string
	title   string
	author  string
	version string
}

/*
A theme structure. It keeps all colors, characters for the theme.
Parent property determines a theme name that is used if a requested
theme object is not declared in the current one. If no parent is
defined then the library uses default built-in theme.
*/
type theme struct {
	parent  string
	title   string
	author  string
	version string
	colors  map[string]term.Attribute
	objects map[string]string
}

const defaultTheme = "default"

func NewThemeManager() *ThemeManager {
	sm := new(ThemeManager)

	sm.Reset()

	return sm
}

func (s *ThemeManager) Reset() {
	s.current = defaultTheme
	s.themes = make(map[string]theme, 0)

	defTheme := theme{parent: "", title: "Default Theme", author: "V. Markelov", version: "1.0"}
	defTheme.colors = make(map[string]term.Attribute, 0)
	defTheme.objects = make(map[string]string, 0)

	defTheme.objects[ObjSingleBorder] = "─│┌┐└┘"
	defTheme.objects[ObjDoubleBorder] = "═║╔╗╚╝"
	defTheme.objects[ObjEdit] = "←→V"
	defTheme.objects[ObjScrollBar] = "|O^v"
	defTheme.objects[ObjViewButtons] = "^↓○[]"
	defTheme.objects[ObjCheckBox] = "[] X?"
	defTheme.objects[ObjRadio] = "() *"
	defTheme.objects[ObjProgressBar] = "░▒"

	defTheme.colors[ColorDisabledText] = ColorBlackBold
	defTheme.colors[ColorDisabledBack] = ColorWhite
	defTheme.colors[ColorText] = ColorWhite
	defTheme.colors[ColorBack] = ColorBlackBold
	defTheme.colors[ColorViewBack] = ColorBlackBold
	defTheme.colors[ColorViewText] = ColorWhite

	defTheme.colors[ColorControlText] = ColorWhite
	defTheme.colors[ColorControlBack] = ColorBlack
	defTheme.colors[ColorControlActiveText] = ColorWhite
	defTheme.colors[ColorControlActiveBack] = ColorMagenta
	defTheme.colors[ColorControlShadow] = ColorBlue
	defTheme.colors[ColorControlDisabledText] = ColorWhite
	defTheme.colors[ColorControlDisabledBack] = ColorBlackBold

	defTheme.colors[ColorEditText] = ColorBlack
	defTheme.colors[ColorEditBack] = ColorWhite
	defTheme.colors[ColorEditActiveText] = ColorBlack
	defTheme.colors[ColorEditActiveBack] = ColorWhiteBold
	defTheme.colors[ColorSelectionText] = ColorYellow
	defTheme.colors[ColorSelectionBack] = ColorBlue

	defTheme.colors[ColorScrollBack] = ColorBlackBold
	defTheme.colors[ColorScrollText] = ColorWhite
	defTheme.colors[ColorThumbBack] = ColorBlackBold
	defTheme.colors[ColorThumbText] = ColorWhite

	defTheme.colors[ColorProgressText] = ColorBlue
	defTheme.colors[ColorProgressBack] = ColorBlackBold
	defTheme.colors[ColorProgressActiveText] = ColorBlack
	defTheme.colors[ColorProgressActiveBack] = ColorBlueBold

	s.themes[defaultTheme] = defTheme
}

func (s *ThemeManager) SysColor(color string) term.Attribute {
	sch, ok := s.themes[s.current]
	if !ok {
		sch = s.themes[defaultTheme]
	}

	clr, okclr := sch.colors[color]
	if !okclr {
		visited := make(map[string]int, 0)
		visited[s.current] = 1
		if !ok {
			visited[defaultTheme] = 1
		}

		for {
			s.LoadTheme(sch.parent)
			sch = s.themes[sch.parent]
			clr, okclr = sch.colors[color]

			if ok {
				break
			} else {
				if _, okSch := visited[sch.parent]; okSch {
					panic("Color + " + color + ". Theme loop detected: " + sch.title + " --> " + sch.parent)
				} else {
					visited[sch.parent] = 1
				}
			}
		}
	}

	return clr
}

func (s *ThemeManager) SysObject(object string) string {
	sch, ok := s.themes[s.current]
	if !ok {
		sch = s.themes[defaultTheme]
	}

	obj, okobj := sch.objects[object]
	if !okobj {
		visited := make(map[string]int, 0)
		visited[s.current] = 1
		if !ok {
			visited[defaultTheme] = 1
		}

		for {
			s.LoadTheme(sch.parent)
			sch = s.themes[sch.parent]
			obj, okobj = sch.objects[object]

			if ok {
				break
			} else {
				if _, okSch := visited[sch.parent]; okSch {
					panic("Object: " + object + ". Theme loop detected: " + sch.title + " --> " + sch.parent)
				} else {
					visited[sch.parent] = 1
				}
			}
		}
	}

	return obj
}

func (s *ThemeManager) ThemeNames() []string {
	var str []string
	str = append(str, defaultTheme)

	path := s.themePath
	if path == "" {
		path = "." + string(os.PathSeparator)
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic("Failed to read theme directory: " + s.themePath)
	}

	for _, f := range files {
		if !f.IsDir() {
			str = append(str, f.Name())
		}
	}

	return str
}

func (s *ThemeManager) CurrentTheme() string {
	return s.current
}

func (s *ThemeManager) SetCurrentTheme(name string) bool {
	if _, ok := s.themes[name]; !ok {
		tnames := s.ThemeNames()
		for _, theme := range tnames {
			if theme == name {
				s.LoadTheme(theme)
				break
			}
		}
	}

	if _, ok := s.themes[name]; ok {
		s.current = name
		return true
	}
	return false
}

func (s *ThemeManager) ThemePath() string {
	return s.themePath
}

func (s *ThemeManager) SetThemePath(path string) {
	if path == s.themePath {
		return
	}

	s.themePath = path
	s.Reset()
}

func (s *ThemeManager) LoadTheme(name string) {
	if _, ok := s.themes[name]; ok {
		return
	}

	theme := theme{parent: "", title: "", author: ""}
	theme.colors = make(map[string]term.Attribute, 0)
	theme.objects = make(map[string]string, 0)

	file, err := os.Open(s.themePath + string(os.PathSeparator) + name)
	if err != nil {
		panic("Failed to open theme " + name + " : " + err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, " ")

		// skip comments
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "/") {
			continue
		}

		// skip invalid lines
		if !strings.Contains(line, "=") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		key := strings.Trim(parts[0], " ")
		value := strings.Trim(parts[1], " ")

		low := strings.ToLower(key)
		if low == "parent" {
			theme.parent = value
		} else if low == "author" {
			theme.author = value
		} else if low == "name" || low == "title" {
			theme.title = value
		} else if low == "version" {
			theme.version = value
		} else if strings.HasSuffix(key, "Back") || strings.HasSuffix(key, "Text") {
			c := StringToColor(value)
			if c%32 == 0 {
				panic("Failed to read color: " + value)
			}
			theme.colors[key] = c
		} else {
			theme.objects[key] = value
		}
	}

	if theme.parent == "" {
		theme.parent = "default"
	}

	s.themes[name] = theme
}

func (s *ThemeManager) ReLoadTheme(name string) {
	if _, ok := s.themes[name]; ok {
		delete(s.themes, name)
	}

	s.LoadTheme(name)
}

func (s *ThemeManager) ThemeInfo(name string) ThemeInfo {
	s.LoadTheme(name)
	var theme ThemeInfo
	if t, ok := s.themes[name]; !ok {
		theme.parent = t.parent
		theme.title = t.title
		theme.version = t.version
	}
	return theme
}
