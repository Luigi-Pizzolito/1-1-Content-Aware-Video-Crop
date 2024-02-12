package main

import (
	"image"
	"image/draw"
	"fmt"
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
                // croppedImg = ResizeWithVerticalLetterbox(croppedImg, squareSize, squareSize)
                croppedImg = ResizeRGBA(croppedImg, squareSize, squareSize)
            }

            if drawUI {
                croppedFrameRW.Lock()
                croppedFrame = croppedImg
                croppedFrameRW.Unlock()
            }
	        copyAndWrite(croppedImg)
        }
    }
}

// ---------------- Helper Functions ----------------

func cropToBBox(src *image.RGBA, cropRect image.Rectangle) *image.RGBA {
    // Create a new RGBA image for the cropped portion
    dst := image.NewRGBA(cropRect)

    // Copy the cropped portion from the source image to the destination image
    draw.Draw(dst, cropRect, src, cropRect.Min, draw.Src)

    return dst
}