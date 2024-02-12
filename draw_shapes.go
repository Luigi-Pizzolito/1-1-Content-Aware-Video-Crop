package main

import (
	"image"
	"image/color"
)

func drawRectangleOutline(img *image.RGBA, rect image.Rectangle, outlineColor color.RGBA, width int) {
    if drawUI && !playOnlyMode {
        // Draw top edge
        drawLine(img, rect.Min.X, rect.Min.Y, rect.Max.X, rect.Min.Y, outlineColor, width)
        // Draw bottom edge
        drawLine(img, rect.Min.X, rect.Max.Y, rect.Max.X, rect.Max.Y, outlineColor, width)
        // Draw left edge
        drawLine(img, rect.Min.X, rect.Min.Y, rect.Min.X, rect.Max.Y, outlineColor, width)
        // Draw right edge
        drawLine(img, rect.Max.X, rect.Min.Y, rect.Max.X, rect.Max.Y, outlineColor, width)
    }
}

func drawLine(img *image.RGBA, x0, y0, x1, y1 int, color color.RGBA, width int) {
    if drawUI && !playOnlyMode {
        // Bresenham's line algorithm to draw the line
        dx := abs(x1 - x0)
        dy := -abs(y1 - y0)
        sx := 1
        sy := 1
        if x0 >= x1 {
            sx = -1
        }
        if y0 >= y1 {
            sy = -1
        }
        err := dx + dy
        for {
            // Set the color of the pixel at the specified coordinates
            for i := x0 - width/2; i <= x0+width/2; i++ {
                for j := y0 - width/2; j <= y0+width/2; j++ {
                    img.SetRGBA(i, j, color)
                }
            }

            if x0 == x1 && y0 == y1 {
                break
            }
            e2 := 2 * err
            if e2 >= dy {
                err += dy
                x0 += sx
            }
            if e2 <= dx {
                err += dx
                y0 += sy
            }
        }
    }
}

func drawPoint(img *image.RGBA, p image.Point, color color.RGBA, radius int) {
    if drawUI && !playOnlyMode {
        x, y := p.X, p.Y
        for i := x - radius; i <= x+radius; i++ {
            for j := y - radius; j <= y+radius; j++ {
                // Calculate the distance from the center
                dx := x - i
                dy := y - j
                distance := dx*dx + dy*dy

                // If the distance is less than or equal to the radius squared, set the color
                if distance <= radius*radius {
                    img.SetRGBA(i, j, color)
                }
            }
        }
    }
}