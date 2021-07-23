package spicefyne

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Shells-com/shells-go/shellsui"
)

func buildItem(labelText string, btnOptions btnOptions) (*widget.Button, *fyne.Container) {
	t := widget.NewLabelWithStyle(labelText, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	t.Wrapping = fyne.TextWrapWord

	var b *widget.Button

	if btnOptions.icon != nil {
		b = widget.NewButtonWithIcon(btnOptions.text, btnOptions.icon, btnOptions.cb)
	} else {
		b = widget.NewButton(btnOptions.text, btnOptions.cb)
	}

	lh := layout.NewHBoxLayout()
	lh.Layout([]fyne.CanvasObject{t, b}, fyne.Size{Width: SideBarWidth})

	sbl := NewSidebarLayout()
	sbc := container.New(sbl, t, b, shellsui.GetMarginRectangleWithHeight(10))
	return b, container.NewVBox(sbc)
}

func buildItemWithToggle(labelText string, btnOptions btnOptions) (*shellsui.ToggleButton, *fyne.Container) {
	t := widget.NewLabelWithStyle(labelText, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	t.Wrapping = fyne.TextWrapWord

	b := shellsui.NewToggleButton(btnOptions.text, btnOptions.icon, btnOptions.textToggled, btnOptions.iconToggled, btnOptions.initState, btnOptions.cb)

	lh := layout.NewHBoxLayout()
	lh.Layout([]fyne.CanvasObject{t, b}, fyne.Size{Width: SideBarWidth})

	sbl := NewSidebarLayout()
	sbc := container.New(sbl, t, b, shellsui.GetMarginRectangleWithHeight(10))
	return b, container.NewVBox(sbc)
}

func build2IconButtonsItem(labelText string, btnIcon1 fyne.Resource, cb1 func(), btnIcon2 fyne.Resource, cb2 func()) *fyne.Container {
	t := widget.NewLabelWithStyle(labelText, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	t.Wrapping = fyne.TextWrapWord

	b1 := widget.NewButtonWithIcon("", btnIcon1, cb1)
	b1.Importance = widget.MediumImportance

	b2 := widget.NewButtonWithIcon("", btnIcon2, cb2)
	b2.Importance = widget.MediumImportance

	c2b := container.NewCenter(container.New(layout.NewHBoxLayout(), b1, shellsui.GetMarginRectangleWithWidth(20), b2))

	lh := layout.NewHBoxLayout()
	lh.Layout([]fyne.CanvasObject{t, c2b}, fyne.Size{Width: SideBarWidth})

	sbl := NewSidebarLayout()
	sbc := container.New(sbl, t, c2b, shellsui.GetMarginRectangleWithHeight(10))
	return container.NewVBox(sbc)
}

func buildLayout(c *control) *fyne.Container {
	zoomButton, zoom := buildItem("Zoom", btnOptions{text: "Zoom", cb: c.zoom})
	c.zoomBtn = zoomButton

	_, resetZoom := buildItem("Reset Zoom", btnOptions{text: "Reset", cb: c.parent.resetZoom})
	_, exit := buildItem("Return to Shells list", btnOptions{text: "Exit", cb: c.parent.terminate})
	h := c.parent.Size().Height

	_, mute := buildItemWithToggle("Mute", btnOptions{
		text:        "Mute",
		icon:        theme.VolumeMuteIcon(),
		textToggled: "Unmute",
		iconToggled: theme.VolumeUpIcon(),
		cb:          c.mute,
		initState:   c.parent.a.Preferences().Bool(PreferencesKeyMute),
	})

	positions := build2IconButtonsItem("Sidebar Position", theme.NavigateBackIcon(), func() {
		c.changePosition("left")
	}, theme.NavigateNextIcon(), func() {
		c.changePosition("right")
	})

	lv := layout.NewVBoxLayout()
	lv.Layout([]fyne.CanvasObject{zoom, resetZoom, mute, positions, exit}, fyne.Size{Height: h})
	cv := container.New(lv, zoom, resetZoom, mute, positions, exit)

	cc := container.New(NewCenterHLayout(), cv)

	bg := canvas.NewRectangle(color.White)
	bg.SetMinSize(fyne.NewSize(SideBarWidth, h))

	return container.NewMax(bg, cc)
}

type btnOptions struct {
	text        string
	icon        fyne.Resource
	textToggled string
	iconToggled fyne.Resource
	cb          func()
	initState   bool
}
