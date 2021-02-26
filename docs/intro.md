# Introduction
CLUI is a UI library for text terminals. It provides a set of standard UI controls to build rather complex UI applications. The library supports theming - you can change the look of a running application easily on the fly (please see an example of theme changing on the fly in the screenshot section of README).

## Limitations
The library supports a wide variety of the terminals that termbox supports and should work smoothly on Windows, Linux, and Mac. But it is not fully compatible with "UNIX on Windows" like MSYS or Git bash. For instance, the library does not work with newer Git bash (included in Git for Windows 2.x) but I have no trouble running the CLUI demo inside older Git bash (included in Git for Windows 1.x - yes it is pretty old, but it is compatible).

## Features
* A rich set of controls out of the box
* Theming - you do not have to set a color for everything: if you have a theme enabled you can use ColorDefault for text and background everywhere
* Automatic widget positioning and sizing depending on parent or terminal size: you do not have to set a widget position manually. The only thing you have to do is to set the mininal width and height of a widget. But even in this case, you can just pass constant AutoSize and the minimal sizes will be calculated by the library (e.g, minimal width of a Label and Button depends on their text length)
* Autoscaling widget on terminal/Window resize. Set scale coefficient while creating a widget and the library automatically will handle the widget resizing. Automatic widget resizing can be forbidden by setting scale coefficient to Fixed - in this case the widget keeps its size minimal, which can be useful for Buttons
* All texts can be colorized: there is a set of simple tags to be used inside text to temporarily change the color of text or background
* Many controls emit events when something happens (e.g, Button emits only OnClick event, ListBox emits events OnKeyPress that can be useful to support filtering or incremental search and OnSelectItem that is emitted when the selection changes)
* While Windows can overlap each other, a widget cannot overlap another widget. This may change in the future, but at this moment it would make it impossible to implement components like ComboBox or PopupMenu
* Windows can be manipulated with both mouse and keyboard - the library provides built-in hotkeys to do common tasks like Window moving or resizing (and an extra key sequence 'CtrlQ CtrlQ' to exit the application)
* A few common dialogs are available out of box: confirmation and selection dialog. Both have customizable button labels

## Contols provided by the library
All controls can be divided into 4 categories:
* Invisible controls - only RadioGroup falls in this category at this moment. These controls are helpers or do logical control grouping. E.g, RadioGroup makes sure that no more than one Radio button of the group is selected
* Top level controls - Window. It is the only visual control that does not have a parent and cannot be a child of any other control. Every Window is a separate dialog/window that a user can move inside the terminal window. Windows can be modal - it may be useful for creating confirmation dialogs. Every application must have at least one visible Window. If a user closes the last Window the application automatically terminates
* Widgets - almost all controls are in this group. Any widget can be both parent and child (though, it has not been thoroughly tested for cases when a parent is not a Frame, e.g, a button contains a few children - it is untested). A widget must be a child of a Window or another widget - it is not possible to display a widget without a parent
* Dialogs - ready to use modal Windows for user interaction. Two kinds of dialogs are available at this moment: ConfirmationDialog - which is a simple dialog with a question and up to 3 buttons with custom prompt texts, and SelectDialog - which is a Window presenting a list of items (dispayed as a ListBox or as a RadioGroup) for a user to select from

### The current list of available controls
* Window - is a top level control with unique features: manual resizing, maximize and minimize abilities, overlapping other Windows
* Label - is for displaying static texts. Text can be displayed in horizontal or vertical direction
* Button  - is a simple push button control
* EditField - is a control to edit text. It is limited to one line of text. There is basic clipboard support: copy and paste
* ListBox - is a scrollable control to display a list of items
* TextView is a scrollable (vertical and horizontal scrolls are supported) viewer for a lot of text. Optional feature: wordwrapping(that disables horizontal scroll) and limiting the maximum number of items that the control contains (that is useful, e.g, to simulate 'tail -f' and show only N last items - if a new line is added to full TextView then the first line is deleted automatically)
* ProgressBar - is a progress indicatior. Both vertical and horizontal direction are supported. For horizontal ProgressBar it is possible to display custom text over a control. Custom text supports a few internal variable like percentage or current value
* Frame - is a decorative control to draw a frame around a group of controls. Without border it can play the role of a spacer between controls for more precise resizing (e.g, to create a Window with a control of fixed size that sticks to the right side, you can create a borderless frame with scale coefficient 1 and then add the control with Fixed size)
* CheckBox - is a tri-state check box control (tri-state is disabled by default)
* Radio - is a simple radio button. It is useless when used as a separate control - it should be attached to RadioGroup
* RadioGroup - is a non-visual control to manage a group of RadioButtons. It makes sure that at any moment of time there is no more than one of the RadioButton is selected
* BarChart - is a chart representing grouped data. It supports displaying bars, real values under bars, and a legend. Bars can be displayed with gaps or one right after each other. BarChart supports custom coloring when drawing, so it is possible, e.g, mark a part of a bar red if its value more than a certain limit (please see barchart.go in demos for real-life examples). Vertical height of bars is always auto-sized, so the highest bar is always occipies the full control height. But bar width can be auto-sized(calculated depending on the number of bars and control width) or defined by the user
* SparkChart - is control similar to BarChart but for dynamic data. The maximum number of displayed bars depends on the control's width and the width of the area on the left of the displayed values. If you add a new value and the control is full then the oldest value is removed and al other bars shift to the left. It is possible to make bars autoscaled (in this case the highest bar occupies the whole control height and other bar heights are recalculted) or you can set a constant maximum (e.g, it may be useful to display CPU load history). Extra feature: display with different color the bars that have the highest value
* GridView - is a table to show structured data. It does not support in place editing, but the GridView emits events for sorting, deleting, adding, and modifying data. Extra features: optional lines between columns and rows, custom draw of any cell
* ConfirmationDialog - is a modal dialog to ask a user confirmation about some action. The dialog can display up to three buttons with custom text. Any button can be set to be the default one
* SelectDialog - is a modal dialog to select an item from a list. Lists can be displayed as a ListBox or as a RadioGroup
