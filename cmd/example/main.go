package main

import (
	"fmt"
	"log"
	"os"

	"github.com/berryboylb/image-processor/pkg/img"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: example <image_path_or_url>")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	fmt.Printf("Loading image from: %s\n", inputPath)

	image, err := img.Load(inputPath)
	if err != nil {
		log.Fatalf("Failed to load image: %v", err)
	}

	// 1. Explicit Conversion using the fluent API (png -> jpeg)
	fmt.Println("Processing: Explicit Conversion via .ToPNG().Save()")
	err = img.NewProcessor(image).
		ToPNG().
		Save("explicit_output.png")
	if err != nil {
		log.Printf("Failed to convert image via fluent API: %v", err)
	}

	// 2. Explicit Conversion using high-level Convert helper
	fmt.Println("Processing: Explicit Conversion (JPEG -> PNG) via img.Convert")

	// 2. SVG to PNG Conversion (if possible)
	// I'll create a dummy SVG for demonstration if the user didn't provide one
	svgPath := "test.svg"
	svgContent := `<svg height="100" width="100"><circle cx="50" cy="50" r="40" stroke="black" stroke-width="3" fill="red" /></svg>`
	_ = os.WriteFile(svgPath, []byte(svgContent), 0644)

	fmt.Println("Processing: SVG rasterization -> output_from_svg.png")
	err = img.Convert(svgPath, "output_from_svg.png")
	if err != nil {
		log.Printf("Failed to convert SVG: %v", err)
	}

	// 3. Resize, Watermark and Save as JPEG (Quality adjustment)
	fmt.Println("Processing: Resize (300x300), Watermark, Quality(50) -> output_low_quality.jpg")
	err = img.NewProcessor(image).
		Resize(300, 300).
		Watermark("Reduced Quality").
		SetQuality(50).
		Save("output_low_quality.jpg")
	if err != nil {
		log.Printf("Failed to save watermarked image: %v", err)
	}

	// 4. Creative transformations: Rotate 180, Flip Vertical
	fmt.Println("Processing: Rotate180, FlipVertical -> output_refined.png")
	err = img.NewProcessor(image).
		Rotate180().
		FlipVertical().
		Save("output_refined.png")
	if err != nil {
		log.Printf("Failed to save refined image: %v", err)
	}

	// 4. Memory Output (ToBytes) - Demonstrating "Download" flow
	fmt.Println("Processing: ToBytes (Memory buffer) for download simulation")
	data, err := img.NewProcessor(image).
		Grayscale().
		ToBytes("jpeg")
	if err != nil {
		log.Fatalf("Failed to get image bytes: %v", err)
	}
	fmt.Printf("Generated %d bytes of JPEG data in memory (ready for download/streaming)\n", len(data))

	fmt.Println("Done!")
}
