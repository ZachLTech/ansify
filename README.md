# Ansify

Ansify is a Go package that converts PNG and JPEG images into ANSI terminal art. It automatically resizes images to fit your terminal width while maintaining aspect ratio, and renders them using colored block characters for a visually appealing result.

I mostly made this as I needed something exactly like it in order to properly complete a [TUI game](https://github.com/ZachLTech/eXit) I'm working on :)

## Features

- Supports PNG and JPEG image formats
- Automatically resizes images to fit terminal width
- Maintains aspect ratio during resizing
- Uses RGB color codes for accurate color representation
- Smooth image scaling using Catmull-Rom algorithm
- Returns the image as a string or prints directly to terminal

## Installation

```bash
go get github.com/ZachLTech/ansify
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "github.com/ZachLTech/ansify"
)

func main() {
    // Print image directly to terminal
    ansify.PrintAnsify("path/to/your/image.jpg")

    // Or get the ANSI string to use elsewhere
    ansiString := ansify.GetAnsify("path/to/your/image.png")
    fmt.Print(ansiString)
}
```

### Functions

#### `PrintAnsify(imageInput string)`
Loads the image, processes it, and prints directly to the terminal.

#### `GetAnsify(imageInput string) string`
Returns the processed image as an ANSI string that can be used elsewhere in your application.

## How It Works

1. The image is loaded and decoded based on its file extension (.png or .jpg/.jpeg)
2. The image is resized to match your terminal width while maintaining aspect ratio
3. Each pixel is converted to RGB values and mapped to ANSI color codes
4. The image is rendered using full block characters (█) with appropriate colors
5. Each pair of vertical pixels is combined into a single block to maintain proper aspect ratio in the terminal

## Limitations

- Only supports PNG and JPEG formats
- Requires a terminal that supports RGB ANSI color codes
- Image quality depends on terminal font, size, and color support

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project falls under the MIT license available [here](./LICENSE)
