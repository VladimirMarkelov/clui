# Hotkeys used in the library
The following hotkeys are built-in ones and cannot be overrided or disabled from an application

### Global hotkeys
- Ctrl+Q Ctrl+Q - exit application

### Window manipulations
- Ctrl+W Ctrl+H - moves the active Window to the bottom of window stack
- Ctrl+W Ctrl+M - maximizes/restores the active Window
- Ctrl+W Ctrl+C - closes the active Window. If the windows is the last visible window of an application then application closes as well
- Ctrl+P <arrow> - changes active Window position: moves Window to the direction of <arrow>
- Ctrl+S <arrow> - changes active Window size: Left and Down increase width and height, Right and Up decrease width and height

Note: Ctrl+P and Ctrl+S are sticky combinations. It means that if you want to move/resize active Window by more tham one character you do not need to press Ctrl+P or Ctrl+S every time. You just press Ctrl+S/P and then press the same arrow key as many times as you need. Sticky mode is off when you press any key other key.

### Control interaction in a Window
- TAB - selects the next control inside active Window
- Alt+PgDn - the same as TAB
- Alt+PgUp - selects the previous control inside active Window
- Space - click Button, Checkbox or RadioGroup control if the control is active
- Ctrl+C - copy text from active EditField (currently is not supported on OSX)
- Ctrl+V - paste text to active EditField - old text is replaced (currently is not supported on OSX)
- Ctrl+R - clears the active EditField
