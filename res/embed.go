package res

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed rightarrow.svg
var rightArrowBin []byte

//go:embed shellslogo.png
var shellsLogoBin []byte

//go:embed Icon-alpha-128x128.png
var shellsIcon []byte

var RightArrow = &fyne.StaticResource{
	StaticName:    "res/rightarrow.svg",
	StaticContent: rightArrowBin,
}

var ShellsLogo = &fyne.StaticResource{
	StaticName:    "res/shellslogo.png",
	StaticContent: shellsLogoBin,
}

var ShellsIcon = &fyne.StaticResource{
	StaticName:    "res/Icon-alpha-128x128.png",
	StaticContent: shellsIcon,
}
