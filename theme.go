package clui

import (
	"bufio"
	"fmt"
	term "github.com/nsf/termbox-go"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"unicode/utf8"
)

/*
ThemeManager support for controls.
The current implementation is limited but later the manager will be
able to load a requested theme on demand and use deep inheritance.
Theme 'default' exists always - it is predefinded and always complete.
User-defined themes may omit any theme section, all omitted items
are loaded from parent theme. The only required property that a user-
defined theme must have is a theme name.

Theme file is a simple text file that has similar to INI file format:
1. Every line started with '#' or '/' is a comment line.
2. Invalid lines - lines that do not contain symbol '=' - are skipped.
3. Valid lines are splitted in two parts:
    key - the text before the first '=' in the line
    value - the text after the first '=' in the line (so, values can
        include '=')
    key and value are trimmed - spaces are removed from both ends.
    If line starts and ends with quote or double quote symbol then
    these symbols are removed, too. It is done to be able to start
    or finish the object with a space rune
4. There is no mandatory keys - all of them are optional
5. Avaiable system keys that used to describe the theme:
    'title' - the theme title
    'author' - theme author
    'version' - theme version
    'parent' - name of the parent theme. If it is not set then the
        'default' is used as a parent
6. Non-system keys are divided into two groups: Colors and Objects
    Colors are the keys that end with 'Back' or 'Text' - background
        and text color, respectively. If theme manager cannot
        value to color it uses black color. See Color*Back * Color*Text
        constants, just drop 'Color' at the beginning of key name.
        Rules of converting text to color:
        1. If the value does not end neither with 'Back' nor with 'Text'
            it is considered as raw attribute value(e.g, 'green bold')
        2. If the value ends with 'Back' or 'Text' it means that one
            of earlier defined attribute must be used. If the current
            scheme does not have that attribute defined (e.g, it is
            defined later in file) then parent theme attribute with
            the same name is used. One can force using parent theme
            colors - just add prefix 'parent.' to color name. This
            may be useful if one wants some parent colors reversed.
            Example:
                ViewBack=ViewText
                ViewText=ViewBack
            this makes both colors the same because ViewBack is defined
            before ViewText. Only ViewBack value is loaded from parent theme.
            Better way is:
                Viewback=parent.ViewText
                ViewText=parent.ViewBack
        Converting text to real color fails and retuns black color if
            a) the string does not look like real color(e.g, typo as in
            'grean bold'), b) parent theme has not loaded yet, c) parent
            theme does not have the color
            with the same name
    Other keys are considered as objects - see Obj* constants, just drop
        'Obj' at the beginning of the key name
    One is not limited with only predefined color and object names.
    The theme can inroduce its own objects, e.g. to provide a runes or
        colors for new control that is not in standard library
To see the real world example of full featured theme, please see
    included theme 'turbovision'
*/
type ThemeManager struct {
	// available theme list
	themes map[string]theme
	// name of the current theme
	current   string
	themePath string
	version   string
}

const defaultTheme = "default"
const themeSuffix = ".theme"

var (
	themeManager *ThemeManager
	thememtx     sync.RWMutex
)

// ThemeDesc is a detailed information about theme:
// title, author, version number
type ThemeDesc struct {
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

// NewThemeManager creates a new theme manager
func initThemeManager() {
	themeManager = new(ThemeManager)
	ThemeReset()
}

// Reset removes all loaded themes from cache and reinitialize
// the default theme
func ThemeReset() {
	thememtx.Lock()
	defer thememtx.Unlock()

	themeManager.current = defaultTheme
	themeManager.themes = make(map[string]theme, 0)

	defTheme := theme{parent: "", title: "Default Theme", author: "Vladimir V. Markelov", version: "1.0"}
	defTheme.colors = make(map[string]term.Attribute, 0)
	defTheme.objects = make(map[string]string, 0)

	defTheme.objects[ObjSingleBorder] = "─│┌┐└┘"
	defTheme.objects[ObjDoubleBorder] = "═║╔╗╚╝"
	defTheme.objects[ObjEdit] = "←→V"
	defTheme.objects[ObjScrollBar] = "░■▲▼◄►"
	defTheme.objects[ObjViewButtons] = "^↓○[]"
	defTheme.objects[ObjCheckBox] = "[] X?"
	defTheme.objects[ObjRadio] = "() *"
	defTheme.objects[ObjProgressBar] = "░▒"
	defTheme.objects[ObjBarChart] = "█─│┌┐└┘┬┴├┤┼"
	defTheme.objects[ObjSparkChart] = "█"
	defTheme.objects[ObjTableView] = "─│┼▼▲"

	defTheme.colors[ColorDisabledText] = ColorBlackBold
	defTheme.colors[ColorDisabledBack] = ColorWhite
	defTheme.colors[ColorText] = ColorWhite
	defTheme.colors[ColorBack] = ColorBlack
	defTheme.colors[ColorViewBack] = ColorBlack
	defTheme.colors[ColorViewText] = ColorWhite

	defTheme.colors[ColorControlText] = ColorWhite
	defTheme.colors[ColorControlBack] = ColorBlack
	defTheme.colors[ColorControlActiveText] = ColorWhite
	defTheme.colors[ColorControlActiveBack] = ColorMagenta
	defTheme.colors[ColorControlShadow] = ColorBlue
	defTheme.colors[ColorControlDisabledText] = ColorWhite
	defTheme.colors[ColorControlDisabledBack] = ColorBlack

	defTheme.colors[ColorButtonText] = ColorWhite
	defTheme.colors[ColorButtonBack] = ColorGreen
	defTheme.colors[ColorButtonActiveText] = ColorWhite
	defTheme.colors[ColorButtonActiveBack] = ColorMagenta
	defTheme.colors[ColorButtonShadow] = ColorBlue
	defTheme.colors[ColorButtonDisabledText] = ColorWhite
	defTheme.colors[ColorButtonDisabledBack] = ColorBlack

	defTheme.colors[ColorEditText] = ColorBlack
	defTheme.colors[ColorEditBack] = ColorWhite
	defTheme.colors[ColorEditActiveText] = ColorBlack
	defTheme.colors[ColorEditActiveBack] = ColorYellow
	defTheme.colors[ColorSelectionText] = ColorYellow
	defTheme.colors[ColorSelectionBack] = ColorBlue

	defTheme.colors[ColorScrollBack] = ColorBlack
	defTheme.colors[ColorScrollText] = ColorWhite
	defTheme.colors[ColorThumbBack] = ColorBlack
	defTheme.colors[ColorThumbText] = ColorWhite

	defTheme.colors[ColorProgressText] = ColorBlue
	defTheme.colors[ColorProgressBack] = ColorBlack
	defTheme.colors[ColorProgressActiveText] = ColorBlack
	defTheme.colors[ColorProgressActiveBack] = ColorBlue
	defTheme.colors[ColorProgressTitleText] = ColorWhite

	defTheme.colors[ColorBarChartBack] = ColorBlack
	defTheme.colors[ColorBarChartText] = ColorWhite

	defTheme.colors[ColorSparkChartBack] = ColorBlack
	defTheme.colors[ColorSparkChartText] = ColorWhite
	defTheme.colors[ColorSparkChartBarBack] = ColorBlack
	defTheme.colors[ColorSparkChartBarText] = ColorCyan
	defTheme.colors[ColorSparkChartMaxBack] = ColorBlack
	defTheme.colors[ColorSparkChartMaxText] = ColorCyanBold

	defTheme.colors[ColorTableText] = ColorWhite
	defTheme.colors[ColorTableBack] = ColorBlack
	defTheme.colors[ColorTableSelectedText] = ColorWhite
	defTheme.colors[ColorTableSelectedBack] = ColorBlack
	defTheme.colors[ColorTableActiveCellText] = ColorWhiteBold
	defTheme.colors[ColorTableActiveCellBack] = ColorBlack
	defTheme.colors[ColorTableLineText] = ColorWhite
	defTheme.colors[ColorTableHeaderText] = ColorWhite
	defTheme.colors[ColorTableHeaderBack] = ColorBlack

	themeManager.themes[defaultTheme] = defTheme
}

// SysColor returns attribute by its id for the current theme.
// The method panics if theme loop is detected - check if
// parent attribute is correct
func SysColor(color string) term.Attribute {
	thememtx.RLock()
	sch, ok := themeManager.themes[themeManager.current]
	if !ok {
		sch = themeManager.themes[defaultTheme]
	}
	thememtx.RUnlock()

	clr, okclr := sch.colors[color]
	if !okclr {
		visited := make(map[string]int, 0)
		visited[themeManager.current] = 1
		if !ok {
			visited[defaultTheme] = 1
		}

		for {
			if sch.parent == "" {
				break
			}
			themeManager.loadTheme(sch.parent)
			thememtx.RLock()
			sch = themeManager.themes[sch.parent]
			clr, okclr = sch.colors[color]
			thememtx.RUnlock()

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

// SysObject returns object look by its id for the current
// theme. E.g, border lines for frame or arrows for scrollbar.
// The method panics if theme loop is detected - check if
// parent attribute is correct
func SysObject(object string) string {
	thememtx.RLock()
	sch, ok := themeManager.themes[themeManager.current]
	if !ok {
		sch = themeManager.themes[defaultTheme]
	}
	thememtx.RUnlock()

	obj, okobj := sch.objects[object]
	if !okobj {
		visited := make(map[string]int, 0)
		visited[themeManager.current] = 1
		if !ok {
			visited[defaultTheme] = 1
		}

		for {
			if sch.parent == "" {
				break
			}

			themeManager.loadTheme(sch.parent)
			thememtx.RLock()
			sch = themeManager.themes[sch.parent]
			obj, okobj = sch.objects[object]
			thememtx.RUnlock()

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

// ThemeNames returns the list of short theme names (file names)
func ThemeNames() []string {
	var str []string
	str = append(str, defaultTheme)

	path := themeManager.themePath
	if path == "" {
		path = "." + string(os.PathSeparator)
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic("Failed to read theme directory: " + themeManager.themePath)
	}

	for _, f := range files {
		name := f.Name()
		if !f.IsDir() && strings.HasSuffix(name, themeSuffix) {
			str = append(str, strings.TrimSuffix(name, themeSuffix))
		}
	}

	return str
}

// CurrentTheme returns name of the current theme
func CurrentTheme() string {
	thememtx.RLock()
	defer thememtx.RUnlock()

	return themeManager.current
}

// SetCurrentTheme changes the current theme.
// Returns false if changing failed - e.g, theme does not exist
func SetCurrentTheme(name string) bool {
	thememtx.RLock()
	_, ok := themeManager.themes[name]
	thememtx.RUnlock()

	if !ok {
		tnames := ThemeNames()
		for _, theme := range tnames {
			if theme == name {
				themeManager.loadTheme(theme)
				break
			}
		}
	}

	thememtx.Lock()
	defer thememtx.Unlock()

	if _, ok := themeManager.themes[name]; ok {
		themeManager.current = name
		return true
	}
	return false
}

// ThemePath returns the current directory with theme inside it
func ThemePath() string {
	return themeManager.themePath
}

// SetThemePath changes the directory that contains themes.
// If new path does not equal old one, theme list reloads
func SetThemePath(path string) {
	if path == themeManager.themePath {
		return
	}

	themeManager.themePath = path
	ThemeReset()
}

// loadTheme loads the theme if it is not in the cache already.
// If theme is in the cache loadTheme does nothing
func (s *ThemeManager) loadTheme(name string) {
	thememtx.Lock()
	defer thememtx.Unlock()

	if _, ok := s.themes[name]; ok {
		return
	}

	theme := theme{parent: defaultTheme, title: "", author: ""}
	theme.colors = make(map[string]term.Attribute, 0)
	theme.objects = make(map[string]string, 0)

	file, err := os.Open(s.themePath + string(os.PathSeparator) + name + themeSuffix)
	if err != nil {
		panic("Failed to open theme " + name + " : " + err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		// skip comments
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "/") {
			continue
		}

		// skip invalid lines
		if !strings.Contains(line, "=") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if (strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) ||
			(strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) {
			toTrim, _ := utf8.DecodeRuneInString(value)
			value = strings.Trim(value, string(toTrim))
		}

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
			// the first case is a reference to existing color (of this or parent theme)
			// the second is the real color
			if strings.HasSuffix(value, "Back") || strings.HasSuffix(value, "Text") {
				clr, ok := theme.colors[value]
				if !ok {
					v := value
					// if color starts with 'parent.' it means the parent color
					// must be used always. It may be useful to load inversed
					// text and background colors of parent theme
					if strings.HasPrefix(v, "parent.") {
						v = strings.TrimPrefix(v, "parent.")
					}
					sch, schOk := s.themes[theme.parent]
					if schOk {
						clr, ok = sch.colors[v]
					} else {
						panic(fmt.Sprintf("%v: Parent theme '%v' not found", name, theme.parent))
					}
				}
				if ok {
					theme.colors[key] = clr
				} else {
					panic(fmt.Sprintf("%v: Failed to find color '%v' by reference", name, value))
				}
			} else {
				c := StringToColor(value)
				if c%32 == 0 {
					panic("Failed to read color: " + value)
				}
				theme.colors[key] = c
			}
		} else {
			theme.objects[key] = value
		}
	}

	s.themes[name] = theme
}

// ReloadTheme refresh cache entry for the theme with new
// data loaded from file. Use it to apply theme changes on
// the fly without resetting manager or restarting application
func ReloadTheme(name string) {
	if name == defaultTheme {
		// default theme cannot be reloaded
		return
	}

	thememtx.Lock()
	if _, ok := themeManager.themes[name]; ok {
		delete(themeManager.themes, name)
	}
	thememtx.Unlock()

	themeManager.loadTheme(name)
}

// ThemeInfo returns detailed info about theme
func ThemeInfo(name string) ThemeDesc {
	themeManager.loadTheme(name)

	thememtx.RLock()
	defer thememtx.RUnlock()

	var theme ThemeDesc
	if t, ok := themeManager.themes[name]; !ok {
		theme.parent = t.parent
		theme.title = t.title
		theme.version = t.version
	}
	return theme
}

// RealColor returns attribute that should be applied to an
// object. By default all attributes equal ColorDefault and
// the real color should be retrieved from the current theme.
// Attribute selection work this way: if color is not ColorDefault,
// it is returned as is, otherwise the function tries to load
// color from the theme.
// clr - current object color
// id - color ID in theme
func RealColor(clr term.Attribute, id string) term.Attribute {
	if clr == ColorDefault {
		clr = SysColor(id)
	}

	if clr == ColorDefault {
		panic("Failed to load color value for " + id)
	}

	return clr
}
