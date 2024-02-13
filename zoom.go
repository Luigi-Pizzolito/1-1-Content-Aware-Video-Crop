package main

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/bamiaux/rez"
)

func ResizeRGBA(img *image.RGBA, maxWidth, maxHeight int) *image.RGBA {
	// Calculate the scaling factor for width and height
	scaleWidth := float64(maxWidth) / float64(img.Bounds().Dx())
	scaleHeight := float64(maxHeight) / float64(img.Bounds().Dy())

	// Use the smaller scaling factor to ensure the entire image fits within the specified dimensions
	scale := scaleWidth
	if scaleHeight < scaleWidth {
		scale = scaleHeight
	}

	// Resize the image using the calculated scaling factor
	resizedImg := image.NewRGBA(image.Rect(0, 0, int(float64(img.Bounds().Dx())*scale), int(float64(img.Bounds().Dy())*scale)))
	err := rez.Convert(resizedImg, img, rez.NewLanczosFilter(2))
	if err != nil {
		panic("Failed to resize img: "+err.Error())
	}

	// Create a new RGBA image with the resized dimensions
	resizedRGBA := image.NewRGBA(image.Rect(0, 0, int(resizedImg.Bounds().Dx()), int(resizedImg.Bounds().Dy())))
	// Fill with black
	draw.Draw(resizedRGBA, resizedRGBA.Bounds(), &image.Uniform{color.Black}, image.ZP, draw.Src)

	// Copy the resized image onto the new RGBA image, with alpha blending
	draw.Draw(resizedRGBA, resizedRGBA.Bounds(), resizedImg, resizedImg.Bounds().Min, draw.Over)

	return resizedRGBA
}

// ZoomOutRGBA zooms out the given image.RGBA by a percentage while keeping the same aspect ratio.
// It fills the background with black.
func ZoomOutRGBA(img *image.RGBA, percentage float64) *image.RGBA {
	// Calculate the dimensions after zooming out
	newWidth := int(float64(img.Bounds().Dx()) * (1 - percentage/100))
	newHeight := int(float64(img.Bounds().Dy()) * (1 - percentage/100))

	// Resize the image using the calculated dimensions
	resizedImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	err := rez.Convert(resizedImg, img, rez.NewLanczosFilter(2))
	if err != nil {
		panic("Failed to resize img: "+err.Error())
	}

	// Create a new RGBA image with the original dimensions
	resizedRGBA := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))

	// Calculate the position to draw the resized image (centered)
	x := (img.Bounds().Dx() - newWidth) / 2
	y := (img.Bounds().Dy() - newHeight) / 2

	// Fill the background with black
	draw.Draw(resizedRGBA, resizedRGBA.Bounds(), &image.Uniform{color.Black}, image.ZP, draw.Src)

	// Draw the resized image onto the new RGBA image
	draw.Draw(resizedRGBA, image.Rect(x, y, x+newWidth, y+newHeight), resizedImg, resizedImg.Bounds().Min, draw.Over)

	return resizedRGBA
}