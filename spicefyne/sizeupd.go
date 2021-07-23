package spicefyne

import (
	"log"
	"time"

	"fyne.io/fyne/v2"
)

func (s *SpiceFyne) szUpdThread() {
	t := time.NewTicker(40 * time.Millisecond)
	run := false
	cnt := 0
	t.Stop()
	var siz fyne.Size

	for {
		select {
		case siz = <-s.szUpd:
			//log.Printf("scheduling update size")
			cnt = 0
			if !run {
				run = true
				t.Reset(40 * time.Millisecond)
			}
		case <-t.C:
			if s.lkSize {
				t.Stop()
				run = false
				break
			}
			//log.Printf("running update size")
			if cnt < 10 {
				cnt += 1
				break
			}
			cnt = 0
			t.Stop()
			run = false

			w, h := siz.Width, siz.Height
			if w < 1280 {
				w = 1280
			}
			if h < 720 {
				h = 720
			}
			if w > siz.Width {
				// we're going to have to stretch height
				zoom := float64(siz.Width) / float64(w) // zoom<1
				h = float32(float64(h) / zoom)
				log.Printf("resize: zoom=%f size=%fx%f → %fx%f", zoom, siz.Width, siz.Height, w, h)
			} else if h > siz.Height {
				// we're going to have to stretch width
				zoom := float64(siz.Height) / float64(h) // zoom<1
				w = float32(float64(w) / zoom)
				log.Printf("resize: zoom=%f size=%fx%f → %fx%f", zoom, siz.Width, siz.Height, w, h)
			}
			if w == float32(s.width) && h == float32(s.height) {
				break
			}
			s.Client.UpdateView(int(w), int(h))
		}
	}
}

func (s *SpiceFyne) UpdateSize(siz fyne.Size) {
	if !s.init {
		return
	}

	//log.Printf("update size queued")
	s.szUpd <- siz
}
