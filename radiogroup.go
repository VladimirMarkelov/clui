package clui

type RadioGroup struct {
	items []*Radio
}

func NewRadioGroup() *RadioGroup {
	c := new(RadioGroup)
	c.items = make([]*Radio, 0)
	return c
}

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

func (c *RadioGroup) AddItem(r *Radio) {
	c.items = append(c.items, r)
	r.SetGroup(c)
}
