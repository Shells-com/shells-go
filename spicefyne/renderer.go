package spicefyne

import (
	"image/color"

	"fyne.io/fyne/v2"
)

type renderer struct {
	s *SpiceFyne
}

func (r *renderer) Layout(s fyne.Size) {
	r.s.size = s
	r.s.UpdateSize(s)
	r.s.output.Move(fyne.Position{0, 0})
	r.s.output.Resize(s)
}

func (r *renderer) MinSize() fyne.Size {
	return fyne.Size{320, 200}
}

func (r *renderer) Refresh() {
}

func (r *renderer) BackgroundColor() color.Color {
	return color.RGBA{0, 0, 0, 0xff}
}

func (r *renderer) Objects() []fyne.CanvasObject {
	if curs := r.s.cursor; curs != nil {
		return []fyne.CanvasObject{r.s.output, curs}
	} else {
		return []fyne.CanvasObject{r.s.output}
	}
}

func (r *renderer) Destroy() {
}
