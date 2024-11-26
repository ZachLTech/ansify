package ansify

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
	"golang.org/x/term"
)

func loadImage(imagePath string) (image.Image, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	imgExt := strings.ToLower(filepath.Ext(imagePath))
	var img image.Image
	var decodeErr error

	switch imgExt {
	case ".jpeg", ".jpg":
		img, decodeErr = jpeg.Decode(file)
	case ".png":
		img, decodeErr = png.Decode(file)
	default:
		return nil, fmt.Errorf("unsupported image format: %s (only PNG or JPG supported)", imgExt)
	}

	if decodeErr != nil {
		return nil, fmt.Errorf("error decoding image: %v", decodeErr)
	}

	return img, nil
}

func resizeImage(img image.Image, width int) (image.Image, error) {
	if img == nil {
		return nil, fmt.Errorf("input image is nil")
	}

	if width <= 0 {
		return nil, fmt.Errorf("invalid width: %d", width)
	}

	bounds := img.Bounds()
	height := (bounds.Dy() * width) / bounds.Dx()
	if height <= 0 {
		return nil, fmt.Errorf("invalid calculated height: %d", height)
	}

	newImage := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(newImage, newImage.Bounds(), img, bounds, draw.Over, nil)
	return newImage, nil
}

func rgbToANSI(r, g, b uint8) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
}

func getAverageColor(img image.Image, x, y int) (uint8, uint8, uint8) {
	c1 := img.At(x, y)
	c2 := img.At(x, y+1)
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	return uint8((r1 + r2) / 512), uint8((g1 + g2) / 512), uint8((b1 + b2) / 512)
}

func mapToBlocks(img image.Image) ([]string, error) {
	if img == nil {
		return nil, fmt.Errorf("input image is nil")
	}

	bounds := img.Bounds()
	height, width := bounds.Max.Y, bounds.Max.X
	result := make([]string, height/2)
	const fullBlock = "█"
	const resetColor = "\x1b[0m"

	for y := bounds.Min.Y; y < height-1; y += 2 {
		var line strings.Builder
		for x := bounds.Min.X; x < width; x++ {
			r, g, b := getAverageColor(img, x, y)
			colorCode := rgbToANSI(r, g, b)
			line.WriteString(colorCode + fullBlock)
		}
		line.WriteString(resetColor)
		result[y/2] = line.String()
	}

	return result, nil
}

func GetAnsify(imageInput string) (string, error) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return "", fmt.Errorf("error getting terminal size: %v", err)
	}

	image, err := loadImage(imageInput)
	if err != nil {
		return "", fmt.Errorf("error loading image: %v", err)
	}

	resized, err := resizeImage(image, width)
	if err != nil {
		return "", fmt.Errorf("error resizing image: %v", err)
	}

	blockLines, err := mapToBlocks(resized)
	if err != nil {
		return "", fmt.Errorf("error mapping blocks: %v", err)
	}

	var sb strings.Builder
	for _, line := range blockLines {
		sb.WriteString(line + "\n")
	}

	return sb.String(), nil
}

func GetAnsifyCustomWidth(imageInput string, termWidth int) (string, error) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return "", fmt.Errorf("error getting terminal size: %v", err)
	}

	image, err := loadImage(imageInput)
	if err != nil {
		return "", fmt.Errorf("error loading image: %v", err)
	}

	resized, err := resizeImage(image, width)
	if err != nil {
		return "", fmt.Errorf("error resizing image: %v", err)
	}

	blockLines, err := mapToBlocks(resized)
	if err != nil {
		return "", fmt.Errorf("error mapping blocks: %v", err)
	}

	var sb strings.Builder
	for _, line := range blockLines {
		sb.WriteString(line + "\n")
	}

	return sb.String(), nil
}

func PrintAnsify(imageInput string) error {
	result, err := GetAnsify(imageInput)
	if err != nil {
		return err
	}
	fmt.Print(result)
	return nil
}
