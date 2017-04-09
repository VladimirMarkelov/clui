# About Standard Widgets

### Creating a widget
Every widget should have parent at creation time otherwise the widget will be invisible and will not recieve any message. Parent and minimal width are the only common arguments of all functions that creates a new widget. Other arguments may vary but there are set of common ones. Generic create function may look like this (except `Frame` that has its own argument 'frameWidth': `BorderThick` or `BorderThin`) - note that it is not real function:
```
CreateWidget(parent, minimalWidth, minimalHeight, title, scale)
```
parent - is a container which created widget belongs;

minimalWidth and minimalHeight - are limits for the widget sizes when its parent is resized. Some widgets does not used minimalHeight because they cannot be higher than 1 symbol - those widgets are Radio, CheckBox, and EditField. If the widget should not have certain minimal size you can set it to `AutoSize`. In this case minimal size will be either calculated on the fly (e.g, Label minimal width is length of its title) or defaults are used (Frame gets width equals 5 and height equals 3 if they are not provided by caller). Minimal sizes can be changed later by calling a widget function `widget.SetConstraints(newMinimalWidth, newMinimalHeight)`. Please note that if you want to change only one value, you can set the other one to constant `KeepValue`, e.g widget.SetConstraints(KeepValue, 10)` changes only minimalHeight to 10 and keeps old minimalWidth value;

title - a text displayed on the widget. Half of widgets does not have it. Titles can include special tags to be display in various colors. A more details about it a bit later;

scale - defines how fast the widget changes its size when its parent is resizing (e.g, if a container has two widgets and the first has scale 1 while the second has scale 2 then the latter widget grows twice faster when the container is growing). Setting scale to a special constant `Fixed` forbids resizing the control, so it always has the constant size.

### Basic widget methods
All methods follow one rule: for every property there is a Set**Property** method to change it and there is a **Property** method to read it. Below only setters are mentioned:
1. Size and position: `SetSize(width, height)`, `SetPos(x, y)`
1. Minimal size: `SetConstraints(minWidth, minHeight)`
1. Title/text/caption: `SetTitle(string)`
1. Disable or enable widget - `SetEnable(bool)`. Disabled controls usually has its own look and does not respond to mouse and keyboard events
1. Activate control - `SetActive(bool)`. A Window can has only one active widget at a time, so `SetActive` deactivates previously activated control before activating a new one
1. Tab control: `SetTabStop(bool)`. Sets if the control can be selected by pressing TAB key or the control is skipped while traversing widgets with keyboard. In any case the widget can be selected with mouse
1. Layout type: `SetPack(PackType)`. Sets packing direction of widget children - Horizontal or Vertical
1. Space between the first(or last) child and widget edge: `SetPaddings(identX, identY)`
1. Space between children: `SetGaps(gapX, gapY)`
1. Set scaling mode: `SetScale(scale)`. The function defines the widget behavior when its parent is resizing. You can forbid widget resizing by setting scale to `Fixed` or you can set the scaling coefficient. Please read more about scaling in the next section
1. Text align: `SetAlign(Align)`. Sets alignment for the widget caption. The property is applied to widgets that do not allow to edit its content: text of `Label`, caption of `Frame` etc. Multiline and editable widgets(`EditField`, `TextView` etc) does not use the property. As for `TableView` - each of its column has its own alignment
1. Basic widget colors: `SetTextColor(Color)` and `SetBackColor(Color)`. It is default colors to display a widget content (see below in section 'Title color tags' a bit more about default color usage). There is a special color `Default` that is set bt default for each widget: if you set widget color to `Default` the widget will use colors from current theme(it works only for standard widget included into the library). So, for simple application you do not need to set colors - the library does it for you

#### Methods that do not follow rule of method naming (since they are not often used)
1. Active widget colors. Set it with `SetActiveTextColor(Color)` and `SetActiveBackColor(Color)` but read with `ActiveColors() (textColor, backColor)`. There are some widgets that a user can interact while those widgets are active. E.g, only active `EditField` can be modified, only active `Button` can be pressed with key Space. Setting your own active color for a widget may help to indentify an active widget easier. And as in basic color case you can use `Default` color if you are fine with the colors provided by current theme - default theme has active colors different from basic ones

### How scaling works
Every container has its starting size that calculated as maximum of two values: its minimal size and sum of minimal sizes of its children. When the container changes its size then a layout manager does children resizing:
1. The difference between new size and starting size(`Delta`) is calculated
1. It calculates the total scale size of all children(`TotalScaleSize`). `Fixed` children have scale coefficient equals 0
1. Each child that has scale greater than 0(except the last one) gets increased by the value: `child.Scale * Delta / TotalScaleSize`. The number is rounded down `Delta` decreased by the value. The latter child gets increased by the value that remains in `Delta`. So, in real life one widget can grow more than the other one even their scales equal (example: two children have scale 1, parent size increased by 3, then the former gets increased by `floor(1 * (3 / (1 + 1))) = 1`, and the latter gets increased by `(3 - 1) = 2`

### Title color tags
Every widget can display colored texts: Labels, Frames, even ListBox items can be colorized. The way of using colors in output strings is similar to te way HTML uses. Every color tag must start with `<Letter:ColorValue>`. `Letter` must be one of 'c', 't', 'f', or 'b' (if `Letter` is not one of the list or colon does not follow the `Letter` then the text displayed as is. So, you do not need to escape angle brackets to display them in text. Use tag `b` to change **b**ackground color, and `c`, `t` or `f` to change **t**ext or **f**oreground color (`c` means **c**olor). `ColroValue` is a one color value or a list of color and its modifier separated with space, `+` or `|`. Supported colors are black, white, green, yellow, blue, magenta, cyan, red, and `default`. The last value `default` means that the corresponding widget color should be used - background for `b:default`, and foreground for `f:default`. Supported modifiers are `bold`(or `bright` instead of it), `underline`(spelling `underlined` supported as well), and `reverse`. Examples of correct tags: `<t:red bold>`, `<t:underline+blue>`.

All tags inside the string works as triggers: when drawing function runs into the tag it changes the current color to tag's one until the end of line of another tag is met. After printing all the string the colors are back to default ones. In other words, colors inside a string are temporary and you do not need to reset colors to normal ones at the end.

Sometimes you need to just hilight only a part of with different color and keep all the rest with default colors. Just enclose the part of text between `<c:Color>` and `<c:default>`. For the default case there is a shortcut - empty tag `<b:>` or `<c:>`. Example: `Label` has text `The <c:green>green<c:> text` displays "green" with green text color while the rest of the string is displayed with `Label` text color that can be green as well.

The library has a set of public functions to use in custom widgets:

`UnColorizeText` - retuns a text with all color tags removed

`AlignColorizedText` - align the text with tags inside an area with certain width

`SliceColorized` - smart text slicing that keeps the tag from the beginning if the tag is not included into slice range. Example: `SliceColorized("abc<c:green>def<c:red>hg", 4, -1) == "<c:green>f<c:red>hg"`

`StringToColor` - get the string in tag format and returns a color value that can be used in functions like `SetTextColor`. Can be useful to read colors from configuration file. `ColorToString` - does the opposite it returns a color description that can be used in tags

