package shellsui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

const WMargin float32 = 10

func GetMarginRectangle() *canvas.Rectangle {
	margin := canvas.NewRectangle(color.Transparent)
	margin.SetMinSize(fyne.NewSize(WMargin, WMargin))
	return margin
}

func GetMarginRectangleWithWidth(width int) *canvas.Rectangle {
	if width < 0 {
		width = int(WMargin)
	}

	margin := canvas.NewRectangle(color.Transparent)
	margin.SetMinSize(fyne.NewSize(float32(width), 1))
	return margin
}

func GetMarginRectangleWithHeight(height int) *canvas.Rectangle {
	if height < 0 {
		height = int(WMargin)
	}

	margin := canvas.NewRectangle(color.Transparent)
	margin.SetMinSize(fyne.NewSize(1, float32(height)))
	return margin
}
