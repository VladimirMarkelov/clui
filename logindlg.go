package clui

const (
	LoginOk = iota
	LoginCanceled
	LoginInvalid
)

// LoginDialog is a login dialog with fields to enter user name and password
// Public properties:
//   * Username - login entered by a user
//   * Password - password entered by a user
//   * Action - how the dialog was closed:
//     - LoginOk - button "OK" was clicked
//     - LoginCanceled - button "Cancel" was clicked or dialog was dismissed
//     - LoginInvalid - invalid credentials were entered. This value appears
//         only in case of callback is used and button "OK" is clicked
//         while entered username or password is incorrect
type LoginDialog struct {
	View     *Window
	Username string
	Password string
	Action   int

	result  int
	onClose func()
	onCheck func(string, string) bool
}

// LoginDialog creates a new login dialog
//  * title - custom dialog title
//  * userName - initial username. Maybe useful if you want to implement
//     a feature "remember me"
// The active control depends on userName: if it is empty then the cursor is
//  in Username field, and in Password field otherwise.
// By default the dialog is closed when button "OK" is clicked. But if you set
//  OnCheck callback the dialog closes only if callback returns true or
//  button "Cancel" is clicked. This is helpful if you do not want to recreate
//  the dialog after every incorrect credentials. So, you define a callback
//  that checks whether pair of Usename and Password is correct and then the
//  button "OK" closed the dialog only if the callback returns true. If the
//  credentials are not valid, then the dialog shows a warning. The warning
//  automatically disappears when a user starts typing in Password or Username
//  field.
func CreateLoginDialog(title, userName string) *LoginDialog {
	dlg := new(LoginDialog)

	dlg.View = AddWindow(15, 8, 10, 4, title)
	WindowManager().BeginUpdate()
	defer WindowManager().EndUpdate()

	dlg.View.SetModal(true)
	dlg.View.SetPack(Vertical)

	userfrm := CreateFrame(dlg.View, 1, 1, BorderNone, Fixed)
	userfrm.SetPaddings(1, 1)
	userfrm.SetPack(Horizontal)
	userfrm.SetGaps(1, 0)
	CreateLabel(userfrm, AutoSize, AutoSize, "User name", Fixed)
	edUser := CreateEditField(userfrm, 20, userName, 1)

	passfrm := CreateFrame(dlg.View, 1, 1, BorderNone, 1)
	passfrm.SetPaddings(1, 1)
	passfrm.SetPack(Horizontal)
	passfrm.SetGaps(1, 0)
	CreateLabel(passfrm, AutoSize, AutoSize, "Password", Fixed)
	edPass := CreateEditField(passfrm, 20, "", 1)
	edPass.SetPasswordMode(true)

	filler := CreateFrame(dlg.View, 1, 1, BorderNone, 1)
	filler.SetPack(Horizontal)
	lbRes := CreateLabel(filler, AutoSize, AutoSize, "", 1)

	blist := CreateFrame(dlg.View, 1, 1, BorderNone, Fixed)
	blist.SetPack(Horizontal)
	blist.SetPaddings(1, 1)
	btnOk := CreateButton(blist, 10, 4, "OK", Fixed)
	btnCancel := CreateButton(blist, 10, 4, "Cancel", Fixed)

	btnCancel.OnClick(func(ev Event) {
		WindowManager().DestroyWindow(dlg.View)
		WindowManager().BeginUpdate()
		dlg.Action = LoginCanceled
		closeFunc := dlg.onClose
		WindowManager().EndUpdate()
		if closeFunc != nil {
			closeFunc()
		}
	})

	btnOk.OnClick(func(ev Event) {
		if dlg.onCheck != nil && !dlg.onCheck(edUser.Title(), edPass.Title()) {
			lbRes.SetTitle("Invalid username or password")
			dlg.Action = LoginInvalid
			return
		}

		dlg.Action = LoginOk
		if dlg.onCheck == nil {
			dlg.Username = edUser.Title()
			dlg.Password = edPass.Title()
		}

		WindowManager().DestroyWindow(dlg.View)
		WindowManager().BeginUpdate()

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

	edUser.OnChange(func(ev Event) {
		lbRes.SetTitle("")
	})
	edPass.OnChange(func(ev Event) {
		lbRes.SetTitle("")
	})

	if userName == "" {
		ActivateControl(dlg.View, edUser)
	} else {
		ActivateControl(dlg.View, edPass)
	}
	return dlg
}

// OnClose sets the callback that is called when the
// dialog is closed
func (d *LoginDialog) OnClose(fn func()) {
	WindowManager().BeginUpdate()
	defer WindowManager().EndUpdate()
	d.onClose = fn
}

// OnCheck sets the callback that is called when the
// button "OK" is clicked. The dialog sends to the callback two arguments:
// username and password. The callback validates the arguments and if
// the credentials are valid it returns true. That means the dialog can be
// closed. If the callback returns false then the dialog remains on the screen.
func (d *LoginDialog) OnCheck(fn func(string, string) bool) {
	WindowManager().BeginUpdate()
	defer WindowManager().EndUpdate()
	d.onCheck = fn
}
