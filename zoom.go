package main

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/nfnt/resize"
)

// ResizeWithVerticalLetterbox resizes the given image while preserving its aspect ratio and adding vertical letterbox bars.
func ResizeWithVerticalLetterbox(img image.Image, width, height int) *image.RGBA {
	// Calculate aspect ratios
	srcWidth := img.Bounds().Dx()
	srcHeight := img.Bounds().Dy()
	srcAspectRatio := float64(srcWidth) / float64(srcHeight)
	dstAspectRatio := float64(width) / float64(height)

	// Initialize new dimensions
	newWidth := width
	newHeight := height

	// Determine whether to scale based on width or height
	if srcAspectRatio > dstAspectRatio {
		// Scale based on width
		newHeight = int(float64(width) / srcAspectRatio)
	} else {
		// Scale based on height
		newWidth = int(float64(height) * srcAspectRatio)
	}

	// Resize the image
	resizedImg := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)

	// Calculate the remaining vertical space after resizing
	remainingHeight := height - newHeight

	// Calculate the height of letterbox bars above and below the image
	letterboxHeight := remainingHeight / 2

	// Create a new RGBA image with the desired dimensions
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(dst, binaryFrame.Bounds(), &image.Uniform{color.RGBA{0,255,0,255}}, image.ZP, draw.Src)

	// Calculate the position to draw the resized image
	x := (width - newWidth) / 2
	y := (height - newHeight) / 2

	// Draw the resized image onto the destination image
	draw.Draw(dst, image.Rect(x, y, x+newWidth, y+newHeight), resizedImg, image.Point{}, draw.Src)

	// Fill letterbox bars with a specified color (e.g., black) above and below the image
	letterboxColor := color.RGBA{0, 255, 0, 255}
	draw.Draw(dst, image.Rect(0, 0, width, letterboxHeight), &image.Uniform{letterboxColor}, image.ZP, draw.Src)
	draw.Draw(dst, image.Rect(0, height-letterboxHeight, width, height), &image.Uniform{letterboxColor}, image.ZP, draw.Src)

	return RemoveAlphaRGBA(dst)
}

// ResizeRGBA scales down the given image.RGBA while preserving its aspect ratio.
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
	resizedImg := resize.Resize(uint(float64(img.Bounds().Dx())*scale), uint(float64(img.Bounds().Dy())*scale), img, resize.Lanczos3)

	// Create a new RGBA image with the resized dimensions
	resizedRGBA := image.NewRGBA(image.Rect(0, 0, int(resizedImg.Bounds().Dx()), int(resizedImg.Bounds().Dy())))

	// Copy the resized image onto the new RGBA image
	draw.Draw(resizedRGBA, resizedRGBA.Bounds(), resizedImg, resizedImg.Bounds().Min, draw.Src)

	return RemoveAlphaRGBA(resizedRGBA)
}

// RemoveAlphaRGBA removes the alpha channel from an RGBA image and returns a new RGBA image.
func RemoveAlphaRGBA(rgba *image.RGBA) *image.RGBA {
	// Create a new RGBA image with the same dimensions as the original image
	rgb := image.NewRGBA(image.Rect(0, 0, rgba.Bounds().Dx(), rgba.Bounds().Dy()))

	// Copy RGB values from the original image to the new image while discarding alpha channel
	for y := rgba.Bounds().Min.Y; y < rgba.Bounds().Max.Y; y++ {
		for x := rgba.Bounds().Min.X; x < rgba.Bounds().Max.X; x++ {
			r, g, b, _ := rgba.At(x, y).RGBA()
			rgb.SetRGBA(x, y, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 255})
		}
	}

	return rgb
}

// ZoomOutRGBA zooms out the given image.RGBA by a percentage while keeping the same aspect ratio.
// It fills the background with black.
func ZoomOutRGBA(img *image.RGBA, percentage float64) *image.RGBA {
	// Calculate the dimensions after zooming out
	newWidth := int(float64(img.Bounds().Dx()) * (1 - percentage/100))
	newHeight := int(float64(img.Bounds().Dy()) * (1 - percentage/100))

	// Resize the image using the calculated dimensions
	resizedImg := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)

	// Create a new RGBA image with the original dimensions
	resizedRGBA := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))

	// Calculate the position to draw the resized image (centered)
	x := (img.Bounds().Dx() - newWidth) / 2
	y := (img.Bounds().Dy() - newHeight) / 2

	// Fill the background with black
	draw.Draw(resizedRGBA, resizedRGBA.Bounds(), &image.Uniform{color.Black}, image.ZP, draw.Src)

	// Draw the resized image onto the new RGBA image
	draw.Draw(resizedRGBA, image.Rect(x, y, x+newWidth, y+newHeight), resizedImg, resizedImg.Bounds().Min, draw.Src)

	return resizedRGBA
}