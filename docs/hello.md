# First application

### Add the library to project
Create an empty application. At first you need to import CLUI library:
```
import (
    ui "github.com/VladimirMarkelov/clui"
)
```
I created an alias 'ui' for the imported library to use shorter name in calls.

### Initialization and finalization
The library must be initialized before creating the first control. Initialization creates control and theme mamagers, intializes termbox library and prepares a main event loop. Finalization just cleans up the terminal - call it before exiting your application. If you forget to call finalization or the application crashes it usually results in cursor disappearing because the library turns off the text cursor at its start.

In a simple application it can be done this way:
```
func main() {
    ui.InitLibrary()
    defer ui.DeinitLibrary()
    ... your other code ...
}
```

### Creating a Window
An UI application without a window is useless. Let's create an empty window. Add the following code after 'defer':
```
view := ui.AddWindow(0, 0, 10, 7, "Hello World!")
```
0, 0 - is the position of the new window. The top left corner in our case
10, 7 - minimal width and height of the window
"Hello World!" - is the window title

### Make the application work
The final step is to start the main event loop that is responsible for displaying and interacting all the UI stuff. Add this line before the final brace:
```
ui.MainLoop()
```
Note: this call must be the last line in the function because no code after this line is executed until the application is closed

### Add more controls
Empty window is boring. Let's create a button that closes the application when anyone clicks it. Add the code between library initialization and calling the main event loop:
```
    btnQuit := ui.CreateButton(view, 15, 4, "Hi", 1)
    btnQuit.OnClick(func(ev ui.Event) {
        go ui.Stop()
    })
```
The first line adds a button to our Window(the first argument is our Window). The button has minimal width 15 and height 4. Button text is 'Hi'. And the scaling coefficient is 1 that means the button will be automatically resized when its parent is resized. Try resizing the window and the button will always fill all the windows because the bitton is the only child of the window.

The second line adds an event callback that is fired when someone clicks the button with mouse or by pressing 'space' key. In the callback we just sends an event to the main loop that application is terminating. ui.Stop() - is a gentle way to exit terminal application.

The full code of the example can be found at ![demos/helloworld.go](/demos/helloworld.go)
