package api

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/image/draw"
)

func GenerateResizedImages(sourceImagePath, outputDirectory string) error {
	file, err := os.Open("./static/test.jpg")
	if err != nil {
		log.Println("1")
		return err
	}

	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Println("2")
		return err
	}

	sizes := []int{100, 200, 400, 800}

	err = os.MkdirAll(outputDirectory, os.ModePerm)
	if err != nil {
		log.Println("3")
		return err
	}

	for _, size := range sizes {
		resizedImg := resizeImage(img, size)

		outputPath := filepath.Join(outputDirectory, fmt.Sprintf("resized_%d.jpg", size))

		outputFile, err := os.Create(outputPath)
		if err != nil {
			log.Println("4")
			return err
		}
		defer outputFile.Close()

		err = jpeg.Encode(outputFile, resizedImg, nil)
		if err != nil {
			log.Println("5")
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
