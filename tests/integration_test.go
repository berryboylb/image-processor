package tests

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"testing"

	"github.com/berryboylb/image-processor/pkg/img"
)

func TestImageProcessingFlow(t *testing.T) {
	// 1. Setup - Create a test image
	testImgPath := "test_input.png"
	createTestImage(testImgPath)
	defer os.Remove(testImgPath)

	// 2. Test Loading
	image, err := img.Load(testImgPath)
	if err != nil {
		t.Fatalf("Failed to load image: %v", err)
	}

	// 3. Test Transformations and Fluent API
	processor := img.NewProcessor(image)
	err = processor.
		Resize(100, 100).
		Grayscale().
		Rotate180().
		FlipVertical().
		SetQuality(80).
		ToPNG().
		Save("test_output.png")

	if err != nil {
		t.Errorf("Failed to process image: %v", err)
	}
	defer os.Remove("test_output.png")

	// 4. Test Explicit Conversion Helper
	err = img.Convert(testImgPath, "test_converted.jpg")
	if err != nil {
		t.Errorf("Failed to convert image via helper: %v", err)
	}
	defer os.Remove("test_converted.jpg")

	// 5. Test In-Memory Conversion (ToBytes)
	pngBytes, err := processor.ToPNG().ToBytes("")
	if err != nil {
		t.Errorf("Failed to convert to bytes: %v", err)
	}
	if len(pngBytes) == 0 {
		t.Error("Generated empty byte slice for PNG")
	}

	// 6. Test SVG Rasterization
	svgPath := "test_input.svg"
	svgContent := `<svg height="100" width="100"><circle cx="50" cy="50" r="40" stroke="black" stroke-width="3" fill="red" /></svg>`
	_ = os.WriteFile(svgPath, []byte(svgContent), 0644)
	defer os.Remove(svgPath)

	err = img.Convert(svgPath, "test_svg_output.png")
	if err != nil {
		t.Errorf("Failed to rasterize SVG: %v", err)
	}
	defer os.Remove("test_svg_output.png")
}

func createTestImage(path string) {
	rect := image.Rect(0, 0, 200, 200)
	img := image.NewRGBA(rect)
	draw.Draw(img, rect, &image.Uniform{color.RGBA{255, 0, 0, 255}}, image.Point{}, draw.Src)

	f, _ := os.Create(path)
	defer f.Close()
	_ = png.Encode(f, img) // Changed from img.Encode(f, "png") to png.Encode(f, img)
}
