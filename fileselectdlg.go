package clui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	term "github.com/nsf/termbox-go"
)

// FileSelectDialog is a dialog to select a file or directory.
// Public properties:
//   * Selected - whether a user has selected an object, or the user canceled
//          or closed the dialog. In latter case all other properties are not
//          used and contain default values
//   * FilePath - full path to the selected file or directory
//   * Exists - if the selected object exists or a user entered manually a
//         name of the object
type FileSelectDialog struct {
	View     *Window
	FilePath string
	Exists   bool
	Selected bool

	fileMasks []string
	currPath  string
	mustExist bool
	selectDir bool

	result  int
	onClose func()

	listBox *ListBox
	edFile  *EditField
	curDir  *Label
}

// Set the cursor the first item in the file list
func (d *FileSelectDialog) selectFirst() {
	d.listBox.SelectItem(0)
}

// Checks if the file name matches the mask. Empty mask, *, and *.* match any file
func (d *FileSelectDialog) fileFitsMask(finfo os.FileInfo) bool {
	if finfo.IsDir() {
		return true
	}

	if d.selectDir {
		return false
	}

	if len(d.fileMasks) == 0 {
		return true
	}

	for _, msk := range d.fileMasks {
		if msk == "*" || msk == "*.*" {
			return true
		}

		matched, err := filepath.Match(msk, finfo.Name())
		if err == nil && matched {
			return true
		}
	}

	return false
}

// Fills the ListBox with the file names from current directory.
// Files which names do not match mask are filtered out.
// If select directory is set, then the ListBox contains only directories.
// Directory names ends with path separator
func (d *FileSelectDialog) populateFiles() error {
	d.listBox.Clear()
	isRoot := filepath.Dir(d.currPath) == d.currPath

	if !isRoot {
		d.listBox.AddItem("..")
	}

	f, err := os.Open(d.currPath)
	if err != nil {
		return err
	}

	finfos, err := f.Readdir(0)
	f.Close()
	if err != nil {
		return err
	}

	fnLess := func(i, j int) bool {
		if finfos[i].IsDir() && !finfos[j].IsDir() {
			return true
		} else if !finfos[i].IsDir() && finfos[j].IsDir() {
			return false
		}

		return strings.ToLower(finfos[i].Name()) < strings.ToLower(finfos[j].Name())
	}

	sort.Slice(finfos, fnLess)

	for _, finfo := range finfos {
		if !d.fileFitsMask(finfo) {
			continue
		}

		if finfo.IsDir() {
			d.listBox.AddItem(finfo.Name() + string(os.PathSeparator))
		} else {
			d.listBox.AddItem(finfo.Name())
		}
	}

	return nil
}

// Tries to find the best fit for the given path.
// It goes up until it gets into the existing directory.
// If all fails it returns working directory.
func (d *FileSelectDialog) detectPath() {
	p := d.currPath
	if p == "" {
		d.currPath, _ = os.Getwd()
		return
	}

	p = filepath.Clean(p)
	for {
		_, err := os.Stat(p)
		if err == nil {
			break
		}

		dirUp := filepath.Dir(p)
		if dirUp == p {
			p, _ = os.Getwd()
			break
		}

		p = dirUp
	}
	d.currPath = p
}

// Goes up in the directory tree if it is possible
func (d *FileSelectDialog) pathUp() {
	dirUp := filepath.Dir(d.currPath)
	if dirUp == d.currPath {
		return
	}
	d.currPath = dirUp
	d.populateFiles()
	d.selectFirst()
}

// Enters the directory
func (d *FileSelectDialog) pathDown(dir string) {
	d.currPath = filepath.Join(d.currPath, dir)
	d.populateFiles()
	d.selectFirst()
}

// Sets the EditField value with the selected item in ListBox if:
//   * a directory is selected and option 'select directory' is set
//   * a file is selected and option 'select directory' is not set
func (d *FileSelectDialog) updateEditBox() {
	s := d.listBox.SelectedItemText()
	if s == "" || s == ".." ||
		(strings.HasSuffix(s, string(os.PathSeparator)) && !d.selectDir) {
		d.edFile.Clear()
		return
	}
	d.edFile.SetTitle(strings.TrimSuffix(s, string(os.PathSeparator)))
}

// CreateFileSelectDialog creates a new file select dialog. It is useful if
// the application needs a way to select a file or directory for reading and
// writing data.
//  * title - custom dialog title (the final dialog title includes this one
//      and the file mask
//  * fileMasks - a list of file masks separated with comma or path separator.
//      If selectDir is true this option is not used. Setting fileMasks to
//      '*', '*.*', or empty string means all files
//  * selectDir - what kind of object to select: file (false) or directory
//      (true). If selectDir is true then the dialog does not show files
//  * mustExists - if it is true then the dialog allows a user to select
//       only existing object ('open file' case). If it is false then the dialog
//       makes possible to enter a name manually and click 'Select' (useful
//       for file 'file save' case)
func CreateFileSelectDialog(title, fileMasks, initPath string, selectDir, mustExist bool) *FileSelectDialog {
	dlg := new(FileSelectDialog)
	cw, ch := term.Size()
	dlg.selectDir = selectDir
	dlg.mustExist = mustExist

	if fileMasks != "" {
		maskList := strings.FieldsFunc(fileMasks,
			func(c rune) bool { return c == ',' || c == ':' })
		dlg.fileMasks = make([]string, 0, len(maskList))
		for _, m := range maskList {
			if m != "" {
				dlg.fileMasks = append(dlg.fileMasks, m)
			}
		}
	}

	dlg.View = AddWindow(10, 4, 20, 16, fmt.Sprintf("%s (%s)", title, fileMasks))
	WindowManager().BeginUpdate()
	defer WindowManager().EndUpdate()

	dlg.View.SetModal(true)
	dlg.View.SetPack(Vertical)

	dlg.currPath = initPath
	dlg.detectPath()
	dlg.curDir = CreateLabel(dlg.View, AutoSize, AutoSize, "", Fixed)
	dlg.curDir.SetTextDisplay(AlignRight)

	flist := CreateFrame(dlg.View, 1, 1, BorderNone, 1)
	flist.SetPaddings(1, 1)
	flist.SetPack(Horizontal)
	dlg.listBox = CreateListBox(flist, 16, ch-20, 1)

	fselected := CreateFrame(dlg.View, 1, 1, BorderNone, Fixed)
	// text + edit field to enter name manually
	fselected.SetPack(Vertical)
	fselected.SetPaddings(1, 0)
	CreateLabel(fselected, AutoSize, AutoSize, "Selected object:", 1)
	dlg.edFile = CreateEditField(fselected, cw-22, "", 1)

	// buttons at the right
	blist := CreateFrame(flist, 1, 1, BorderNone, Fixed)
	blist.SetPack(Vertical)
	blist.SetPaddings(1, 1)
	btnOpen := CreateButton(blist, AutoSize, AutoSize, "Open", Fixed)
	btnSelect := CreateButton(blist, AutoSize, AutoSize, "Select", Fixed)
	btnCancel := CreateButton(blist, AutoSize, AutoSize, "Cancel", Fixed)

	btnCancel.OnClick(func(ev Event) {
		WindowManager().DestroyWindow(dlg.View)
		WindowManager().BeginUpdate()
		dlg.Selected = false
		closeFunc := dlg.onClose
		WindowManager().EndUpdate()
		if closeFunc != nil {
			closeFunc()
		}
	})

	btnSelect.OnClick(func(ev Event) {
		WindowManager().DestroyWindow(dlg.View)
		WindowManager().BeginUpdate()
		dlg.Selected = true

		dlg.FilePath = filepath.Join(dlg.currPath, dlg.listBox.SelectedItemText())
		if dlg.edFile.Title() != "" {
			dlg.FilePath = filepath.Join(dlg.currPath, dlg.edFile.Title())
		}
		_, err := os.Stat(dlg.FilePath)
		dlg.Exists = !os.IsNotExist(err)

		closeFunc := dlg.onClose
		WindowManager().EndUpdate()
		if closeFunc != nil {
			closeFunc()
		}
	})

	dlg.View.OnClose(func(ev Event) bool {
		if dlg.result == DialogAlive {
			dlg.result = DialogClosed
			if ev.X != 1 {
				WindowManager().DestroyWindow(dlg.View)
			}
			if dlg.onClose != nil {
				dlg.onClose()
			}
		}
		return true
	})

	dlg.listBox.OnSelectItem(func(ev Event) {
		item := ev.Msg
		if item == ".." {
			btnSelect.SetEnabled(false)
			btnOpen.SetEnabled(true)
			return
		}
		isDir := strings.HasSuffix(item, string(os.PathSeparator))
		if isDir {
			btnOpen.SetEnabled(true)
			if !dlg.selectDir {
				btnSelect.SetEnabled(false)
				return
			}
		}
		if !isDir {
			btnOpen.SetEnabled(false)
		}

		btnSelect.SetEnabled(true)
		if (isDir && dlg.selectDir) || !isDir {
			dlg.edFile.SetTitle(item)
		} else {
			dlg.edFile.Clear()
		}
	})

	btnOpen.OnClick(func(ev Event) {
		s := dlg.listBox.SelectedItemText()
		if s != ".." && (s == "" || !strings.HasSuffix(s, string(os.PathSeparator))) {
			return
		}

		if s == ".." {
			dlg.pathUp()
		} else {
			dlg.pathDown(s)
		}
	})

	dlg.edFile.OnChange(func(ev Event) {
		s := ""
		lowCurrText := strings.ToLower(dlg.listBox.SelectedItemText())
		lowEditText := strings.ToLower(dlg.edFile.Title())
		if !strings.HasPrefix(lowCurrText, lowEditText) {
			itemIdx := dlg.listBox.PartialFindItem(dlg.edFile.Title(), false)
			if itemIdx == -1 {
				if dlg.mustExist {
					btnSelect.SetEnabled(false)
				}
				return
			}

			s, _ = dlg.listBox.Item(itemIdx)
			dlg.listBox.SelectItem(itemIdx)
		} else {
			s = dlg.listBox.SelectedItemText()
		}

		isDir := strings.HasSuffix(s, string(os.PathSeparator))
		equal := strings.EqualFold(s, dlg.edFile.Title()) ||
			strings.EqualFold(s, dlg.edFile.Title()+string(os.PathSeparator))

		btnOpen.SetEnabled(isDir || s == "..")
		if isDir {
			btnSelect.SetEnabled(dlg.selectDir && (equal || !dlg.mustExist))
			return
		}

		btnSelect.SetEnabled(equal || !dlg.mustExist)
	})

	dlg.edFile.OnKeyPress(func(key term.Key, c rune) bool {
		if key == term.KeyArrowUp {
			idx := dlg.listBox.SelectedItem()
			if idx > 0 {
				dlg.listBox.SelectItem(idx - 1)
				dlg.updateEditBox()
			}
			return true
		}
		if key == term.KeyArrowDown {
			idx := dlg.listBox.SelectedItem()
			if idx < dlg.listBox.ItemCount()-1 {
				dlg.listBox.SelectItem(idx + 1)
				dlg.updateEditBox()
			}
			return true
		}

		return false
	})

	dlg.listBox.OnKeyPress(func(key term.Key) bool {
		if key == term.KeyBackspace || key == term.KeyBackspace2 {
			dlg.pathUp()
			return true
		}

		if key == term.KeyCtrlM {
			s := dlg.listBox.SelectedItemText()
			if s == ".." {
				dlg.pathUp()
				return true
			}

			if strings.HasSuffix(s, string(os.PathSeparator)) {
				dlg.pathDown(s)
				return true
			}

			return false
		}

		return false
	})

	dlg.curDir.SetTitle(dlg.currPath)
	dlg.populateFiles()
	dlg.selectFirst()

	ActivateControl(dlg.View, dlg.listBox)
	return dlg
}

// OnClose sets the callback that is called when the
// dialog is closed
func (d *FileSelectDialog) OnClose(fn func()) {
	WindowManager().BeginUpdate()
	defer WindowManager().EndUpdate()
	d.onClose = fn
}
