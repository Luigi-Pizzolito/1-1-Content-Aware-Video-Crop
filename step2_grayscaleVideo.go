package main

import (
	"image"
	"image/color"
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
            if drawUI {
                grayFrameRW.Lock()
                grayFrame = convertGrayToRGBA(maskFrame)
                grayFrameRW.Unlock()
            }
			procFrameQueue2 <- &ImageGrayMask{Mask:maskFrame,Image:img}
        }
    }
}

// ---------------- Helper Functions ----------------

// Function to convert a color to grayscale
func toGray(c color.Color) color.Color {
    r, g, b, _ := c.RGBA()

    // Calculate the grayscale value
    gray := uint8((19595*r + 38470*g + 7471*b + 1<<15) >> 24)

    return color.RGBA{gray, gray, gray, 255}
}

func convertToMonochrome(src *image.RGBA) *image.Gray {
    bounds := src.Bounds()
    width, height := bounds.Dx(), bounds.Dy()

    // Create a new RGBA image
    monochromeImg := image.NewGray(image.Rect(0, 0, width, height))

    // Iterate over each pixel in the original image
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            // Get the color of the pixel
            pixelColor := src.At(x, y)

            // Convert the color to grayscale
            grayValue := toGray(pixelColor)

            // Set the pixel color in the monochrome image
            monochromeImg.Set(x, y, grayValue)
        }
    }

    return monochromeImg
}

func convertGrayToRGBA(grayImg *image.Gray) *image.RGBA {
    bounds := grayImg.Bounds()
    width, height := bounds.Dx(), bounds.Dy()

    // Create a new RGBA image with the same bounds as the grayscale image
    rgbaImg := image.NewRGBA(image.Rect(0, 0, width, height))

    // Convert grayscale image to RGBA using a cast
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            // Cast the grayscale value to uint8 and set it to all RGB channels
            value := grayImg.GrayAt(x, y).Y
            rgbaImg.Pix[y*rgbaImg.Stride+x*4] = value
            rgbaImg.Pix[y*rgbaImg.Stride+x*4+1] = value
            rgbaImg.Pix[y*rgbaImg.Stride+x*4+2] = value
            rgbaImg.Pix[y*rgbaImg.Stride+x*4+3] = 255 // Alpha channel (fully opaque)
        }
    }

    return rgbaImg
}