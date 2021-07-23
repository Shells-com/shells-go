package shellsui

import (
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

const LabelSize float32 = 16

func GetShellsLabel(text string) *canvas.Text {
	label := canvas.NewText(text, theme.ForegroundColor())
	label.TextSize = LabelSize
	return label
}
