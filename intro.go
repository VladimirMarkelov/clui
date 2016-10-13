/*
Package clui is an UI library to create simple interactive console applications.
Inspired by Borland TurboVision.

The library includes a set of simple but useful controls to create a multi-Windows
application with mouse support easily. Available at this moment controls:
    - EditField is one line edit text control
    - Label is static control to output single- or multi-line texts. The control
        supports multicolor output with tags
    - Frame is a border decoration and container
    - ListBox is a control to display a set of items with scrollbar
    - CheckBox is a control to allow to choose between 2 or 3 choices
    - RadioGroup is a control to select one item of the available
    - ProgressBar is a control to show some task progress
    - Button is push button control
    - Dialogs for a user to confirm something or to choose an item
    - more controls to come later
* Drag-n-drop with mouse is not supported due to limitations of some terminals.

Built-in theme support feature. Change the control look instantly without
restarting the application

Predefined hotkeys(hardcoded).
One touch combinations:
    - TAB to select the next control in the current View
    - Alt+PgDn, Alt+PgUp to select next or previous control in the current View,
        respectively. Hotkeys added because it is not possible to catch
        TAB contol with Shift or Ctrl modifier
    - Space to click a Button, CheckBox or RadioGroup if the control is active
    - Ctrl+R to clear the EditField
    - Arrows, Home, and End to move cursor inside EditField and ListBox
Sequences:
    At first one should press a sequence start combination:
        Ctrl+W to execute any View related command
        Ctrl+P to begin changing View position
        Ctrl+S to begin changing View size
        Ctrl+Q is used only to quit application - press Ctrl+Q twice
    And the next pressed key is considered as a subcomand:
        Ctrl+S and Ctrl+P processes only arrows. These commands supports
            multi-press - while one presses the same arrow after Ctrl+*, he/she
            does not need to press Ctrl+P or Ctrl+S before pressing each arrow
        Ctrl+W allow to:
            Ctrl+H moves the View to the bottom of View stack and activates a
                View below the View
            Ctrl+M maximizes or restores the current View
            Ctrl+C closes the current View. If it is the only View the application
                closes
*/
package clui
