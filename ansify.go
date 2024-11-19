package ansify

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/draw"
	"golang.org/x/term"
)

func loadImage(image string) image.Image {
	file, err := os.Open(image)
	if err != nil {
		fmt.Printf("error while opening file %v\n", err)
	}
	defer file.Close()

	imgExt := filepath.Ext(image)

	if imgExt == ".jpeg" || imgExt == ".jpg" {
		img, err := jpeg.Decode(file)
		if err != nil {
			fmt.Printf("error while decoding image %v\n", err)
		}

		return img
	} else if imgExt == ".png" {
		img, err := png.Decode(file)
		if err != nil {
			fmt.Printf("error while decoding image %v\n", err)
		}

		return img
	} else {
		panic("There was an error and the file provided was likely not the correct format. (Only PNG or JPG)\n")
	}
}

func resizeImage(img image.Image, width int) image.Image {
	bounds := img.Bounds()
	height := (bounds.Dy() * width) / bounds.Dx()
	newImage := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(newImage, newImage.Bounds(), img, bounds, draw.Over, nil)
	return newImage
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

func mapToBlocks(img image.Image) []string {
	bounds := img.Bounds()
	height, width := bounds.Max.Y, bounds.Max.X
	result := make([]string, height/2)

	const fullBlock = "█"
	const resetColor = "\x1b[0m"

	for y := bounds.Min.Y; y < height-1; y += 2 {
		line := ""
		for x := bounds.Min.X; x < width; x++ {
			r, g, b := getAverageColor(img, x, y)
			colorCode := rgbToANSI(r, g, b)
			line += colorCode + fullBlock
		}
		line += resetColor
		result[y/2] = line
	}

	return result
}

func GetAnsify(imageInput string) string {
	var renderedImage string
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Printf("error while getting terminal size %v\n", err)
	}

	image := loadImage(imageInput)
	image = resizeImage(image, width)

	blockLines := mapToBlocks(image)
	for _, line := range blockLines {
		renderedImage += fmt.Sprintf(line)
	}

	return renderedImage
}

func PrintAnsify(imageInput string) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Printf("error while getting terminal size %v\n", err)
	}

	image := loadImage(imageInput)
	image = resizeImage(image, width)

	blockLines := mapToBlocks(image)
	for _, line := range blockLines {
		fmt.Println(line)
	}
}
