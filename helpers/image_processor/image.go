package image_processor

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/nickalie/go-webpbin"
	"golang.org/x/image/draw"
)


func GenerateResizedImages(sourceImagePath, outputDirectory string) error {
	file, err := os.Open("./static/test.jpg")
	if err != nil {
		return err
	}

	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return err
	}

	sizes := []int{300, 450, 700, 950, 1200}

	err = os.MkdirAll(outputDirectory, os.ModePerm)
	if err != nil {
		return err
	}

	for _, size := range sizes {
		resizedImg := resizeImage(img, size)

		outputPath := filepath.Join(outputDirectory, fmt.Sprintf("resized_%d.webp", size))

		outputFile, err := os.Create(outputPath)

		if err != nil {
			return err
		}

		defer outputFile.Close()

		err = webpbin.Encode(outputFile, resizedImg)
		if err != nil {
			return err
		}

		fmt.Printf("Resized image saved at %s\n", outputPath)
	}

	return nil
}

func resizeImage(img image.Image, width int) image.Image {
	bounds := img.Bounds()
	aspectRatio := float64(bounds.Dx()) / float64(bounds.Dy())
	height := int(float64(width) / aspectRatio)

	resizedImg := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.ApproxBiLinear.Scale(resizedImg, resizedImg.Bounds(), img, bounds, draw.Over, nil)

	return resizedImg
}
