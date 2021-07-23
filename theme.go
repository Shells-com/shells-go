package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type shellsTheme struct {
	background                                                                    color.Color
	button, disabledButton, text, placeholder, hover, shadow, disabled, scrollBar color.Color
	regular, bold, italic, boldItalic, monospace                                  fyne.Resource
}

func getShellsTheme() fyne.Theme {
	r := &shellsTheme{
		background:     color.NRGBA{0xff, 0xff, 0xff, 0xff},
		button:         color.Transparent,
		disabled:       color.NRGBA{0x0, 0x0, 0x0, 0x42},
		disabledButton: color.NRGBA{0xe5, 0xe5, 0xe5, 0xff},
		text:           color.NRGBA{0x21, 0x21, 0x21, 0xff},
		placeholder:    color.NRGBA{0x88, 0x88, 0x88, 0xff},
		hover:          color.NRGBA{0x0, 0x0, 0x0, 0x0f},
		scrollBar:      color.NRGBA{0x0, 0x0, 0x0, 0x99},
		shadow:         color.NRGBA{0x0, 0x0, 0x0, 0x33},

		regular:    theme.DefaultTextFont(),
		bold:       theme.DefaultTextBoldFont(),
		italic:     theme.DefaultTextItalicFont(),
		boldItalic: theme.DefaultTextBoldItalicFont(),
		monospace:  theme.DefaultTextMonospaceFont(),
	}
	return theme.FromLegacy(r) //we will have to implement the new version of theme
}

func (c shellsTheme) BackgroundColor() color.Color {
	return c.background
}

func (c shellsTheme) ButtonColor() color.Color {
	return c.button
}

func (c shellsTheme) DisabledButtonColor() color.Color {
	return c.disabledButton
}

func (c shellsTheme) HyperlinkColor() color.Color {
	// Deprecated: Hyperlinks now use the primary color for consistency.
	return c.PrimaryColor()
}

func (c shellsTheme) TextColor() color.Color {
	return c.text
}

func (c shellsTheme) DisabledTextColor() color.Color {
	return c.disabled
}

func (c shellsTheme) IconColor() color.Color {
	// Deprecated: Icons now use the text colour for consistency.
	return c.TextColor()
}

func (c shellsTheme) DisabledIconColor() color.Color {
	// Deprecated: Disabled icons match disabled text color for consistency.
	return c.DisabledTextColor()
}

func (c shellsTheme) PlaceHolderColor() color.Color {
	return c.placeholder
}

func (c shellsTheme) PrimaryColor() color.Color {
	return color.NRGBA{R: 0x52, G: 0xbd, B: 0x15, A: 0xff}
}

func (c shellsTheme) HoverColor() color.Color {
	return c.hover
}

func (c shellsTheme) FocusColor() color.Color {
	return c.PrimaryColor()
}

func (c shellsTheme) ScrollBarColor() color.Color {
	return c.scrollBar
}

func (c shellsTheme) ShadowColor() color.Color {
	return c.shadow
}

func (c shellsTheme) TextSize() int {
	return 14
}

func (c shellsTheme) TextFont() fyne.Resource {
	return c.regular
}

func (c shellsTheme) TextBoldFont() fyne.Resource {
	return c.bold
}

func (c shellsTheme) TextItalicFont() fyne.Resource {
	return c.italic
}

func (c shellsTheme) TextBoldItalicFont() fyne.Resource {
	return c.boldItalic
}

func (c shellsTheme) TextMonospaceFont() fyne.Resource {
	return c.monospace
}

func (c shellsTheme) Padding() int {
	return 4
}

func (c shellsTheme) IconInlineSize() int {
	return 20
}

func (c shellsTheme) ScrollBarSize() int {
	return 16
}

func (c shellsTheme) ScrollBarSmallSize() int {
	return 3
}
