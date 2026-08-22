// Package avatarimg normalizes uploaded avatar images.
package avatarimg

import (
	"bytes"
	"errors"
	"fmt"
	"image"

	"github.com/disintegration/imaging"

	// Registers the webp decoder with image.Decode. Go has webp support
	// neither in the standard library nor, for encoding, anywhere at all,
	// which is why a webp upload is stored as PNG.
	_ "golang.org/x/image/webp"
)

const (
	// Size is the side of the stored avatar in pixels. It doubles the
	// largest place the avatar is rendered in, so that retina screens have
	// pixels to spare.
	Size = 512

	// MimeType is the content type every normalized avatar is stored with.
	MimeType = "image/png"

	// maxDimension caps the side of an accepted image. The upload is already
	// capped in bytes, but a megabyte of PNG unpacks into gigabytes of
	// pixels, so the header is inspected before anything is decoded.
	maxDimension = 8192
)

var (
	// ErrUnsupportedFormat reports data that is not an image in any format
	// the decoders registered here understand.
	ErrUnsupportedFormat = errors.New("unsupported image format")

	// ErrTooLarge reports an image whose declared dimensions are beyond
	// anything an avatar needs, and which is therefore not worth decoding.
	ErrTooLarge = errors.New("image dimensions are too large")
)

// Normalize turns an uploaded image into a stored avatar: it applies the EXIF
// orientation, crops the largest centered square and scales it down to Size,
// then encodes the result as PNG. An image smaller than Size keeps its size —
// upscaling invents no detail.
func Normalize(data []byte) ([]byte, error) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnsupportedFormat, err)
	}

	if cfg.Width > maxDimension || cfg.Height > maxDimension {
		return nil, fmt.Errorf("%w: %dx%d", ErrTooLarge, cfg.Width, cfg.Height)
	}

	img, err := imaging.Decode(bytes.NewReader(data), imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnsupportedFormat, err)
	}

	side := min(img.Bounds().Dx(), img.Bounds().Dy())
	img = imaging.CropCenter(img, side, side)

	if side > Size {
		img = imaging.Resize(img, Size, Size, imaging.Lanczos)
	}

	buf := new(bytes.Buffer)
	if err := imaging.Encode(buf, img, imaging.PNG); err != nil {
		return nil, fmt.Errorf("encode avatar: %w", err)
	}

	return buf.Bytes(), nil
}
