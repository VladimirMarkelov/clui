package clui

// RadioGroup is non-interactive invisible object. It manages
// set of Radio buttons: at a time no more than one radio
// button from a group can be selected
type RadioGroup struct {
	items []*Radio
}

// NewRadioGroup creates a new RadioGroup
func CreateRadioGroup() *RadioGroup {
	c := new(RadioGroup)
	c.items = make([]*Radio, 0)
	return c
}

// Selected returns the number of currently selected radio
// button inside the group or -1 if no button is selected
func (c *RadioGroup) Selected() int {
	selected := -1

	for id, item := range c.items {
		if item.Selected() {
			selected = id
			break
		}
	}

	return selected
}

// SelectItem makes the radio selected. The function returns false
// if it failed to find the radio in the radio group
func (c *RadioGroup) SelectItem(r *Radio) bool {
	found := false

	for _, item := range c.items {
		if item == r {
			found = true
			item.SetSelected(true)
		} else {
			item.SetSelected(false)
		}
	}

	return found
}

// SetSelected selects the radio by its number. The function
// returns false if the number is invalid for the radio group
func (c *RadioGroup) SetSelected(id int) bool {
	found := false
	if id < 0 || id >= len(c.items) {
		return false
	}

	for idx, item := range c.items {
		if idx == id {
			found = true
			item.SetSelected(true)
		} else {
			item.SetSelected(false)
		}
	}

	return found
}

// AddItem add a new radio button to group
func (c *RadioGroup) AddItem(r *Radio) {
	c.items = append(c.items, r)
	r.SetGroup(c)
}
