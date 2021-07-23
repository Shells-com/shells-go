package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type enterEntry struct {
	widget.Entry

	onEnter func()
}

func newEnterEntry() *enterEntry {
	entry := &enterEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *enterEntry) TypedKey(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyReturn:
		if e.onEnter != nil {
			e.onEnter()
			return
		}
		fallthrough
	default:
		e.Entry.TypedKey(key)
	}
}
