package clui

import (
	"testing"
	мИнт "./пакИнтерфейсы"
)

func TestListBox(t *testing.T) {
	width, height := 10, 5
	lbox := CreateListBox(nil, width, height, мИнт.Fixed)

	w, h := lbox.Size()
	if w != width {
		t.Errorf("Width invalid: %v instead of %v", w, width)
	}
	if h != height {
		t.Errorf("Width invalid: %v instead of %v", height, height)
	}

	lbox.AddItem("Item1")
	lbox.AddItem("Item2")
	lbox.AddItem("Item3")

	if lbox.ItemCount() != 3 {
		t.Errorf("Item count must be %v instead of %v", 3, lbox.ItemCount())
	}

	n := lbox.FindItem("Item2", false)
	if n != 1 {
		t.Errorf("Item2 is not found")
	}
	n = lbox.FindItem("item2", true)
	if n != 1 {
		t.Errorf("item2 is not found")
	}
	lbox.SelectItem(n)
	str := lbox.SelectedItemText()
	if str != "Item2" {
		t.Errorf("The second item text must be %v, found %v", "Item2", str)
	}
	n = lbox.FindItem("item4", false)
	if n != -1 {
		t.Errorf("item4 should not be found")
	}
	lbox.RemoveItem(1)
	if lbox.ItemCount() != 2 {
		t.Errorf("After deleting an item the list box item count should decrease (%v)", lbox.ItemCount())
	}
	str = lbox.SelectedItemText()
	if str != "Item3" {
		t.Errorf("The second item text must be %v, found %v", "Item3", str)
	}
	n = lbox.FindItem("Item2", false)
	if n != -1 {
		t.Errorf("Item2 should not be found")
	}
	n = lbox.FindItem("item3", true)
	if n != 1 {
		t.Errorf("Item3 should #%v instead of %v", 1, n)
	}
	lbox.Clear()
	if lbox.ItemCount() != 0 {
		t.Errorf("Clear failed")
	}
}
