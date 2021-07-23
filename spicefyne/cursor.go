package spicefyne

import "image"

type fyneCursor struct {
	img  image.Image
	x, y int
}

func (f *fyneCursor) Image() (image.Image, int, int) {
	return f.img, f.x, f.y
}
