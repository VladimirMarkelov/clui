package clui

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
	colors  map[ColorType]Color
	objects map[ObjectType]rune
	parent  string
}

const defaultTheme = "default"

// Theme color constants
const (
	// Font color for
	ColorActiveText = iota
	ColorInactiveText
	ColorActiveBack
	ColorInactiveBack
	// Font color for disabled controls
	ColorGrayText
	// Background color for inactive interactive control (e.g, for Button, CheckBox etc)
	ColorControlBack
	// Background color for active interactive control (e.g, for Button, CheckBox etc)
	ColorControlActiveBack
	// Font color for enabled controls (e.g, Button, Label etc)
	ColorControlText
	// Background color for control shadow (e.g, Button)
	ColorControlShadow
	// Background color for Window and controls with transparent background like Label or Frame
	ColorViewBack
	// Background color for text containers (e.g, EditField, ListBox etc)
	ColorEditBack
	// Font color for text containers (e.g, EditField, ListBox etc)
	ColorEditText
	// Font color for selected item (e.g, ListBox)
	ColorSelectionText
	// Background color for selected item (e.g, ListBox)
	ColorSelectionBack
	// Font color for ScrollBar
	ColorScroll
	// Font color for ScrollBar thumb
	ColorScrollThumb
	// Background color for ScrollBar
	ColorScrollBack
	// Font color for buttons in Window title
	ColorIcon
	// Font color for filled ProgressBar part
	ColorProgressOn
	// Font color for empty ProgressBar part
	ColorProgressOff
	// Background color for filled ProgressBar part
	ColorProgressOnBack
	// Background color for empty ProgressBar part
	ColorProgressOffBack
	// Background color for Menu
	ColorMenuBack
	// Font color for Menu
	ColorMenuText
)

// Symbols to represent some objects
const (
	// Single Border elements
	//  HorizontalLine, VerticalLine, UpperLeft corner, UpperRight Corner, BottomLeft Corner, BottomRight Corner
	ObjSingleBorderHLine    = iota // H V UL UR DL DR
	ObjSingleBorderVLine           // H V UL UR DL DR
	ObjSingleBorderULCorner        // H V UL UR DL DR
	ObjSingleBorderURCorner        // H V UL UR DL DR
	ObjSingleBorderDLCorner        // H V UL UR DL DR
	ObjSingleBorderDRCorner        // H V UL UR DL DR

	// Double Border elements H V UL UR DL DR
	ObjDoubleBorderHLine
	ObjDoubleBorderVLine
	ObjDoubleBorderULCorner
	ObjDoubleBorderURCorner
	ObjDoubleBorderDLCorner
	ObjDoubleBorderDRCorner

	// ObjScrollBar           // bar thumb upArrow downArrow
	ObjScrollBar       // bar thumb upArrow downArrow
	ObjScrollThumb     // bar thumb upArrow downArrow
	ObjScrollUpArrow   // bar thumb upArrow downArrow
	ObjScrollDownArrow // bar thumb upArrow downArrow

	// ObjIcons  // hide close start end
	ObjIconMinimize
	ObjIconDestroy
	ObjIconOpen
	ObjIconClose

	// ObjEdit                // < > {arrows for edit field} ► {wordwrap sign}
	ObjEditLeftArrow  // <
	ObjEditRightArrow // >
	ObjEditWordWrap   // ►

	// ObjCheckbox            // [ ] o x ?
	ObjCheckboxOpen      // [
	ObjCheckboxClose     // ]
	ObjCheckboxChecked   // x
	ObjCheckboxUnchecked // o
	ObjCheckboxUnknown   // ?

	// ObjRadiobutton         // ( ) o x
	ObjRadioOpen       // (
	ObjRadioClose      // )
	ObjRadioSelected   // x
	ObjRadioUnselected // o

	// ObjCombobox            // v
	ObjComboboxDropDown // v

	// ObjProgressBar         // ░▒ (off filled)
	ObjProgressBarFull  // ▒ filled
	ObjProgressBarEmpty // ░ empty
)

type (
	ColorType  int
	ObjectType int
)

func NewThemeManager() *ThemeManager {
	sm := new(ThemeManager)
	sm.current = defaultTheme
	sm.themes = make(map[string]theme, 0)

	defTheme := theme{parent: ""}
	defTheme.colors = make(map[ColorType]Color, 0)
	defTheme.objects = make(map[ObjectType]rune, 0)

	defTheme.colors[ColorActiveText] = ColorBrightWhite
	defTheme.colors[ColorInactiveText] = ColorWhite
	defTheme.colors[ColorActiveBack] = ColorBlack
	defTheme.colors[ColorInactiveBack] = ColorBlack
	defTheme.colors[ColorGrayText] = ColorWhite
	defTheme.colors[ColorControlActiveBack] = ColorMagenta
	defTheme.colors[ColorControlShadow] = ColorBlue
	defTheme.colors[ColorControlBack] = ColorBlack
	defTheme.colors[ColorControlText] = ColorWhite
	defTheme.colors[ColorViewBack] = ColorBrightBlack
	defTheme.colors[ColorEditBack] = ColorWhite
	defTheme.colors[ColorEditText] = ColorBlack
	defTheme.colors[ColorSelectionText] = ColorWhite
	defTheme.colors[ColorSelectionBack] = ColorBlack
	defTheme.colors[ColorScroll] = ColorBlack
	defTheme.colors[ColorScrollThumb] = ColorBlack
	defTheme.colors[ColorScrollBack] = ColorWhite
	defTheme.colors[ColorIcon] = ColorWhite
	defTheme.colors[ColorProgressOn] = ColorBrightBlue
	defTheme.colors[ColorProgressOff] = ColorBlue
	defTheme.colors[ColorProgressOnBack] = ColorBrightBlack
	defTheme.colors[ColorProgressOffBack] = ColorBrightBlack
	defTheme.colors[ColorMenuBack] = ColorBlack
	defTheme.colors[ColorMenuText] = ColorWhite

	// defTheme.objects[ObjSingleBorder] = "─│┌┐└┘"
	defTheme.objects[ObjSingleBorderHLine] = '─'
	defTheme.objects[ObjSingleBorderVLine] = '│'
	defTheme.objects[ObjSingleBorderULCorner] = '┌'
	defTheme.objects[ObjSingleBorderURCorner] = '┐'
	defTheme.objects[ObjSingleBorderDLCorner] = '└'
	defTheme.objects[ObjSingleBorderDRCorner] = '┘'

	// defTheme.objects[ObjDoubleBorder] = "═║╔╗╚╝"
	defTheme.objects[ObjDoubleBorderHLine] = '═'
	defTheme.objects[ObjDoubleBorderVLine] = '║'
	defTheme.objects[ObjDoubleBorderULCorner] = '╔'
	defTheme.objects[ObjDoubleBorderURCorner] = '╗'
	defTheme.objects[ObjDoubleBorderDLCorner] = '╚'
	defTheme.objects[ObjDoubleBorderDRCorner] = '╝'

	// defTheme.objects[ObjScrollBar] = "|O^v"
	defTheme.objects[ObjScrollBar] = '|'
	defTheme.objects[ObjScrollThumb] = 'O'
	defTheme.objects[ObjScrollUpArrow] = '^'
	defTheme.objects[ObjScrollDownArrow] = 'v'

	// defTheme.objects[ObjIcons] = "↓○[]"
	defTheme.objects[ObjIconMinimize] = '↓'
	defTheme.objects[ObjIconDestroy] = '○'
	defTheme.objects[ObjIconOpen] = '['
	defTheme.objects[ObjIconClose] = ']'

	// defTheme.objects[ObjEdit] = "←→►"
	defTheme.objects[ObjEditLeftArrow] = '←'
	defTheme.objects[ObjEditRightArrow] = '→'
	defTheme.objects[ObjEditWordWrap] = '►'

	// defTheme.objects[ObjCheckbox] = "[] X?"
	defTheme.objects[ObjCheckboxOpen] = '['
	defTheme.objects[ObjCheckboxClose] = ']'
	defTheme.objects[ObjCheckboxChecked] = 'X'
	defTheme.objects[ObjCheckboxUnchecked] = ' '
	defTheme.objects[ObjCheckboxUnknown] = '?'

	// defTheme.objects[ObjRadiobutton] = "() *"
	defTheme.objects[ObjRadioOpen] = '('
	defTheme.objects[ObjRadioClose] = ')'
	defTheme.objects[ObjRadioSelected] = '*'
	defTheme.objects[ObjRadioUnselected] = ' '

	defTheme.objects[ObjComboboxDropDown] = 'V'

	// defTheme.objects[ObjProgressBar] = "░▒"
	defTheme.objects[ObjProgressBarFull] = '▒'
	defTheme.objects[ObjProgressBarEmpty] = '░'

	sm.themes[defaultTheme] = defTheme

	return sm
}

func (s *ThemeManager) GetSysColor(color ColorType) Color {
	sch, ok := s.themes[s.current]
	if !ok {
		sch = s.themes[defaultTheme]
	}

	clr, okclr := sch.colors[color]
	if !okclr && sch.parent != "" {
		sch = s.themes[sch.parent]
		clr, okclr = sch.colors[color]
		if !okclr {
			clr = ColorDefault
		}
	}

	return clr
}

func (s *ThemeManager) GetSysObject(object ObjectType) rune {
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
