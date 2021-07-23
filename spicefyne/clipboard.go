package spicefyne

import (
	"bytes"
	"context"
	"errors"
	"image/png"
	"log"
	"time"

	"github.com/KarpelesLab/goclip"
	"github.com/Shells-com/spice"
)

type clipOption struct {
	spice         *SpiceFyne
	selection     spice.SpiceClipboardSelection
	clipboardType spice.SpiceClipboardFormat
}

func (o *clipOption) Type() goclip.Type {
	switch o.clipboardType {
	case spice.VD_AGENT_CLIPBOARD_UTF8_TEXT:
		return goclip.Text
	case spice.VD_AGENT_CLIPBOARD_IMAGE_PNG:
		return goclip.Image
	case spice.VD_AGENT_CLIPBOARD_IMAGE_BMP:
		return goclip.Image
	case spice.VD_AGENT_CLIPBOARD_IMAGE_TIFF:
		return goclip.Image
	case spice.VD_AGENT_CLIPBOARD_IMAGE_JPG:
		return goclip.Image
	default:
		return goclip.Invalid
	}
}

func (o *clipOption) Mime() string {
	switch o.clipboardType {
	case spice.VD_AGENT_CLIPBOARD_UTF8_TEXT:
		return "text/plain;charset=utf-8"
	case spice.VD_AGENT_CLIPBOARD_IMAGE_PNG:
		return "image/png"
	case spice.VD_AGENT_CLIPBOARD_IMAGE_BMP:
		return "image/bmp"
	case spice.VD_AGENT_CLIPBOARD_IMAGE_TIFF:
		return "image/tiff"
	case spice.VD_AGENT_CLIPBOARD_IMAGE_JPG:
		return "image/jpeg"
	default:
		return "application/octet-string"
	}
}

func (o *clipOption) Data(ctx context.Context) ([]byte, error) {
	// fetch from spice
	return o.spice.main.RequestClipboard(o.selection, o.clipboardType)
}

func (s *SpiceFyne) ClipboardGrabbed(selection spice.SpiceClipboardSelection, clipboardTypes []spice.SpiceClipboardFormat) {
	//s.main.RequestClipboard(selection, clipboardTypes[0])
	var opts []goclip.DataOption

	for _, opt := range clipboardTypes {
		opts = append(opts, &clipOption{selection: selection, clipboardType: opt, spice: s})
	}

	data := &goclip.StaticData{TargetBoard: board_vd2go(selection), Options: opts}
	err := goclip.CopyTo(context.Background(), board_vd2go(selection), data)
	if err != nil {
		log.Printf("copy operation failed: %s", err)
	}
}

func (s *SpiceFyne) ClipboardFetch(selection spice.SpiceClipboardSelection, clipboardType spice.SpiceClipboardFormat) ([]byte, error) {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	data, err := goclip.PasteFrom(ctx, board_vd2go(selection))
	if err != nil {
		log.Printf("failed to paste: %s", err)
		return nil, err
	}

	//log.Printf("grabbed data from goclip: %+v", data)

	// typically we do not want to offer formats such as bmp or tiff as these are heavy

	switch clipboardType {
	case spice.VD_AGENT_CLIPBOARD_UTF8_TEXT:
		res, err := data.ToText(ctx)
		if err != nil {
			return nil, err
		}
		return []byte(res), nil
	case spice.VD_AGENT_CLIPBOARD_IMAGE_PNG:
		res, err := data.GetFormat(ctx, "image/png")
		if err == nil {
			return res, nil
		}
		log.Printf("spicefyne: failed to fetch png, will try to convert: %s", err)

		// fetch image, convert
		img, err := data.ToImage(ctx)
		if err != nil {
			return nil, err
		}
		buf := &bytes.Buffer{}
		err = png.Encode(buf, img)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	case spice.VD_AGENT_CLIPBOARD_IMAGE_BMP:
		return data.GetFormat(ctx, "image/bmp")
	case spice.VD_AGENT_CLIPBOARD_IMAGE_TIFF:
		return data.GetFormat(ctx, "image/tiff")
	case spice.VD_AGENT_CLIPBOARD_IMAGE_JPG:
		return data.GetFormat(ctx, "image/jpeg")
	default:
		return nil, errors.New("unsupported format")
	}
}

func (s *SpiceFyne) ClipboardRelease(selection spice.SpiceClipboardSelection) {
	log.Printf("spicefyne: TODO release clipboard")
}
