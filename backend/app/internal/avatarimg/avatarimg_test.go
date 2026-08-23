package avatarimg

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// solidPNG makes a w×h image filled with c and encodes it as PNG.
func solidPNG(t *testing.T, w, h int, c color.Color) []byte {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.Set(x, y, c)
		}
	}

	buf := new(bytes.Buffer)
	require.NoError(t, png.Encode(buf, img))
	return buf.Bytes()
}

// solidJPEG makes a w×h image filled with c and encodes it as JPEG.
func solidJPEG(t *testing.T, w, h int, c color.Color) []byte {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.Set(x, y, c)
		}
	}

	buf := new(bytes.Buffer)
	require.NoError(t, jpeg.Encode(buf, img, nil))
	return buf.Bytes()
}

// halvesJPEG makes a w×h image split into a top and a bottom half and encodes
// it as JPEG.
func halvesJPEG(t *testing.T, w, h int, top, bottom color.Color) []byte {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			if y < h/2 {
				img.Set(x, y, top)
			} else {
				img.Set(x, y, bottom)
			}
		}
	}

	buf := new(bytes.Buffer)
	require.NoError(t, jpeg.Encode(buf, img, nil))
	return buf.Bytes()
}

// withEXIFOrientation injects an APP1 segment holding a single Orientation tag
// into an encoded JPEG. Go's jpeg encoder writes no EXIF at all, so a photo
// straight from a phone can only be imitated by splicing the segment in.
func withEXIFOrientation(t *testing.T, data []byte, orientation byte) []byte {
	t.Helper()
	require.Equal(t, []byte{0xFF, 0xD8}, data[:2], "not a JPEG")

	exif := []byte{
		'E', 'x', 'i', 'f', 0, 0,
		'I', 'I', 42, 0, // TIFF header, little-endian
		8, 0, 0, 0, // IFD0 follows the header
		1, 0, // a single entry
		0x12, 0x01, // tag 0x0112, Orientation
		3, 0, // type SHORT
		1, 0, 0, 0, // one value
		orientation, 0, 0, 0,
		0, 0, 0, 0, // no next IFD
	}

	out := new(bytes.Buffer)
	out.Write(data[:2])                                   // SOI
	out.Write([]byte{0xFF, 0xE1, 0, byte(len(exif) + 2)}) // APP1 and its length
	out.Write(exif)
	out.Write(data[2:])
	return out.Bytes()
}

// pngHeader builds a PNG carrying nothing but an IHDR chunk declaring w×h, so
// that the picture is never actually allocated: a decoder that reads the
// header first has everything it needs to refuse.
func pngHeader(w, h uint32) []byte {
	ihdr := make([]byte, 4+13)
	copy(ihdr, "IHDR")
	binary.BigEndian.PutUint32(ihdr[4:], w)
	binary.BigEndian.PutUint32(ihdr[8:], h)
	ihdr[12] = 8 // bit depth
	ihdr[13] = 2 // color type: truecolor

	out := new(bytes.Buffer)
	out.Write([]byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A})
	binary.Write(out, binary.BigEndian, uint32(len(ihdr)-4))
	out.Write(ihdr)
	binary.Write(out, binary.BigEndian, crc32.ChecksumIEEE(ihdr))
	return out.Bytes()
}

// stripes makes a w×h image split into three vertical stripes of equal width,
// so that a centered square crop keeps the middle one only.
func stripes(w, h int, left, middle, right color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			switch {
			case x < w/3:
				img.Set(x, y, left)
			case x < 2*w/3:
				img.Set(x, y, middle)
			default:
				img.Set(x, y, right)
			}
		}
	}

	return img
}

// stripesPNG encodes stripes as PNG.
func stripesPNG(t *testing.T, w, h int, left, middle, right color.Color) []byte {
	t.Helper()

	buf := new(bytes.Buffer)
	require.NoError(t, png.Encode(buf, stripes(w, h, left, middle, right)))
	return buf.Bytes()
}

// stripesJPEG encodes stripes as JPEG, which decodes into subsampled YCbCr
// planes rather than into the RGBA a PNG gives back.
func stripesJPEG(t *testing.T, w, h int, left, middle, right color.Color) []byte {
	t.Helper()

	buf := new(bytes.Buffer)
	require.NoError(t, jpeg.Encode(buf, stripes(w, h, left, middle, right), nil))
	return buf.Bytes()
}

// assertColor compares a pixel to the expected color, tolerating the rounding
// a resampling filter introduces.
func assertColor(t *testing.T, img image.Image, x, y int, expected color.Color) {
	t.Helper()

	require.True(t, image.Pt(x, y).In(img.Bounds()),
		"point (%d, %d) is outside of %v", x, y, img.Bounds())

	er, eg, eb, _ := expected.RGBA()
	ar, ag, ab, _ := img.At(x, y).RGBA()

	const tolerance = 2 << 8
	assert.InDelta(t, er, ar, tolerance, "red at (%d, %d)", x, y)
	assert.InDelta(t, eg, ag, tolerance, "green at (%d, %d)", x, y)
	assert.InDelta(t, eb, ab, tolerance, "blue at (%d, %d)", x, y)
}

// decodePNG decodes normalized output, asserting it is really a PNG.
func decodePNG(t *testing.T, data []byte) image.Image {
	t.Helper()

	img, format, err := image.Decode(bytes.NewReader(data))
	require.NoError(t, err)
	require.Equal(t, "png", format)
	return img
}

// opaqueImage hides the SubImage method of the image it wraps, standing in for
// a decoder that hands out a type the crop cannot share pixels with.
type opaqueImage struct{ image.Image }

func TestNormalize(t *testing.T) {
	t.Run("should downscale a large image to MaxSide", func(t *testing.T) {
		out, err := Normalize(zerolog.Nop(), solidPNG(t, 1000, 1000, color.RGBA{R: 10, G: 20, B: 30, A: 255}))
		require.NoError(t, err)

		img := decodePNG(t, out)
		assert.Equal(t, MaxSide, img.Bounds().Dx())
		assert.Equal(t, MaxSide, img.Bounds().Dy())
	})

	t.Run("should crop a wide image to its center instead of squashing it", func(t *testing.T) {
		red := color.RGBA{R: 255, A: 255}
		green := color.RGBA{G: 255, A: 255}
		blue := color.RGBA{B: 255, A: 255}

		out, err := Normalize(zerolog.Nop(), stripesPNG(t, 1800, 600, red, green, blue))
		require.NoError(t, err)

		img := decodePNG(t, out)
		assertColor(t, img, 5, MaxSide/2, green)
		assertColor(t, img, MaxSide-5, MaxSide/2, green)
	})

	t.Run("should crop a wide jpeg to its center as well", func(t *testing.T) {
		red := color.RGBA{R: 255, A: 255}
		green := color.RGBA{G: 255, A: 255}
		blue := color.RGBA{B: 255, A: 255}

		// An odd width leaves the crop at an odd offset, where the chroma
		// planes of a subsampled JPEG are indexed a pixel off the luma one.
		out, err := Normalize(zerolog.Nop(), stripesJPEG(t, 1802, 600, red, green, blue))
		require.NoError(t, err)

		img := decodePNG(t, out)
		assertColor(t, img, MaxSide/4, MaxSide/2, green)
		assertColor(t, img, 3*MaxSide/4, MaxSide/2, green)
	})

	t.Run("should re-encode a jpeg as png", func(t *testing.T) {
		out, err := Normalize(zerolog.Nop(), solidJPEG(t, 800, 800, color.RGBA{R: 10, G: 20, B: 30, A: 255}))
		require.NoError(t, err)

		// decodePNG fails unless the output really is a PNG.
		img := decodePNG(t, out)
		assert.Equal(t, MaxSide, img.Bounds().Dx())
	})

	t.Run("should accept webp and store it as png", func(t *testing.T) {
		data, err := os.ReadFile(filepath.Join("testdata", "stripes.webp"))
		require.NoError(t, err)

		out, err := Normalize(zerolog.Nop(), data)
		require.NoError(t, err)

		img := decodePNG(t, out)
		assert.Equal(t, MaxSide, img.Bounds().Dx())
		assertColor(t, img, MaxSide/2, MaxSide/2, color.RGBA{G: 255, A: 255})
	})

	t.Run("should rotate a jpeg according to its EXIF orientation", func(t *testing.T) {
		red := color.RGBA{R: 255, A: 255}
		blue := color.RGBA{B: 255, A: 255}

		// Orientation 6 means "rotate 90° clockwise to display", which moves
		// the top half of the source to the right half of the picture.
		data := withEXIFOrientation(t, halvesJPEG(t, 1200, 800, red, blue), 6)

		out, err := Normalize(zerolog.Nop(), data)
		require.NoError(t, err)

		img := decodePNG(t, out)
		assertColor(t, img, 20, 20, blue)
		assertColor(t, img, MaxSide-20, 20, red)
	})

	t.Run("should log an upload heavier than an ordinary photo", func(t *testing.T) {
		log := new(bytes.Buffer)

		// Thirteen megapixels: past what a phone produces, inside the cap.
		_, err := Normalize(zerolog.New(log),
			solidPNG(t, 3600, 3600, color.RGBA{R: 10, G: 20, B: 30, A: 255}))
		require.NoError(t, err)

		assert.Contains(t, log.String(), "3600x3600")
	})

	t.Run("should say nothing about an ordinary upload", func(t *testing.T) {
		log := new(bytes.Buffer)

		_, err := Normalize(zerolog.New(log),
			solidPNG(t, 1000, 1000, color.RGBA{R: 10, G: 20, B: 30, A: 255}))
		require.NoError(t, err)

		assert.Empty(t, log.String())
	})

	t.Run("should reject an image with too many pixels", func(t *testing.T) {
		// Neither side is outlandish on its own; the area is.
		_, err := Normalize(zerolog.Nop(), pngHeader(8000, 8000))
		require.ErrorIs(t, err, ErrTooLarge)
	})

	t.Run("should reject data that is not an image", func(t *testing.T) {
		_, err := Normalize(zerolog.Nop(), []byte("definitely not an image"))
		require.ErrorIs(t, err, ErrUnsupportedFormat)
	})

	t.Run("should not upscale an image smaller than MaxSide", func(t *testing.T) {
		out, err := Normalize(zerolog.Nop(), solidPNG(t, 100, 100, color.RGBA{R: 10, G: 20, B: 30, A: 255}))
		require.NoError(t, err)

		img := decodePNG(t, out)
		assert.Equal(t, 100, img.Bounds().Dx())
		assert.Equal(t, 100, img.Bounds().Dy())
	})
}

func TestCropSquare(t *testing.T) {
	t.Run("should take the centered square of a wide image", func(t *testing.T) {
		out := cropSquare(image.NewNRGBA(image.Rect(0, 0, 100, 40)))

		assert.Equal(t, 40, out.Bounds().Dx())
		assert.Equal(t, 40, out.Bounds().Dy())
		assert.Equal(t, image.Pt(30, 0), out.Bounds().Min)
	})

	t.Run("should take the centered square of a tall image", func(t *testing.T) {
		out := cropSquare(image.NewNRGBA(image.Rect(0, 0, 40, 100)))

		assert.Equal(t, 40, out.Bounds().Dx())
		assert.Equal(t, 40, out.Bounds().Dy())
		assert.Equal(t, image.Pt(0, 30), out.Bounds().Min)
	})

	t.Run("should share the pixels of the source instead of copying them", func(t *testing.T) {
		src := image.NewNRGBA(image.Rect(0, 0, 100, 40))
		out := cropSquare(src)

		red := color.RGBA{R: 255, A: 255}
		src.Set(50, 20, red)

		assertColor(t, out, 50, 20, red)
	})

	t.Run("should copy an image that cannot give out a sub-image", func(t *testing.T) {
		red := color.RGBA{R: 255, A: 255}
		src := image.NewNRGBA(image.Rect(0, 0, 100, 40))
		src.Set(50, 20, red)

		out := cropSquare(opaqueImage{src})
		assert.Equal(t, 40, out.Bounds().Dx())
		assertColor(t, out, 20, 20, red)
	})
}
