// Package avatarimg normalizes uploaded avatar images.
package avatarimg

import (
	"bytes"
	"errors"
	"fmt"
	"image"

	"github.com/disintegration/imaging"
	"github.com/rs/zerolog"

	// Registers the webp decoder with image.Decode. Go has webp support
	// neither in the standard library nor, for encoding, anywhere at all,
	// which is why a webp upload is stored as PNG.
	_ "golang.org/x/image/webp"
)

const (
	// MaxSide is the longest side a stored avatar can have, in pixels. It
	// doubles the largest place the avatar is rendered in, so that retina
	// screens have pixels to spare.
	MaxSide = 512

	// MimeType is the content type every normalized avatar is stored with.
	MimeType = "image/png"

	// maxPixels caps the area of an accepted image, because area is what
	// decoding costs and the weight of the file says nothing about it: 105kb
	// of flat PNG peaks at 195mb of heap, the picture being held decoded and
	// cropped at once. Hence the header is read before anything is decoded.
	// Twenty five megapixels is around as much as a real photo fits into the
	// upload limit anyway, and doubling it doubles the worst case.
	maxPixels = 25_000_000

	// heavyPixels marks an upload worth a line in the log. An ordinary phone
	// photo is twelve megapixels; past that it is either a proper camera or
	// somebody feeling out what the server will swallow.
	heavyPixels = 12_000_000
)

var (
	// ErrUnsupportedFormat reports data that is not an image in any format
	// the decoders registered here understand.
	ErrUnsupportedFormat = errors.New("unsupported image format")

	// ErrTooLarge reports an image with more pixels than an avatar could ever
	// need, and which is therefore not worth decoding.
	ErrTooLarge = errors.New("image dimensions are too large")
)

// Normalize turns an uploaded image into a stored avatar: it applies the EXIF
// orientation, crops the largest centered square and scales it down to
// MaxSide, then encodes the result as PNG. A picture smaller than that keeps
// its size — upscaling invents no detail. An upload heavier than an ordinary
// photo is noted in the log before it is decoded.
func Normalize(logger zerolog.Logger, data []byte) ([]byte, error) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnsupportedFormat, err)
	}

	pixels := int64(cfg.Width) * int64(cfg.Height)
	if pixels > maxPixels {
		return nil, fmt.Errorf("%w: %dx%d", ErrTooLarge, cfg.Width, cfg.Height)
	}

	if pixels > heavyPixels {
		logger.Warn().Msgf("heavy avatar upload: %dx%d, %d kb",
			cfg.Width, cfg.Height, len(data)/1024)
	}

	img, err := imaging.Decode(bytes.NewReader(data), imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnsupportedFormat, err)
	}

	side := min(img.Bounds().Dx(), img.Bounds().Dy())
	img = imaging.CropCenter(img, side, side)

	if side > MaxSide {
		img = imaging.Resize(img, MaxSide, MaxSide, imaging.Lanczos)
	}

	buf := new(bytes.Buffer)
	if err := imaging.Encode(buf, img, imaging.PNG); err != nil {
		return nil, fmt.Errorf("encode avatar: %w", err)
	}

	return buf.Bytes(), nil
}
