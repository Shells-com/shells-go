package shellsui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type SocialButton struct {
	widget.Button

	icon  fyne.Resource
	color color.Color

	tapAnim *fyne.Animation
}

func NewSocialButton(icon fyne.Resource, color color.Color, tapped func()) *SocialButton {
	btn := &SocialButton{
		color: color,
		icon:  icon,
	}
	btn.ExtendBaseWidget(btn)
	btn.SetIcon(icon)
	btn.OnTapped = tapped
	return btn
}

// CreateRenderer is a private method to Fyne which links this widget to its renderer
func (b *SocialButton) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)

	background := canvas.NewCircle(b.color)
	tapBG := canvas.NewCircle(color.Transparent)
	b.tapAnim = newButtonTapAnimation(tapBG, b)
	b.tapAnim.Curve = fyne.AnimationEaseOut
	objects := []fyne.CanvasObject{
		background,
		tapBG,
	}
	r := &buttonRenderer{
		objects:    objects,
		background: background,
		tapBG:      tapBG,
		button:     b,
		layout:     layout.NewHBoxLayout(),
	}
	r.updateIcon()
	r.background.FillColor = b.color
	r.background.Refresh()
	return r
}

type buttonRenderer struct {
	fyne.WidgetRenderer

	icon       *canvas.Image
	objects    []fyne.CanvasObject
	background *canvas.Circle
	tapBG      *canvas.Circle
	button     *SocialButton
	layout     fyne.Layout
}

func (r *buttonRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Layout the components of the button widget
func (r *buttonRenderer) Layout(size fyne.Size) {
	var inset fyne.Position
	bgSize := size
	inset = fyne.NewPos(theme.Padding()/2, theme.Padding()/2)
	bgSize = size.Subtract(fyne.NewSize(theme.Padding(), theme.Padding()))

	r.background.Move(inset)
	r.background.Resize(bgSize)

	hasIcon := r.icon != nil
	if !hasIcon {
		// Nothing to layout
		return
	}
	iconSize := fyne.NewSize(30, 30)

	// Icon Only
	r.icon.Move(alignedPosition(iconSize, size))
	r.icon.Resize(iconSize)

}

// MinSize calculates the minimum size of a button.
// This is based on the contained text, any icon that is set and a standard
// amount of padding added.
func (r *buttonRenderer) MinSize() (size fyne.Size) {
	iconSize := fyne.NewSize(40, 40)
	size.Width += iconSize.Width
	size.Height = iconSize.Height
	size = size.Add(r.padding())
	return
}

func (r *buttonRenderer) Refresh() {
	r.updateIcon()
	r.background.FillColor = r.button.color
	r.background.Refresh()
	r.Layout(r.button.Size())
	canvas.Refresh(r.button) //maybe wrong
}

func (r *buttonRenderer) padding() fyne.Size {
	return fyne.NewSize(theme.Padding()*4, theme.Padding()*4)
}

func (r *buttonRenderer) updateIcon() {
	if r.button.Icon != nil {
		if r.icon == nil {
			r.icon = canvas.NewImageFromResource(r.button.Icon)
			r.icon.FillMode = canvas.ImageFillContain
			r.SetObjects([]fyne.CanvasObject{r.background, r.tapBG, r.icon})
		}
		r.icon.Resource = r.button.Icon
		r.icon.Refresh()
		r.icon.Show()
	} else if r.icon != nil {
		r.icon.Hide()
	}
}

func (r *buttonRenderer) SetObjects(objects []fyne.CanvasObject) {
	r.objects = objects
}

func alignedPosition(objectSize, layoutSize fyne.Size) (pos fyne.Position) {
	pos.Y = (layoutSize.Height - objectSize.Height) / 2
	pos.X = (layoutSize.Width - objectSize.Width) / 2
	return pos
}

func newButtonTapAnimation(bg *canvas.Circle, w fyne.Widget) *fyne.Animation {
	return fyne.NewAnimation(canvas.DurationStandard, func(done float32) {
		mid := (w.Size().Width - theme.Padding()) / 2
		size := mid * done
		bg.Resize(fyne.NewSize(size*2, w.Size().Height-theme.Padding()))
		bg.Move(fyne.NewPos(mid-size, theme.Padding()/2))

		r, g, bb, a := theme.PressedColor().RGBA()
		aa := uint8(a)
		fade := aa - uint8(float32(aa)*done)
		bg.FillColor = &color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(bb), A: fade}
		canvas.Refresh(bg)
	})
}
