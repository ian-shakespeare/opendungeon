package media

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/disintegration/imaging"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gen2brain/heic"
	"golang.org/x/image/webp"
)

var (
	ErrUnknownContentType     = errors.New("unknown content type")
	ErrUnsupportedImageFormat = errors.New("unsupported image format")
	ErrConvertFailed          = errors.New("convert failed")
)

func ConvertToAvatar(r io.Reader) (io.Reader, error) {
	b := new(bytes.Buffer)
	tr := io.TeeReader(r, b)

	mt, err := mimetype.DetectReader(tr)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnknownContentType, err.Error())
	}

	mr := io.MultiReader(b, r)

	var img image.Image
	switch {
	case mt.Is("image/jpeg"):
		img, err = jpeg.Decode(mr)
	case mt.Is("image/png"):
		img, err = png.Decode(mr)
	case mt.Is("image/heic"), mt.Is("image/heif"):
		img, err = heic.Decode(mr)
	case mt.Is("image/webp"):
		img, err = webp.Decode(mr)
	default:
		return nil, ErrUnsupportedImageFormat
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedImageFormat, err.Error())
	}

	resized := imaging.Resize(img, 128, 128, imaging.Linear)
	blurred := imaging.Blur(resized, 0.75)

	b.Reset()
	if err := png.Encode(b, blurred); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConvertFailed, err.Error())
	}

	return b, nil
}
