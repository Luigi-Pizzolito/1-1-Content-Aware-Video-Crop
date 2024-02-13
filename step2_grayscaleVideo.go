package main

import (
	"image"
    "image/draw"
)

// ---------------- Step Function ----------------

func grayscaleVideo() {
	for {
        select {
        case img, ok := <-procFrameQueue1:
            if !ok {
                // Channel is closed
                return
            }

            // Process the received image
			maskFrame := convertToMonochrome(img)
            if drawUI && !playOnlyMode {
                grayFrameRW.Lock()
                grayFrame = convertGrayToRGBA(maskFrame)
                grayFrameRW.Unlock()
            }
			procFrameQueue2 <- &ImageGrayMask{Mask:maskFrame,Image:img}
        }
    }
}

// ---------------- Helper Functions ----------------

func convertToMonochrome(img *image.RGBA) *image.Gray {
    monoFrame := image.NewGray(img.Bounds())
    draw.Draw(monoFrame, monoFrame.Bounds(), img, img.Bounds().Min, draw.Src)
    return monoFrame
}

func convertGrayToRGBA(img *image.Gray) *image.RGBA {
    rgbaFrame := image.NewRGBA(img.Bounds())
    draw.Draw(rgbaFrame, rgbaFrame.Bounds(), img, img.Bounds().Min, draw.Src)
    return rgbaFrame
}