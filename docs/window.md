# About Window
Window is the main control of the library - you cannot run a CLUI-based application without any Window. At the same time an application can have unlimited number of Windows at a time. Window has a lot in common with other controls yet it have a set of its own unique features.

### Creating new Window
To create a new Window, use function
```
AddWindow(positionX, positionY, minWidth, minHeight, title)
```
instead of `CreateWindow`. All other controls must be create with Create**ControlType** functions. In fact, `CreateWindow` function can be used but it only creates a new Window but does not do anything else. So, a Window created with `CreateWindow` is not dispalyed on the screen and a user cannot manipulate it. `AddWindow` does all the internal job: registers Windows, makes it visible and updates the screen to apply changes. Both minWidth and minHeight can be set to `AutoSize` - in this case minimal width is 10 and heigth is 5.

### Window features
* It is the only control that can be moved or resized manually(using mouse or keyboard). There is a set of predefined hotkeys to do common tasks witout mouse. All hotkeys are key sequences - you have to press the first combination, release it and then press the second key (all default hotkeys are easy to remember, I added info why the keys was used for the certain action):
  * CtrlS "arrow key" - change active Window size in direction of pressed key. CtrlS - Ctrl+**S**ize
  * CtrlP "arrow key" - move active Window in the direction of pressed key. CtrlP - Ctrl+**P**osition
  * CtrlW CtrlM - maximize or restore active Window. CtrlW - Ctrl+**W**indow, CtrlM - Ctrl+**M**aximize
  * CtrlW CtrlH - move active Window to bottom of window stack and makes active the next window in the list. It does nothing if there is only one Window on the screen. CtrlW - Ctrl+**W**indow, CtrlH - Ctrl+**H**ide
* Windows have borders that indicates its activity: currently active Window has double border, while all others have single border
* Every Window has 'icons' at the right to corner to manipulate Window with mouse. The available 'icons' (it is a default set - from left to right): move to background, maximize/restore, and close. Please note that closing the last Window terminates application
* Window grabs **TAB** key control to support moving to the next child using keyboard
