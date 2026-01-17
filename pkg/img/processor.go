package img

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/bmp"
	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/tiff"
)

// Processor wraps an image.Image to provide chainable processing methods.
type Processor struct {
	Img          image.Image
	Quality      int    // Quality/Compression level (1-100). Default is 90.
	TargetFormat string // Explicit output format (e.g., "png", "jpeg").
}

// NewProcessor creates a new Processor from an image.
func NewProcessor(img image.Image) *Processor {
	return &Processor{
		Img:          img,
		Quality:      90,
		TargetFormat: "", // Default to inferred from extension
	}
}

// Resize resizes the image to the specified width and height using Catmull-Rom interpolation.
func (p *Processor) Resize(width, height int) *Processor {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	// CatmullRom is a high-quality resampling filter.
	xdraw.CatmullRom.Scale(dst, dst.Rect, p.Img, p.Img.Bounds(), xdraw.Over, nil)
	p.Img = dst
	return p
}

// Grayscale converts the image to grayscale.
func (p *Processor) Grayscale() *Processor {
	bounds := p.Img.Bounds()
	dst := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dst.Set(x, y, p.Img.At(x, y))
		}
	}
	p.Img = dst
	return p
}

// Watermark adds a text watermark to the bottom-right corner of the image.
func (p *Processor) Watermark(text string) *Processor {
	bounds := p.Img.Bounds()
	// Create a mutable RGBA image from the current image
	m := image.NewRGBA(bounds)
	draw.Draw(m, bounds, p.Img, bounds.Min, draw.Src)

	// Add text
	col := color.RGBA{255, 255, 255, 128} // Semi-transparent white
	point := fixed.Point26_6{
		X: fixed.I(bounds.Max.X - len(text)*7 - 10), // Rough estimation of text width
		Y: fixed.I(bounds.Max.Y - 10),
	}

	d := &font.Drawer{
		Dst:  m,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(text)

	p.Img = m
	return p
}

// SetQuality sets the quality or compression level for the output (1-100).
func (p *Processor) SetQuality(quality int) *Processor {
	if quality < 1 {
		quality = 1
	}
	if quality > 100 {
		quality = 100
	}
	p.Quality = quality
	return p
}

// Rotate180 rotates the image 180 degrees.
func (p *Processor) Rotate180() *Processor {
	bounds := p.Img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// New X = original Width - 1 - original X
			// New Y = original Height - 1 - original Y
			dst.Set(bounds.Max.X-1-x, bounds.Max.Y-1-y, p.Img.At(x, y))
		}
	}

	p.Img = dst
	return p
}

// Rotate90 rotates the image 90 degrees clockwise.
func (p *Processor) Rotate90() *Processor {
	bounds := p.Img.Bounds()
	// New dimensions: width becomes height, height becomes width
	newBounds := image.Rect(0, 0, bounds.Dy(), bounds.Dx())
	dst := image.NewRGBA(newBounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// New X = original Height - 1 - original Y
			// New Y = original X
			dst.Set(bounds.Max.Y-1-y, x, p.Img.At(x, y))
		}
	}

	p.Img = dst
	return p
}

// FlipHorizontal flips the image horizontally.
func (p *Processor) FlipHorizontal() *Processor {
	bounds := p.Img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// New X = original Width - 1 - original X
			dst.Set(bounds.Max.X-1-x, y, p.Img.At(x, y))
		}
	}

	p.Img = dst
	return p
}

// FlipVertical flips the image vertically.
func (p *Processor) FlipVertical() *Processor {
	bounds := p.Img.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// New Y = original Height - 1 - original Y
			dst.Set(x, bounds.Max.Y-1-y, p.Img.At(x, y))
		}
	}

	p.Img = dst
	return p
}

// Brightness adjusts the brightness of the image by a percentage (-100 to 100).
func (p *Processor) Brightness(percentage int) *Processor {
	bounds := p.Img.Bounds()
	dst := image.NewRGBA(bounds)

	// Clamp percentage
	if percentage < -100 {
		percentage = -100
	}
	if percentage > 100 {
		percentage = 100
	}

	scale := float64(percentage) / 100.0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := p.Img.At(x, y)
			r, g, b, a := originalColor.RGBA()

			// RGBA() returns 0-65535, we need to convert back to 0-255 for standard processing
			nr := uint8(r >> 8)
			ng := uint8(g >> 8)
			nb := uint8(b >> 8)
			na := uint8(a >> 8)

			// Apply brightness adjustment
			adjust := func(val uint8) uint8 {
				newVal := float64(val) + (255.0 * scale)
				if newVal > 255 {
					return 255
				}
				if newVal < 0 {
					return 0
				}
				return uint8(newVal)
			}

			dst.Set(x, y, color.RGBA{adjust(nr), adjust(ng), adjust(nb), na})
		}
	}

	p.Img = dst
	return p
}

// --- Format Conversion & Output Methods ---

// ToPNG sets the target output format to PNG.
func (p *Processor) ToPNG() *Processor {
	p.TargetFormat = "png"
	return p
}

// ToJPEG sets the target output format to JPEG.
func (p *Processor) ToJPEG() *Processor {
	p.TargetFormat = "jpeg"
	return p
}

// ToGIF sets the target output format to GIF.
func (p *Processor) ToGIF() *Processor {
	p.TargetFormat = "gif"
	return p
}

// ToBMP sets the target output format to BMP.
func (p *Processor) ToBMP() *Processor {
	p.TargetFormat = "bmp"
	return p
}

// ToTIFF sets the target output format to TIFF.
func (p *Processor) ToTIFF() *Processor {
	p.TargetFormat = "tiff"
	return p
}

// Save saves the image to the specified path.
// If a target format was set via To<Format>(), it ignores the file extension.
// Otherwise, it infers the format from the file extension.
func (p *Processor) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	format := p.TargetFormat
	if format == "" {
		format = filepath.Ext(path)
	}
	return p.Encode(file, format)
}

// ToBytes returns the image as a byte slice.
// If format is empty, it uses the TargetFormat set on the processor.
func (p *Processor) ToBytes(format string) ([]byte, error) {
	if format == "" {
		format = p.TargetFormat
	}
	if format == "" {
		return nil, fmt.Errorf("no format specified for conversion")
	}

	var buf bytes.Buffer
	err := p.Encode(&buf, format)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Encode writes the image to the provided writer in the specified format.
// If format is empty, it uses the TargetFormat set on the processor.
func (p *Processor) Encode(w io.Writer, format string) error {
	if format == "" {
		format = p.TargetFormat
	}
	format = strings.TrimPrefix(strings.ToLower(format), ".")

	switch format {
	case "jpg", "jpeg":
		return jpeg.Encode(w, p.Img, &jpeg.Options{Quality: p.Quality})
	case "png":
		level := png.DefaultCompression
		if p.Quality <= 25 {
			level = png.BestSpeed
		} else if p.Quality >= 76 {
			level = png.BestCompression
		}
		encoder := png.Encoder{CompressionLevel: level}
		return encoder.Encode(w, p.Img)
	case "gif":
		return gif.Encode(w, p.Img, nil)
	case "bmp":
		return bmp.Encode(w, p.Img)
	case "tiff", "tif":
		return tiff.Encode(w, p.Img, nil)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// Image returns the current image.Image from the processor.
func (p *Processor) Image() image.Image {
	return p.Img
}

// --- Top-Level Helpers ---

// Convert loads an image from src and saves it to dst, handling format conversion.
func Convert(src, dst string) error {
	img, err := Load(src)
	if err != nil {
		return fmt.Errorf("load failed: %w", err)
	}

	return NewProcessor(img).Save(dst)
}
