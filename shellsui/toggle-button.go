package shellsui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type ToggleButton struct {
	widget.Button
	text        string
	icon        fyne.Resource
	textToggled string
	iconToggled fyne.Resource
	toggled     bool
}

func NewToggleButton(text string, icon fyne.Resource, textToggled string, iconToggled fyne.Resource, initState bool, tapped func()) *ToggleButton {
	btn := &ToggleButton{
		text:        text,
		icon:        icon,
		textToggled: textToggled,
		iconToggled: iconToggled,
	}
	btn.ExtendBaseWidget(btn)
	if initState {
		btn.SetText(textToggled)
		btn.SetIcon(iconToggled)
	} else {
		btn.SetText(text)
		btn.SetIcon(icon)
	}
	btn.Button.OnTapped = tapped
	return btn
}

func (b *ToggleButton) Tapped(*fyne.PointEvent) {
	if b.toggled {
		b.Button.SetText(b.text)
		b.Button.SetIcon(b.icon)
	} else {
		b.Button.SetText(b.textToggled)
		b.Button.SetIcon(b.iconToggled)
	}
	b.toggled = !b.toggled
	b.OnTapped()
}
