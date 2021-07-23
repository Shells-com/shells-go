package spicefyne

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type autoLayout struct {
	direction      *string
	sidebarVisible *bool
}

func (l *autoLayout) Layout(o []fyne.CanvasObject, s fyne.Size) {
	first := true
	for _, obj := range o {
		if first {
			obj.Resize(s)
			first = false
			continue
		}
		obj.Resize(fyne.Size{Width: obj.Size().Width, Height: s.Height})
		h := obj.Size().Height
		var c fyne.Position
		var vsize float32
		if *l.sidebarVisible {
			vsize = float32(SideBarWidth)
		} else {
			vsize = float32(0)
		}
		if *l.direction == "left" {
			c = fyne.NewPos(SideBarIconButton-obj.Size().Width+vsize, (s.Height-h)/2)
		} else {
			c = fyne.NewPos(s.Width-SideBarIconButton-vsize, (s.Height-h)/2)
		}
		obj.Move(c)
	}
}

func (l *autoLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.Size{Width: 320, Height: 200}
}

type centerHLayout struct {
}

func NewCenterHLayout() fyne.Layout {
	return &centerHLayout{}
}

func (c *centerHLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, child := range objects {
		childMin := child.MinSize()
		child.Resize(childMin)
		child.Move(fyne.NewPos(SideBarPadding, 0))
	}
}

func (c *centerHLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		minSize = minSize.Max(child.MinSize())
	}

	return minSize
}

type sidebarLayout struct {
}

func NewSidebarLayout() fyne.Layout {
	return &sidebarLayout{}
}

func (s *sidebarLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	total := float32(0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		total += child.MinSize().Height
	}

	x, y := float32(0), float32(0)

	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		height := child.MinSize().Height
		child.Move(fyne.NewPos(x, y))
		y += theme.Padding() + height
		child.Resize(fyne.NewSize(SideBarWidth-(SideBarPadding*2), height))
	}
}

func (s *sidebarLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	addPadding := false

	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		minSize.Width = fyne.Max(child.MinSize().Width, minSize.Width)
		minSize.Height += child.MinSize().Height
		if addPadding {
			minSize.Height += theme.Padding()
		}
		addPadding = true

	}

	return minSize
}

// as to be use as child of fyne.Window
type muteIconLayout struct {
}

func NewMuteIconLayout() fyne.Layout {
	return &muteIconLayout{}
}

func (s *muteIconLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, child := range objects {
		child.Resize(fyne.NewSize(MuteIconSize, MuteIconSize))
		x := size.Width - MuteIconSize - MuteIconMargin
		y := size.Height - MuteIconSize - MuteIconMargin
		child.Move(fyne.NewPos(x, y))
	}
}

func (s *muteIconLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(MuteIconSize, MuteIconSize)
}
