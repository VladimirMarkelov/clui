package clui

import (
	"fmt"
	"github.com/VladimirMarkelov/termbox-go"
	xs "github.com/huandu/xstrings"
	"log"
	"os"
	"strings"
)

// Internal representation of Window object that keeps and manages control inside it.
type view struct {
	posX, posY    int
	width, height int
	minW, minH    int // min size constraints
	maxW, maxH    int // max size constraints
	id            WinId
	title         string
	borderStyle   BorderStyle
	icons         BorderIcon
	enabled       bool
	active        bool
	canvas        *FrameBuffer
	controls      []Control
	parent        *Composer
	lastCtrlId    WinId //last Id used for control

	originals      map[WinId]Coord
	originalWidth  int
	originalHeight int

	// pack support
	layout                LayoutType
	pack                  PackType
	lockUpdate            bool // it is true while new controls are adding
	padSide, padTop       int
	padX, padY            int
	scale                 int
	children              []WinId
	packer                Packer
	lastX, lastY          int
	currWidth, currHeight int

	// helpers
	logger *log.Logger
}

func NewView(composer *Composer, id WinId, posX, posY, width, height int, title string) *view {
	d := new(view)
	d.minW, d.minH = 10, 5

	if width < d.minW {
		width = d.minW
	}
	if height < d.minH {
		height = d.minH
	}

	d.SetTitle(title)
	d.SetSize(width, height)
	d.SetPos(posX, posY)
	d.SetEnabled(true)
	d.controls = make([]Control, 0)
	d.originals = make(map[WinId]Coord)
	d.parent = composer
	d.id = id
	d.active = false
	d.originalWidth = width
	d.originalHeight = height
	d.lastCtrlId = 0
	d.children = make([]WinId, 0)
	d.layout = LayoutManual

	d.lastX, d.lastY = -1, -1
	d.currWidth, d.currHeight = -1, -1
	d.padSide, d.padTop, d.padX, d.padY = 0, 0, 1, 0
	d.borderStyle = BorderSingle

	file, _ := os.OpenFile("debug.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	d.logger = log.New(file, fmt.Sprintf("WND[%v]", title), log.Ldate|log.Ltime|log.Lshortfile)

	return d
}

func (d *view) SetTitle(title string) {
	d.title = title
}

func (d *view) GetTitle() string {
	return d.title
}

func (d *view) SetSize(width, height int) {
	if width > 1000 || width < d.minW {
		panic(fmt.Sprintf("Invalid width: %v", width))
	}
	if height > 200 || height < d.minH {
		panic(fmt.Sprintf("Invalid height: %v", height))
	}

	d.width = width
	d.height = height

	d.canvas = NewFrameBuffer(width, height)
}

func (d *view) SetPos(x, y int) {
	d.posX = x
	d.posY = y
}

func (d *view) GetId() WinId {
	return d.id
}

func (d *view) GetPos() (int, int) {
	return d.posX, d.posY
}

func (d *view) GetSize() (int, int) {
	return d.width, d.height
}

func (d *view) GetBorderStyle() BorderStyle {
	return d.borderStyle
}

func (d *view) DrawControls() {
	if d.packer != nil {
		d.packer.Redraw(d)
	}
	for _, ctrl := range d.controls {
		if ctrl.GetVisible() {
			ctrl.Redraw(d)
		}
	}
}
func (d *view) DrawDecoration() {
	if d.canvas == nil {
		return
	}

	tm := d.parent.GetThemeManager()
	fg, bg := ColorWhite, ColorBlack
	if d.active {
		fg = tm.GetSysColor(ColorActiveText)
		bg = tm.GetSysColor(ColorViewBack)
	} else {
		fg = tm.GetSysColor(ColorInactiveText)
		bg = tm.GetSysColor(ColorViewBack)
	}

	d.canvas.DrawBorder(d, tm, fg, bg)
	d.canvas.DrawBorderIcons(d, tm, fg, bg)
	d.canvas.DrawTitle(d, fg, bg)
}

func (d *view) Redraw() {
	if d.canvas == nil {
		return
	}

	tm := d.parent.GetThemeManager()
	bg := tm.GetSysColor(ColorViewBack)

	d.canvas.Clear(bg)
	d.DrawDecoration()
	d.DrawControls()
}

func (d *view) GetBorderIcons() BorderIcon {
	return d.icons
}

func (d *view) SetBorderIcons(icons BorderIcon) {
	if d.icons != icons {
		d.icons = icons
		d.DrawDecoration()
	}
}

func (d *view) isInside(screenX, screenY int) bool {
	if screenX >= d.posX && screenX < d.posX+d.width && screenY >= d.posY && screenY < d.posY+d.height {
		return true
	} else {
		return false
	}
}

func (d *view) GetScreenSymbol(screenX, screenY int) (Symbol, bool) {
	if !d.isInside(screenX, screenY) {
		return Symbol{ch: ' '}, false
	} else {
		return d.canvas.GetSymbol(screenX-d.posX, screenY-d.posY)
	}
}

func (d *view) GetEnabled() bool {
	return d.enabled
}

func (d *view) SetEnabled(enable bool) {
	d.enabled = enable
}

func (d *view) GetActive() bool {
	return d.active
}

func (d *view) SetActive(active bool) {
	d.active = active
}

func (d *view) HitTest(screenX, screenY int) HitResult {
	if !d.isInside(screenX, screenY) {
		return HitOutside
	}

	if screenX == d.posX {
		if screenY == d.posY {
			return HitTopLeft
		} else if screenY == d.posY+d.height-1 {
			return HitBottomLeft
		} else {
			return HitLeftBorder
		}
	}

	if screenX == d.posX+d.width-1 {
		if screenY == d.posY {
			return HitTopRight
		} else if screenY == d.posY+d.height-1 {
			return HitBottomRight
		} else {
			return HitRightBorder
		}
	}

	if screenY == d.posY+d.height-1 {
		return HitBottomBorder
	}

	if screenY == d.posY {
		dx := -3
		if d.icons&IconClose != 0 {
			if screenX == d.posX+d.width+dx {
				return HitButtonClose
			}
			dx--
		}
		if d.icons&IconBottom != 0 {
			if screenX == d.posX+d.width+dx {
				return HitButtonBottom
			}
			dx--
		}

		return HitHeader
	}

	return HitInside
}

func (d *view) GetConstraints() (int, int) {
	if d.packer != nil {
		minW, minH := d.packer.GetConstraints()
		if minW < d.minW {
			minW = d.minW
		}
		if minH < d.minH {
			minH = d.minH
		}
		return minW, minH
	}
	return d.minW, d.minH
}

func (d *view) AddControl(control Control) WinId {
	d.controls = append(d.controls, control)

	id := control.GetId()
	x, y := control.GetPos()
	w, h := control.GetSize()

	c := Coord{x: x, y: y, w: w, h: h}
	d.originals[id] = c

	d.logger.Printf("Control %v got id %v (%vx%v)", control.GetText(), id, w, h)

	return id
}

func (d *view) RemoveControl(control Control) {
	id := control.GetId()
	_, ok := d.originals[id]
	if ok {
		delete(d.originals, id)
	}

	newList := make([]Control, 0)
	for _, ctrl := range d.controls {
		if ctrl.GetId() != id {
			newList = append(newList, ctrl)
		}
	}
	d.controls = newList

	if len(d.children) > 0 {
		newKids := make([]WinId, 0)
		for _, cid := range d.children {
			if cid != id {
				newKids = append(newKids, cid)
			}
		}

		d.children = newKids
	}
}

func (d *view) GetControl(id WinId) Control {
	for _, ctrl := range d.controls {
		if ctrl.GetId() == id {
			return ctrl
		}
	}

	return nil
}

// ------------ Canvas methods ---------------
func (d *view) DrawText(x, y, w int, text string, fg, bg Color) {
	if text == "" || w < 1 {
		return
	}
	if xs.Len(text) > w {
		text = xs.Slice(text, 0, w)
	}
	d.canvas.DrawText(d, x, y, text, fg, bg)
}

func (d *view) DrawVerticalText(x, y, h int, text string, fg, bg Color) {
	if text == "" || h < 1 {
		return
	}

	for idx, r := range text {
		if idx >= h {
			break
		}

		d.canvas.DrawText(d, x, y+idx, string(r), fg, bg)
	}
}

func (d *view) DrawAlignedText(x, y, w int, text string, fg, bg Color, align Align) {
	if text == "" || w < 1 {
		return
	}
	length := xs.Len(text)
	if length < w {
		if align == AlignCenter {
			d.DrawText(x+int((w-length)/2), y, w, text, fg, bg)
		} else if align == AlignLeft {
			d.DrawText(x, y, w, text, fg, bg)
		} else {
			d.DrawText(x+w-length, y, length, text, fg, bg)
		}
	} else {
		str := ""
		if align == AlignCenter {
			dx := int((length - w) / 2)
			str = xs.Slice(text, dx, dx+w)
		} else if align == AlignLeft {
			str = xs.Slice(text, 0, w)
		} else {
			str = xs.Slice(text, length-w, -1)
		}
		d.DrawText(x, y, w, str, fg, bg)
	}
}

func (d *view) DrawRune(x, y int, r rune, fg, bg Color) {
	d.canvas.DrawText(d, x, y, string(r), fg, bg)
}

func (d *view) DrawFrame(x, y, width, height int, bs BorderStyle, fg, bg Color) {
	if bs == BorderNone {
		return
	}

	tm := d.parent.GetThemeManager()

	var cH, cV, cUL, cUR, cDL, cDR rune
	if bs == BorderSingle {
		cH = tm.GetSysObject(ObjSingleBorderHLine)
		cV = tm.GetSysObject(ObjSingleBorderVLine)
		cUL = tm.GetSysObject(ObjSingleBorderULCorner)
		cUR = tm.GetSysObject(ObjSingleBorderURCorner)
		cDL = tm.GetSysObject(ObjSingleBorderDLCorner)
		cDR = tm.GetSysObject(ObjSingleBorderDRCorner)
	} else {
		cH = tm.GetSysObject(ObjDoubleBorderHLine)
		cV = tm.GetSysObject(ObjDoubleBorderVLine)
		cUL = tm.GetSysObject(ObjDoubleBorderULCorner)
		cUR = tm.GetSysObject(ObjDoubleBorderURCorner)
		cDL = tm.GetSysObject(ObjDoubleBorderDLCorner)
		cDR = tm.GetSysObject(ObjDoubleBorderDRCorner)
	}

	if width > 1 && height > 1 {
		d.DrawRune(x, y, cUL, fg, bg)
		d.DrawRune(x, y+height-1, cDL, fg, bg)
		d.DrawRune(x+width-1, y, cUR, fg, bg)
		d.DrawRune(x+width-1, y+height-1, cDR, fg, bg)
		for dx := 1; dx < width-1; dx++ {
			d.DrawRune(x+dx, y, cH, fg, bg)
			d.DrawRune(x+dx, y+height-1, cH, fg, bg)
		}
		for dy := 1; dy < height-1; dy++ {
			d.DrawRune(x, y+dy, cV, fg, bg)
			d.DrawRune(x+width-1, y+dy, cV, fg, bg)
		}
	} else if width == 1 {
		for dy := 0; dy < height; dy++ {
			d.DrawRune(x, y+dy, cV, fg, bg)
		}
	} else if height == 1 {
		for dx := 0; dx < width; dx++ {
			d.DrawRune(x+dx, y, cH, fg, bg)
		}
	}
}

func (d *view) ClearRect(x, y, w, h int, bg Color) {
	if w < 1 || h < 1 {
		return
	}

	s := strings.Repeat(" ", w)

	for i := y; i < y+h; i++ {
		d.canvas.DrawText(d, x, i, s, ColorWhite, bg)
	}
}

func (d *view) SetCursorPos(control Control, x, y int) {
	if !d.active {
		return
	}

	xc, yc := -1, -1
	wc, hc := 0, 0

	for _, ctrl := range d.controls {
		if ctrl.GetId() == control.GetId() {
			xc, yc = ctrl.GetPos()
			wc, hc = ctrl.GetSize()
			break
		}
	}

	if xc >= 0 && yc >= 0 && x >= 0 && x < wc && y >= 0 && y < hc {
		wx, wy := d.mapViewToScreen(xc, yc)
		d.parent.SetCursorPos(wx+x, wy+y)
	}
}

//---------------- internal -----------------------

func (d *view) mapViewToScreen(x, y int) (int, int) {
	wx, wy := d.GetPos()
	bs := d.GetBorderStyle()
	if bs != BorderNone {
		wx++
		wy++
	}

	return x + wx, y + wy
}

func (d *view) mapScreenToView(x, y int) (int, int) {
	wx, wy := d.GetPos()
	bs := d.GetBorderStyle()
	if bs != BorderNone {
		wx++
		wy++
	}

	return x - wx, y - wy
}

func (d *view) deactivateControls() {
	for _, ctrl := range d.controls {
		ctrl.SetActive(false)
	}
}

func (d *view) getActiveControl() Control {
	for _, ctrl := range d.controls {
		if ctrl.GetActive() {
			return ctrl
		}
	}

	return nil
}

func (d *view) activateNextControl(forward bool) bool {
	if len(d.controls) == 0 {
		return false
	}

	idx := -1
	for i := 0; i < len(d.controls); i++ {
		if d.controls[i].GetActive() {
			idx = i
			break
		}
	}

	if idx == -1 && forward {
		idx = len(d.controls)
	}

	var newidx, inc int
	if forward {
		newidx = idx + 1
		inc = 1
	} else {
		newidx = idx - 1
		inc = -1
	}

	for {
		if newidx == idx {
			break
		}

		if newidx < 0 {
			newidx = len(d.controls) - 1
		}
		if newidx >= len(d.controls) {
			newidx = 0
		}

		if d.controls[newidx].GetTabStop() && d.controls[newidx].GetVisible() && d.controls[newidx].GetEnabled() {
			break
		}

		newidx += inc
	}

	if idx == newidx {
		return false
	} else {
		d.ActivateControl(d.controls[newidx])
		return true
	}
}

func (d *view) ActivateControl(control Control) bool {
	id := control.GetId()
	activated := false
	for _, ctrl := range d.controls {
		if ctrl.GetId() == id {
			activated = true
			if !ctrl.GetActive() {
				event := Event{Type: EventActivate, X: 1}
				ctrl.ProcessEvent(event)
			}
			ctrl.SetActive(true)
		} else {
			if ctrl.GetActive() {
				event := Event{Type: EventActivate, X: 0}
				ctrl.ProcessEvent(event)
			}
			ctrl.SetActive(false)
		}
	}

	return activated
}

func (d *view) controlAtPos(screenX, screenY int) Control {
	posX, posY := d.mapScreenToView(screenX, screenY)

	for id := len(d.controls) - 1; id >= 0; id-- {
		ctrl := d.controls[id]

		if ctrl.GetVisible() {
			w, h := ctrl.GetSize()
			x, y := ctrl.GetPos()

			if posX >= x && posX < x+w && posY >= y && posY < y+h {
				return ctrl
			}
		}
	}

	return nil
}

func (d *view) recalculateManual() {
	winW, winH := d.GetSize()

	for _, ctrl := range d.controls {
		anchor := ctrl.GetAnchors()
		if anchor == AnchorNone {
			continue
		}

		orig, ok := d.originals[ctrl.GetId()]
		if !ok {
			d.logger.Printf("No originals for %v", ctrl.GetId())
			continue
		}

		newX, newY := orig.x, orig.y
		newW, newH := orig.w, orig.h
		dx := winW - d.originalWidth
		dy := winH - d.originalHeight

		if anchor&AnchorRight != 0 && anchor&AnchorLeft == 0 {
			// right side align
			newX += dx
		}
		if anchor&AnchorRight != 0 && anchor&AnchorLeft != 0 {
			// full width
			newW = orig.w + dx
		}
		if anchor&AnchorBottom != 0 && anchor&AnchorTop == 0 {
			// bottom align
			newY += dy
		}
		if anchor&AnchorTop != 0 && anchor&AnchorBottom != 0 {
			// full width
			newH = orig.h + dy
		}

		if newH > 0 && newW > 0 && newX >= 0 && newY >= 0 {
			ctrl.SetPos(newX, newY)
			ctrl.SetSize(newW, newH)
		}
	}
}

func (d *view) recalculateDynamic() {
	if d.packer != nil {
		newW, newH := d.GetSize()
		oldW, oldH := d.GetConstraints()
		dx, dy := newW-oldW, newH-oldH
		d.packer.ResizeChidren(dx, dy)
	}
}

func (d *view) recalculateControls() {
	if d.pack == PackFixed {
		d.recalculateManual()
	} else {
		d.recalculateDynamic()
		d.packer.RepositionChildren()
	}
}

func (d *view) hideAllExtraControls() {
	for _, ctrl := range d.controls {
		ctrl.HideChildren()
	}
}

func (d *view) ProcessEvent(ev Event) bool {
	switch ev.Type {
	case EventKey, EventMouseScroll, EventMouseClick, EventMousePress, EventMouseRelease, EventMouseMove, EventMouse:
		if ev.Type == EventKey && ev.Key == termbox.KeyTab {
			d.activateNextControl(ev.Mod&termbox.ModShift == 0)
			return true
		}
		if ev.Type == EventMouse || ev.Type == EventMouseClick {
			cunder := d.controlAtPos(ev.X, ev.Y)
			if cunder == nil {
				return true
			}
			d.ActivateControl(cunder)
		}
		ctrl := d.getActiveControl()
		x, y := ev.X, ev.Y
		copyEv := ev
		if ev.Type != EventMouseScroll {
			copyEv.X, copyEv.Y = d.mapScreenToView(x, y)
		}
		if ctrl != nil {
			ctrl.ProcessEvent(copyEv)
		}
	case EventActivate:
		if ev.X == 0 {
			d.parent.SetCursorPos(-1, -1)
		}
	case EventResize:
		d.hideAllExtraControls()
		d.recalculateControls()
	}

	return true
}

func (d *view) SendEvent(ev InternalEvent) {
	// now just send to composer
	ev.view = d.GetId()
	d.parent.SendEvent(ev)
}

func (d *view) Theme() *ThemeManager {
	return d.parent.GetThemeManager()
}

func (d *view) GetNextControlId() WinId {
	d.lastCtrlId++
	id := d.lastCtrlId
	return id
}

// ----- Packer ----------------------

func (d *view) AddPack(pt PackType) Packer {
	if d.packer != nil {
		panic("View already has a packer")
	}

	if len(d.controls) > 0 {
		panic("Cannot enable pack mode if a packer contains any control")
	}

	if pt != PackFixed {
		d.layout = LayoutDynamic
		d.pack = pt
		pid := d.GetNextControlId()
		d.packer = NewContainer(d, nil, pid, 0, 0, d.width, d.height, Props{})
		d.packer.SetPackType(pt)
	}

	d.lockUpdate = true

	return d.packer
}

func (d *view) PackEnd() (int, int) {
	if d.packer == nil || d.pack == PackFixed {
		panic("PackEnd can be used only if any dynamic packer is created before")
	}

	d.lockUpdate = false
	w, h := d.CalculateSize()

	// add window border
	w += 2
	h += 2
	d.logger.Printf("Pack ends with %vx%v", w, h)

	if w > 0 && h > 0 && (w > d.minW || h > d.minH) {
		d.SetConstraints(w, h)
	}

	d.recalculateControls()

	return w, h
}

func (d *view) CalculateSize() (int, int) {
	if d.packer != nil {
		return d.packer.CalculateSize()
	}

	return -1, -1
}

func (d *view) GetLayout() LayoutType {
	return d.layout
}

func (d *view) GetNextPosition() (int, int) {
	return d.lastX, d.lastY
}

func (d *view) SetNextPosition(x, y int) {
	d.lastX, d.lastY = x, y
}

func (d *view) SetConstraints(w, h int) {
	if w >= 10 {
		d.minW = w
	}
	if h >= 5 {
		d.minH = h
	}

	if d.width < w || d.height < h {
		d.SetSize(w, h)
	}
}

func (d *view) SetPaddings(pSide, pTop, pX, pY int) {
	if len(d.children) > 0 {
		panic("Cannot change padding if a child is added")
	}

	if pSide != DoNotChange {
		d.padSide = pSide
	}
	if pTop != DoNotChange {
		d.padTop = pTop
	}
	if pX != DoNotChange {
		d.padX = pX
	}
	if pY != DoNotChange {
		d.padY = pSide
	}
}

func (d *view) GetPaddings() (int, int, int, int) {
	return d.padSide, d.padTop, d.padX, d.padY
}

func (d *view) GetScale() int {
	return d.scale
}

func (d *view) SetScale(scale int) {
	// do nothing - does not make sense to set scale for view
}

//-------------- debug -----------------------------
func (d *view) Logger() *log.Logger {
	return d.logger
}
