package main

import (
	"fmt"
	"image"
	"image/draw"
)

// ---------------- Step Function ----------------

func cropVideo() {
	for {
		select {
		case mask, ok := <-procFrameQueue4:
			if !ok {
				// Channel is closed
				fmt.Println("Wrote Output Image Sequence")
				return
			}

			// Process the received image
			croppedImg := cropToBBox(mask.Image, mask.CBox)
			//? Experimental zoom out and aspect ratio when bbox > cropRect
			if mask.zoom {
				// squareWidth := squareSize
				// if widthSize != 0 {
				// 	squareWidth = widthSize
				// }

				// croppedImg = ResizeWithVerticalLetterbox(croppedImg, squareSize, squareSize)
				croppedImg = ResizeRGBA(croppedImg, squareSize, squareSize)
			}

			if drawUI {
				croppedFrameRW.Lock()
				croppedFrame = croppedImg
				croppedFrameRW.Unlock()
			}

			if !playOnlyMode {
				copyAndWrite(croppedImg)
			}
		}
	}
}

// ---------------- Helper Functions ----------------

func cropToBBox(src *image.RGBA, cropRect image.Rectangle) *image.RGBA {
	// Create a new RGBA image for the cropped portion
	dst := image.NewRGBA(cropRect)

	// Fill background with black
	// draw.Draw(dst, dst.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 255}}, image.ZP, draw.Src)

	// Copy the cropped portion from the source image to the destination image
	draw.Draw(dst, cropRect, src, cropRect.Min, draw.Src)

	return dst
}
