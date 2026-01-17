# Image Processing Library

A creative and robust Go library for image processing. Convert between formats, resize, compress, and apply artistic filters with a simple fluent API.

## Features

- **Generic Format Conversion**: Seamlessly convert between PNG, JPEG, GIF, BMP, TIFF, and **SVG**.
- **Explicit Target Formats**: Use `ToPNG()`, `ToJPEG()`, `ToGIF()`, `ToBMP()`, or `ToTIFF()` to set the output format explicitly.
- **Remote & Local Loading**: Load images from your disk or directly from a URL.
- **Explicit Convert Helper**: Use `img.Convert(src, dst)` for immediate file-to-file conversion.
- **High-Quality Resizing**: Powered by Catmull-Rom interpolation for professional results.
- **Compression Control**:
  - **JPEG**: Fine-tune quality from 1 to 100.
  - **PNG**: Choose between Default, Fast, Best, or No Compression.
- **Creative Filters**:
  - **Grayscale**: Classic black and white effect.
  - **Watermarking**: Add text overlays to protect your images.
  - **Transformations**: Rotate (90, 180 degrees) and Flip (Horizontal, Vertical).
  - **Brightness**: Adjust luminosity with a simple percentage.
- **Fluent API**: Chain methods for elegant, readable code.
- **Memory Efficient**: Export images as raw byte slices for direct streaming or database storage.

## Installation

```bash
go get github.com/user/image-lib
```

## Quick Start

```go
package main

import (
    "github.com/user/image-lib/pkg/img"
    "log"
)

func main() {
    // Load local or remote image
    image, err := img.Load("input.jpg")
    if err != nil {
        log.Fatal(err)
    }

    // Process and save
    err = img.NewProcessor(image).
        Resize(800, 600).
        Grayscale().
        Watermark("Protected").
        SetQuality(80).
        Save("output.png") // Format converted automatically

    if err != nil {
        log.Fatal(err)
    }
}
```

### Fluent Conversion

Convert and transform in one chain:

```go
err := img.NewProcessor(myImg).
    ToJPEG().
    SetQuality(85).
    Save("output.jpg")
```

### Explicit Conversion

The easiest way to convert an image from one format to another:

```go
err := img.Convert("input.png", "output.jpg")
```

### SVG Rasterization

Convert an SVG to a PNG or JPEG:

```go
err := img.Convert("logo.svg", "logo.png")
```

### Advanced Processing

## Memory flow (Downloads/Streaming)

Use `ToBytes` to handle images without saving to disk:

```go
data, err := img.NewProcessor(image).
    SetQuality(75).
    ToBytes("jpeg")

// serve 'data' as download or stream it
```

## Supported Transformations

| Method              | Description                                  |
| :------------------ | :------------------------------------------- |
| `Resize(w, h)`      | Scales image with high-quality interpolation |
| `Grayscale()`       | Converts to 8-bit grayscale                  |
| `Watermark(text)`   | Adds semi-transparent text to bottom-right   |
| `Rotate90()`        | 90° clockwise rotation                       |
| `Rotate180()`       | 180° rotation                                |
| `FlipHorizontal()`  | Mirror image horizontally                    |
| `FlipVertical()`    | Flip image upside down                       |
| `Brightness(%)`     | Adjust brightness (-100 to 100)              |
| `SetQuality(1-100)` | Output quality or compression level          |

### Export & Conversion

| Method              | Description                                     |
| :------------------ | :---------------------------------------------- |
| `Save(path)`        | Infers format from extension and saves to disk  |
| `ToBytes(format)`   | Returns image as `[]byte` in specified format   |
| `Encode(w, format)` | Writes image to `io.Writer` in specified format |

## License

MIT
