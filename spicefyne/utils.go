package spicefyne

import (
	"fmt"

	"github.com/KarpelesLab/goclip"
	"github.com/Shells-com/spice"
)

func GetVDAgentClipboardTypeByFormatType(ft goclip.Type) (spice.SpiceClipboardFormat, error) {
	switch ft {
	case goclip.Text:
		return spice.VD_AGENT_CLIPBOARD_UTF8_TEXT, nil
	case goclip.Image:
		return spice.VD_AGENT_CLIPBOARD_IMAGE_PNG, nil
	}
	return spice.VD_AGENT_CLIPBOARD_NONE, fmt.Errorf("spice/clipboard: format type %d is not available", ft)
}

func board_go2vd(b goclip.Board) spice.SpiceClipboardSelection {
	switch b {
	case goclip.Default:
		return spice.VD_AGENT_CLIPBOARD_SELECTION_CLIPBOARD
	case goclip.PrimarySelection:
		return spice.VD_AGENT_CLIPBOARD_SELECTION_PRIMARY
	case goclip.SecondarySelection:
		return spice.VD_AGENT_CLIPBOARD_SELECTION_SECONDARY
	default:
		return spice.VD_AGENT_CLIPBOARD_SELECTION_CLIPBOARD
	}
}

func board_vd2go(b spice.SpiceClipboardSelection) goclip.Board {
	switch b {
	case spice.VD_AGENT_CLIPBOARD_SELECTION_CLIPBOARD:
		return goclip.Default
	case spice.VD_AGENT_CLIPBOARD_SELECTION_PRIMARY:
		return goclip.PrimarySelection
	case spice.VD_AGENT_CLIPBOARD_SELECTION_SECONDARY:
		return goclip.SecondarySelection
	default:
		return goclip.InvalidBoard
	}
}

func mime_go2sp(mime string) (spice.SpiceClipboardFormat, bool) {
	switch mime {
	case "text/plain", "text/plain;charset=utf-8":
		return spice.VD_AGENT_CLIPBOARD_UTF8_TEXT, true
	case "image/png":
		return spice.VD_AGENT_CLIPBOARD_IMAGE_PNG, true
	case "image/bmp":
		return spice.VD_AGENT_CLIPBOARD_IMAGE_BMP, true
	case "image/tiff":
		return spice.VD_AGENT_CLIPBOARD_IMAGE_TIFF, true
	case "image/jpeg":
		return spice.VD_AGENT_CLIPBOARD_IMAGE_JPG, true
	default:
		return spice.VD_AGENT_CLIPBOARD_NONE, false
	}
}
