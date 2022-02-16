package spicefyne

import (
	"image"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/KarpelesLab/goclip"
	"github.com/Shells-com/spice"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type SpiceFyne struct {
	widget.BaseWidget
	*spice.Client

	w      fyne.Window
	a      fyne.App
	main   *spice.ChMain
	in     *spice.ChInputs
	output *canvas.Image
	cursor *canvas.Image
	width  int
	height int
	size   fyne.Size // size as shown. Will need to adapt
	zoom   float64
	init   bool
	clip   *goclip.Monitor
	ctrl   *control

	cursorObj *fyneCursor

	expanded bool

	// focus
	focused bool

	// mouse
	mouseX int
	mouseY int

	// scroll
	scrollY int

	// hot
	hotX, hotY int

	//zoom
	zoomed    bool
	zoomMode  bool
	zoomedPos fyne.Position

	// size updater
	szUpd  chan fyne.Size
	lkSize bool
}

func New(w fyne.Window, a fyne.App, c spice.Connector, password string) (*SpiceFyne, error) {
	final := &SpiceFyne{w: w, a: a, zoom: 1, szUpd: make(chan fyne.Size, 4)}
	final.ExtendBaseWidget(final)

	// create initial dummy image
	img := image.NewRGBA(image.Rect(0, 0, 800, 600))
	for p := range img.Pix {
		if p%4 == 3 {
			img.Pix[p] = 0xff
		} else {
			img.Pix[p] = 0
		}
	}
	final.output = canvas.NewImageFromImage(img)
	//final.output.ScaleMode = canvas.ImageScalePixels
	final.output.FillMode = canvas.ImageFillContain
	final.output.ScaleMode = canvas.ImageScaleFastest

	// configure visibility/etc
	final.output.Show()

	//final.cursor = canvas.NewImageFromImage(image.NewRGBA(image.Rect(0, 0, 16, 16)))
	//final.cursor.Resize(fyne.Size{64, 64})
	//final.cursor.Show()

	final.ctrl = newControl(final)

	spiceClient, err := spice.New(c, final, password)
	if err != nil {
		return nil, err
	}
	final.Client = spiceClient
	go final.szUpdThread()

	w.SetPadded(false)
	w.SetContent(container.NewMax(container.New(&autoLayout{direction: &final.ctrl.sidebarPosition, sidebarVisible: &final.ctrl.visible}, final, final.ctrl.main), final.ctrl.muteIconContainer))
	w.Canvas().Focus(final)

	if mon, err := goclip.NewMonitor(); err == nil {
		final.clip = mon
		mon.Subscribe(func(data goclip.Data) error {
			// grab clipboard now
			var types []spice.SpiceClipboardFormat
			typSkip := make(map[spice.SpiceClipboardFormat]bool)

			fmts, err := data.GetAllFormats()
			if err != nil {
				log.Printf("spicefyne: failed to fetch clipboard infos: %s", err)
				return nil
			}

			log.Printf("spicefyne: got clipboard %s", data)

			if data.Type() == goclip.Image {
				// always have png in priority if image
				types = append(types, spice.VD_AGENT_CLIPBOARD_IMAGE_PNG)
				typSkip[spice.VD_AGENT_CLIPBOARD_IMAGE_PNG] = true
			}

			for _, fmt := range fmts {
				fmt, ok := mime_go2sp(fmt.Mime())
				if !ok {
					continue
				}
				if _, found := typSkip[fmt]; found {
					continue
				}

				types = append(types, fmt)
				typSkip[fmt] = true
			}

			final.main.SendGrabClipboard(board_go2vd(data.Board()), types)
			return nil
		})
	}

	return final, nil
}

func (s *SpiceFyne) DisplayInit(img image.Image) {
	s.output.Image = img
	bounds := img.Bounds()
	s.width = bounds.Dx()
	s.height = bounds.Dy()

	if !s.init {
		if fullscreen := os.Getenv("SHELLS_FULLSCREEN"); fullscreen != "" {
			// special option: if SHELLS_FULLSCREEN=1 is passed in env, pass into fullscreen with passed size (eg 1920x1080)
			fs := strings.SplitN(fullscreen, "x", 2)
			if len(fs) == 2 {
				w, _ := strconv.ParseInt(fs[0], 10, 32)
				h, _ := strconv.ParseInt(fs[1], 10, 32)
				if w > 0 && h > 0 {
					// ok, good
					s.w.Resize(fyne.Size{Width: float32(w), Height: float32(h)})
					s.w.SetFullScreen(true)
					s.init = true

					s.UpdateSize(fyne.Size{Width: float32(w), Height: float32(h)})
				}
			}
		}
		if !s.init {
			s.w.Resize(fyne.Size{Width: float32(s.width), Height: float32(s.height)})
			s.w.CenterOnScreen()
			s.init = true

			// call updatesize with the window size in case it's different
			s.UpdateSize(s.w.Canvas().Size())
		}
	}
	canvas.Refresh(s.output)
}

func (s *SpiceFyne) DisplayRefresh() {
	canvas.Refresh(s.output)
}

func (s *SpiceFyne) SetEventsTarget(in *spice.ChInputs) {
	s.in = in
}

func (s *SpiceFyne) SetMainTarget(main *spice.ChMain) {
	s.main = main
}

func (s *SpiceFyne) CreateRenderer() fyne.WidgetRenderer {
	return &renderer{s}
}

// TODO only care for those events if MouseDown() isn't triggered (ie. mobile)
func (s *SpiceFyne) Tapped(t *fyne.PointEvent) {
}

func (s *SpiceFyne) TappedSecondary(*fyne.PointEvent) {
}

func (s *SpiceFyne) relPos(pos fyne.Position) (uint32, uint32) {
	if float32(s.width) == s.size.Width && float32(s.height) == s.size.Height {
		// easy
		s.mouseX, s.mouseY = int(pos.X), int(pos.Y)
		s.zoom = 1
		return uint32(pos.X), uint32(pos.Y)
	}
	x, y := pos.X, pos.Y

	// say I have a 100x100 image shown in a 50x200 canvas.Image, if I have a click at position 25x100 I should be getting 50x50
	ratioA := float64(s.width) / float64(s.height)
	ratioB := float64(s.size.Width) / float64(s.size.Height)

	if ratioA > ratioB {
		// top/bottom padded, width accurate
		zoom := float64(s.size.Width) / float64(s.width)
		realHeight := int(float64(s.height) * zoom)
		realTop := (s.size.Height - float32(realHeight)) / 2
		s.zoom = zoom
		x = float32(float64(x) / zoom)
		y = float32(float64(y-realTop) / zoom)
		//log.Printf("Ratio A bigger zoom=%f realHeight=%d realTop=%d", zoom, realHeight, realTop)
	} else {
		// left/right padded
		zoom := float64(s.size.Height) / float64(s.height)
		realWidth := int(float64(s.width) * zoom)
		realLeft := (s.size.Width - float32(realWidth)) / 2
		s.zoom = zoom
		x = float32(float64(x-realLeft) / zoom)
		y = float32(float64(y) / zoom)
		//log.Printf("Ratio B bigger zoom=%f realWidth=%d realLeft=(%d - %d)/2=%d", zoom, realWidth, s.size.Width, realWidth, realLeft)
	}

	if x < 0 || y < 0 || x > float32(s.width) || y > float32(s.height) {
		return uint32(s.mouseX), uint32(s.mouseY)
	}

	s.mouseX, s.mouseY = int(x), int(y)
	return uint32(x), uint32(y)
}

func (s *SpiceFyne) mouseBtn(v desktop.MouseButton) uint8 {
	// 0=left, 1=middle, 2=right, 3=up, 4=down, 5=side, 6=extra
	switch v {
	case desktop.MouseButtonPrimary:
		return 0
	case desktop.MouseButtonSecondary:
		return 2
	case desktop.MouseButtonTertiary:
		return 1
	default:
		return 0
	}
}

func (s *SpiceFyne) zoomIn(x uint32, y uint32) {
	size := s.Size()
	var posX float32
	var posY float32

	if float32(x)-size.Width/4 <= 0 {
		posX = 0
	} else if float32(x)+size.Width/4 >= size.Width*2 {
		posX = -size.Width
	} else {
		posX = float32(x) - size.Width/2
		if posX > 0 {
			posX = -(posX + size.Width/2)
		}
	}

	if float32(y)-size.Height/4 <= 0 {
		posY = 0
	} else if float32(y)+size.Height/4 >= size.Height*2 {
		posY = -size.Height
	} else {
		posY = float32(y) - size.Height/2
		if posY > 0 {
			posY = -(posY + size.Height/2)
		}
	}
	s.output.Resize(fyne.Size{Width: size.Width * 2, Height: size.Height * 2})
	s.output.Move(fyne.NewPos(posX, posY))
	s.zoomedPos = fyne.NewPos(posX, posY)
}

func (s *SpiceFyne) resetZoom() {
	size := s.Size()
	s.output.Resize(size)
	s.output.Move(fyne.NewPos(0, 0))
	s.zoomed = false
	s.ctrl.zoomBtn.Enable()
}

func (s *SpiceFyne) MouseDown(ev *desktop.MouseEvent) {
	if in := s.in; in != nil {
		x, y := s.relPos(ev.Position)
		if s.ctrl.visible {
			if x > SideBarWidth {
				s.ctrl.toggle()
			}
			return
		}
		if s.zoomMode {
			s.zoomIn(x, y)
			s.zoomMode = false
			s.zoomed = true
			s.ctrl.zoomBtn.Disable()
			return
		}

		if s.zoomed {
			x = x + uint32(s.zoomedPos.X)
			y = y + uint32(s.zoomedPos.Y)
		}

		in.MouseDown(s.mouseBtn(ev.Button), x, y)
	}
}

func (s *SpiceFyne) MouseUp(ev *desktop.MouseEvent) {
	if in := s.in; in != nil {
		x, y := s.relPos(ev.Position)
		in.MouseUp(s.mouseBtn(ev.Button), x, y)
	}
}

func (s *SpiceFyne) MouseMoved(ev *desktop.MouseEvent) {
	if in := s.in; in != nil {
		x, y := s.relPos(ev.Position)
		if s.ctrl.visible || s.zoomMode {
			return
		}

		if s.zoomed {
			x = (x + uint32(math.Abs(float64(s.zoomedPos.X)))) / 2
			y = (y + uint32(math.Abs(float64(s.zoomedPos.Y)))) / 2
			pos := ev.Position
			pos.X = float32(x)
			pos.Y = float32(y)
			s.updatedMouse(pos)

		} else {
			s.updatedMouse(ev.Position)
		}
		in.MousePosition(x, y)
	}
}

/*
func (s *SpiceFyne) Dragged(ev *fyne.DragEvent) {
	x, y := s.relPos(ev.Position)
	log.Printf("mouse dragged, %dx%d", x, y)
}
func (s *SpiceFyne) DragEnd() {
}
*/

func (s *SpiceFyne) MouseIn(*desktop.MouseEvent) {
}

func (s *SpiceFyne) MouseOut() {
}

func (s *SpiceFyne) Scrolled(ev *fyne.ScrollEvent) {
	// scroll happened
	log.Printf("scroll dx=%f dy=%f", ev.Scrolled.DX, ev.Scrolled.DY)

	s.scrollY += int(ev.Scrolled.DY)

	if -10 < s.scrollY && s.scrollY < 10 {
		return
	}

	if in := s.in; in != nil {
		x, y := s.relPos(ev.Position)
		if s.scrollY > 0 {
			// going up
			for s.scrollY > 0 {
				in.MouseDown(4, x, y)
				in.MouseUp(4, x, y)
				s.scrollY -= 10
			}
		} else if s.scrollY < 0 {
			// going down
			for s.scrollY < 0 {
				in.MouseDown(5, x, y)
				in.MouseUp(5, x, y)
				s.scrollY += 10
			}
		}
	}
}

func (s *SpiceFyne) KeyDown(ev *fyne.KeyEvent) {
	if in := s.in; in != nil {
		if k, ok := scancodeLookup(ev.Physical.ScanCode); ok && k != nil {
			in.OnKeyDown(k)
			return
		}
		k, ok := fyneScancodeMap[ev.Name]
		if !ok {
			log.Printf("unhandled key down = %s", ev.Name)
			return
		}
		in.OnKeyDown(k)
	}
}

func (s *SpiceFyne) KeyUp(ev *fyne.KeyEvent) {
	if in := s.in; in != nil {
		if k, ok := scancodeLookup(ev.Physical.ScanCode); ok && k != nil {
			in.OnKeyUp(k)
			return
		}
		k, ok := fyneScancodeMap[ev.Name]
		if !ok {
			log.Printf("unhandled key up = %s", ev.Name)
			return
		}
		in.OnKeyUp(k)
	}
}

// focus related
func (s *SpiceFyne) FocusGained() {
	if s.clip != nil {
		err := s.clip.Poll()
		if err != nil {
			log.Printf("spice/goclip: error during poll of clipboard information: %s", err)
		}
	}
	s.focused = true
}

func (s *SpiceFyne) FocusLost() {
	s.focused = false
}

func (s *SpiceFyne) Focused() bool {
	return s.focused
}

// TODO have a keyboard mode where we just emulate whatever is input by the user
func (s *SpiceFyne) TypedRune(rune) {
}

func (s *SpiceFyne) TypedKey(*fyne.KeyEvent) {
}

func (s *SpiceFyne) SetCursor(i image.Image, x, y uint16) {
	if s.cursor == nil {
		if i == nil {
			s.cursorObj = nil
		} else {
			s.cursorObj = &fyneCursor{i, int(x), int(y)}
		}
	} else if i == nil {
		s.cursor.Image = image.NewRGBA(image.Rect(0, 0, 16, 16))
	} else {
		s.cursor.Image = i
		s.hotX = int(x)
		s.hotY = int(y)
		bounds := i.Bounds()
		s.cursor.Resize(fyne.Size{float32(bounds.Dx()), float32(bounds.Dy())})
		canvas.Refresh(s.cursor)
	}
}

func (s *SpiceFyne) updatedMouse(p fyne.Position) {
	if s.cursor != nil {
		s.cursor.Move(p)
	}
}

func (s *SpiceFyne) Cursor() desktop.Cursor {
	if s.cursor == nil {
		if s.cursorObj == nil {
			return desktop.DefaultCursor
		} else {
			return s.cursorObj
		}
	} else {
		return desktop.HiddenCursor
	}
}

func (s *SpiceFyne) terminate() {
	dialog.ShowConfirm("Exit?", "Are you sure you want to quit?", func(f bool) {
		if f {
			os.Exit(0)
		}
	}, s.w)
}

// implement fyne.Tabbable
func (s *SpiceFyne) AcceptsTab() bool {
	return true
}
