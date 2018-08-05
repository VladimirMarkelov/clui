# File Picker

Use file picker if you need to select a single file or directory.

<img src="/docs/img/fselect.png" alt="File picker dialog">

The dialog includes:

- title that shows the custom text followed by file masks
- current directory
- list of directories and files inside the current directory
- edit field to a) enter a name of a new file or directory b) quick search. When edit field has focus you can use arrows up and down to select previous or next object
- button **Open** enters the selected directory. The button is useful in select directory mode
- button **Select** closes the dialog and returns a path to the selected object. If edit field is empty the path is the path to the object selected in the file list box, otherwise the path is made from current directory and edit field text
- button **Cancel** closes the dialog and does not return path to selected object

### Returned values

After the dialog is closed a few its properties contains information what object a user has selected:

- **Selected** contains information about how the dialog was closed: true - a user has selected an object and clicked **Select**, false - a user canceled the dialog without selecting any object
- **Exists** is true if a user has selected existing file or directory, and it is false if a user entered name of a new object and clicked **Select**. The latter is possible only if option **mustExist** is set to false
- **FilePath** is a full path to the selected object

### API

To show a dialog, call the function
```
func CreateFileSelectDialog(title, fileMasks, initPath string, selectDir, mustExist bool) *FileSelectDialog
```

Function arguments:

- **title** is a custom dialog title. It should not contain file masks because the dialog always shows title and the file masks follows it
- **fileMasks** is list of file masks separated with comma or OS path separator (';' - for Windows, ':' - for Linux). Empty, "*", and "*.*" mean *all files*
- **initPath** sets the starting directory for the dialog. If it is empty then the dialog uses the current working directory. If the **intiPath** does not exist, then the dialog looks up for the first existing directory in the directory tree starting from **initPath**. In case it fails to find any existing directory, the dialog opens the current working directory
- **selectDir** - set it to *true* if you want to select a directory instead of a file. In case of **selectDir** is *true*, the file list box does not display regular files
- **mustExist** set it to *true* if you want a user to select only existing object. If it is *false* then the dialog allows a user to enter any name into edit field and click **Select**. The latter is useful for "File save" dialog.

Do not forget to set a callback that is called after the dialog is closed:
```
bookFile := ""
dlg := CreateFileSelectDialog("Select a book to read", "*.fb2,*.epub,*.txt", "", false, true)
dlg.OnClose(func() {
    if !dlg.Selected {
        // a user canceled the dialog
        return
    }
    bookFile = dlg.FilePath
})
```

Please, check the ![dialog demo](/demos/fileselect/fselect.go) for more details.

