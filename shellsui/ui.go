package shellsui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/Shells-com/shells-go/res"
)

var logoContainer *fyne.Container

func getLogoContainer() *fyne.Container {
	if logoContainer == nil {
		logo := canvas.NewImageFromResource(res.ShellsLogo)
		logo.SetMinSize(fyne.Size{Width: 242, Height: 76})

		cLogo := container.New(layout.NewCenterLayout(), logo)

		logoContainer = container.New(
			layout.NewBorderLayout(
				nil,
				GetMarginRectangleWithHeight(40),
				nil,
				nil,
			),
			cLogo,
		)
	}

	return logoContainer
}

func GetMainContainer(content *fyne.Container) *fyne.Container {
	return container.New(
		layout.NewBorderLayout(
			GetMarginRectangle(),
			GetMarginRectangle(),
			GetMarginRectangle(),
			GetMarginRectangle(),
		),
		container.NewVBox(
			getLogoContainer(),
			content,
		),
	)
}
