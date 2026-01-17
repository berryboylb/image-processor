package img

import (
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// Load loads an image from a local path or a remote URL.
func Load(path string) (image.Image, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return loadRemote(path)
	}
	return loadLocal(path)
}

func loadLocal(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if strings.HasSuffix(strings.ToLower(path), ".svg") {
		return decodeSVG(file)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	return img, nil
}

func decodeSVG(ioReader io.Reader) (image.Image, error) {
	icon, err := oksvg.ReadIconStream(ioReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read SVG: %w", err)
	}

	w, h := int(icon.ViewBox.W), int(icon.ViewBox.H)
	if w <= 0 || h <= 0 {
		w, h = 512, 512 // Default fallback size
	}

	img := image.NewRGBA(image.Rect(0, 0, w, h))
	scanner := rasterx.NewScannerGV(w, h, img, img.Bounds())
	raster := rasterx.NewDasher(w, h, scanner)
	icon.Draw(raster, 1.0)

	return img, nil
}

func loadRemote(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download image, status code: %d", resp.StatusCode)
	}

	if strings.HasSuffix(strings.ToLower(url), ".svg") {
		return decodeSVG(resp.Body)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	return img, nil
}
