package spicefyne

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	resources "github.com/Shells-com/shells-go/res"
)

type control struct {
	main            *fyne.Container
	setW            chan int
	visible         bool
	parent          *SpiceFyne
	sidebarPosition string
	tgtW            float32

	muteIconContainer *fyne.Container
	sidebar           *fyne.Container
	expBtn            *widget.Button
	expBtnContainer   *fyne.Container
	zoomBtn           *widget.Button
	muteIcon          *widget.Icon
}

func newControl(w *SpiceFyne) *control {
	res := &control{
		setW:   make(chan int),
		parent: w,
	}

	pos := w.a.Preferences().String(PreferencesKeyPosition)
	if len(pos) == 0 {
		pos = "left"
	}
	res.sidebarPosition = pos

	// expand button
	res.expBtn = widget.NewButtonWithIcon("", resources.ShellsIcon, res.toggle)
	res.expBtn.Importance = widget.LowImportance
	res.expBtn.Resize(fyne.NewSize(SideBarIconButton, SideBarIconButton))

	sidebar := buildLayout(res)
	res.sidebar = sidebar

	expBtnContainer := container.NewCenter(res.expBtn)
	res.expBtnContainer = expBtnContainer

	if res.sidebarPosition == "left" {
		c := container.NewHBox(sidebar, expBtnContainer)
		res.main = c
	} else {
		c := container.NewHBox(expBtnContainer, sidebar)
		res.main = c
	}

	isMuted := w.a.Preferences().Bool(PreferencesKeyMute)
	muteIcon := widget.NewIcon(theme.VolumeMuteIcon())
	muteIcon.Hidden = !isMuted
	res.muteIcon = muteIcon
	lmi := NewMuteIconLayout()
	res.muteIconContainer = container.New(lmi, muteIcon)

	go res.anim()

	return res
}

func (c *control) anim() {
	t := time.NewTicker(10 * time.Millisecond)
	tgtW := (0 - c.main.Size().Width) + SideBarWidth

	// anim task
	for {
		select {
		case w := <-c.setW:
			tgtW = float32(w)
			t.Reset(10 * time.Millisecond)
		case <-t.C:
			// get width
			p := c.main.Position()
			if p.X == tgtW {
				t.Stop()
				break
			}

			// update width
			if c.sidebarPosition == "left" {
				if p.X > tgtW {
					// need to reduce
					p.X -= 7
					if p.X < tgtW {
						p.X = tgtW
					}
					c.main.Move(p)
				} else {
					p.X += 7
					if p.X > tgtW {
						p.X = tgtW
					}
					c.main.Move(p)
				}
			} else {
				if p.X < tgtW {
					// need to increase
					p.X += 7
					if p.X > tgtW {
						p.X = tgtW
					}
					c.main.Move(p)
				} else {
					p.X -= 7
					if p.X < tgtW {
						p.X = tgtW
					}
					c.main.Move(p)
				}
			}

		}
	}
}

// set visibility and calculate where the main container will go with the anim
func (c *control) toggle() {
	c.visible = !c.visible
	if c.visible {
		if c.sidebarPosition == "left" {
			c.setW <- 0
		} else {
			c.setW <- int(c.parent.Size().Width - SideBarIconButton - SideBarWidth)
		}
	} else {
		if c.sidebarPosition == "left" {
			c.setW <- int((0 - c.main.Size().Width) + SideBarIconButton)
		} else {
			c.setW <- int(c.parent.size.Width - SideBarIconButton)
		}
	}
}

func (c *control) zoom() {
	if c.parent.zoomMode {
		c.parent.zoomMode = false
	} else {
		c.parent.zoomMode = true
		c.toggle()
	}
}

func (c *control) mute() {
	c.parent.Client.ToggleMute()
	isMuted := c.parent.Client.GetMute()
	c.parent.a.Preferences().SetBool(PreferencesKeyMute, isMuted)
	c.muteIcon.Hidden = !isMuted
}

func (c *control) changePosition(direction string) {
	c.sidebarPosition = direction
	c.parent.a.Preferences().SetString(PreferencesKeyPosition, direction)
	ps := c.main.Position()
	size := c.parent.Size()
	// to set the button on the right/left side, we must reorder the sidebar/button order in the container
	c.main.Remove(c.sidebar)
	c.main.Remove(c.expBtnContainer)
	if direction == "left" {
		ps.X = 0
		c.main.Add(c.sidebar)
		c.main.Add(c.expBtnContainer)
	} else {
		ps.X = size.Width - SideBarWidth - SideBarIconButton
		c.main.Add(c.expBtnContainer)
		c.main.Add(c.sidebar)
	}
	// move the main container to the right/left side in the open state
	c.main.Move(ps)
}
