package main

import (
	"image"
	"image/color"
	"sort"
)

// ---------------- Step Function ----------------

func binaryVideo() {
	for {
        select {
        case mask, ok := <-procFrameQueue2:
            if !ok {
                // Channel is closed
                return
            }

            // Process the received image
			threshold := otsuThreshold(mask.Mask)
			binaryImg := applyThreshold(mask.Mask, threshold)
            binaryFrameM := convertGrayToRGBA(binaryImg)

            // Apply connected component labeling algorithm
	        boundingBoxes, _ := connectedComponentLabeling(binaryImg)
            // Print bounding boxes and label matrix
            // Convert map to slice of rectangles
            var rectangles []image.Rectangle
            for _, rect := range boundingBoxes {
                rectangles = append(rectangles, rect)
            }
            sort.Sort(ByArea(rectangles))

            subjects := make([]image.Rectangle, 0)
            computedOverlapRects := rectangles
            for _, bbox := range rectangles {
                largestArea := rectangles[0].Dx()*rectangles[0].Dy()
                area := bbox.Dx()*bbox.Dy()
                perCovered := 100*float64(float64(area)/float64(largestArea))
                // boxes classified as big have at least 40% of the area of the biggest box
                if perCovered >= 20 {
                    
                    drawRectangleOutline(binaryFrameM, bbox, color.RGBA{255,0,255,255},1)

                    // merge any other boxes that overlap big box into big bound
                    bbox, computedOverlapRects = CombineBounds(bbox, computedOverlapRects, mask.Image.Bounds().Dy()/20)
                    drawRectangleOutline(binaryFrameM, bbox, color.RGBA{0,255,255,255},1)

                    ctr := calculateCentroid(binaryImg, bbox)
                    drawPoint(binaryFrameM, ctr, color.RGBA{0,255,255,255}, 3)

                    subjects = append(subjects, bbox)
                } else {
                    drawRectangleOutline(binaryFrameM, bbox, color.RGBA{170,0,255,255},1)
                }
            }
            

            subjects_bbox := TotalBoundingBox(subjects)
            //? extra padding??
            // subjects_bbox = ExpandRectangle(subjects_bbox, mask.Image.Bounds().Dy()/20)
            


            if drawUI && !playOnlyMode {
				binaryFrameRW.Lock()
				binaryFrame = binaryFrameM
				binaryFrameRW.Unlock()
			}
            subj = len(subjects)
			procFrameQueue3 <- &ImageGrayBoundMask{Mask:binaryImg,Image:mask.Image,Bbox:subjects_bbox,Subj:subj}
        }
    }
}

// ---------------- Helper Functions ----------------

func otsuThreshold(img *image.Gray) uint8 {
    // Calculate histogram
    histogram := make([]int, 256)
    totalPixels := float64(img.Bounds().Dx() * img.Bounds().Dy())
    for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
        for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
            intensity := img.GrayAt(x, y).Y
            histogram[intensity]++
        }
    }

    // Normalize histogram
    var normalizedHistogram [256]float64
    for i := range histogram {
        normalizedHistogram[i] = float64(histogram[i]) / totalPixels
    }

    // Calculate cumulative sum and cumulative mean
    var cumulativeSum [256]float64
    var cumulativeMean [256]float64
    cumulativeSum[0] = normalizedHistogram[0]
    cumulativeMean[0] = 0
    for i := 1; i < 256; i++ {
        cumulativeSum[i] = cumulativeSum[i-1] + normalizedHistogram[i]
        cumulativeMean[i] = cumulativeMean[i-1] + float64(i)*normalizedHistogram[i]
    }

    // Calculate between-class variance and find optimal threshold
    var maxBetweenClassVariance float64
    var optimalThreshold uint8
    for t := 0; t < 256; t++ {
        omega := cumulativeSum[t]
        mu := cumulativeMean[t]
        omega2 := 1 - omega
        mu2 := (cumulativeMean[255] - mu) / omega2
        betweenClassVariance := omega * omega2 * (mu - mu2) * (mu - mu2)
        if betweenClassVariance > maxBetweenClassVariance {
            maxBetweenClassVariance = betweenClassVariance
            optimalThreshold = uint8(t)
        }
    }

    return optimalThreshold
}

func applyThreshold(img *image.Gray, threshold uint8) *image.Gray {
    binaryImg := image.NewGray(img.Bounds())
    black := color.Gray{Y: 0}
    white := color.Gray{Y: 255}

    for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
        for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
            if img.GrayAt(x, y).Y < threshold {
                binaryImg.SetGray(x, y, black)
            } else {
                binaryImg.SetGray(x, y, white)
            }
        }
    }
    return binaryImg
}

// TotalBoundingBox finds the total bounding box given a slice of image.Rectangles.
func TotalBoundingBox(rectangles []image.Rectangle) image.Rectangle {
	if len(rectangles) == 0 {
		// Return zero rectangle if there are no input rectangles
		return image.ZR
	}

	// Initialize min and max coordinates
	minX, minY := rectangles[0].Min.X, rectangles[0].Min.Y
	maxX, maxY := rectangles[0].Max.X, rectangles[0].Max.Y

	// Iterate over rectangles to find min and max coordinates
	for _, rect := range rectangles {
		if rect.Min.X < minX {
			minX = rect.Min.X
		}
		if rect.Min.Y < minY {
			minY = rect.Min.Y
		}
		if rect.Max.X > maxX {
			maxX = rect.Max.X
		}
		if rect.Max.Y > maxY {
			maxY = rect.Max.Y
		}
	}

	// Construct and return the total bounding box
	return image.Rect(minX, minY, maxX, maxY)
}

// ExpandRectangle expands the given rectangle by x pixels from its center.
func ExpandRectangle(rect image.Rectangle, x int) image.Rectangle {
	// Expand the rectangle
	newMinX := rect.Min.X - x
	newMinY := rect.Min.Y - x
	newMaxX := rect.Max.X + x
	newMaxY := rect.Max.Y + x

	// Return the expanded rectangle
	return image.Rect(newMinX, newMinY, newMaxX, newMaxY)
}

// CombineBounds combines the bounds of all rectangles in the slice
// that overlap with the given rectangle and removes them from the slice.
func CombineBounds(rect image.Rectangle, rects []image.Rectangle, tol int) (combined image.Rectangle, remainingRects []image.Rectangle) {
    // Initialize the combined rectangle with the rectangle passed as an argument
    combined = rect

    // Iterate over the rectangles in the slice
    for _, r := range rects {
        // Check if the current rectangle overlaps with the given rectangle
        if ExpandRectangle(combined, tol).Overlaps(r) {
            // Combine the bounds of overlapping rectangles
            combined = combined.Union(r)
        } else {
            // Add the non-overlapping rectangle to the remaining rectangles slice
            remainingRects = append(remainingRects, r)
        }
    }

    return combined, remainingRects
}

// Connected component labeling algorithm
func connectedComponentLabeling(binaryMap *image.Gray) (map[int]image.Rectangle, [][]int) {
	labelMatrix := make([][]int, binaryMap.Bounds().Dy())
	for i := range labelMatrix {
		labelMatrix[i] = make([]int, binaryMap.Bounds().Dx())
	}

	boundingBoxes := make(map[int]image.Rectangle)

	label := 1
	for y := binaryMap.Bounds().Min.Y; y < binaryMap.Bounds().Max.Y; y++ {
		for x := binaryMap.Bounds().Min.X; x < binaryMap.Bounds().Max.X; x++ {
			if binaryMap.GrayAt(x, y).Y > 0 && labelMatrix[y][x] == 0 {
				boundingBox := dfs(binaryMap, labelMatrix, x, y, label)
				boundingBoxes[label] = boundingBox
				label++
			}
		}
	}

	return boundingBoxes, labelMatrix
}

// Depth-first search to label connected components
func dfs(binaryMap *image.Gray, labelMatrix [][]int, x, y, label int) image.Rectangle {
    minX, minY, maxX, maxY := x, y, x, y

    stack := []image.Point{{x, y}}
    for len(stack) > 0 {
        current := stack[len(stack)-1]
        stack = stack[:len(stack)-1]

        cx, cy := current.X, current.Y
        labelMatrix[cy][cx] = label

        if cx < minX {
            minX = cx
        }
        if cx > maxX {
            maxX = cx
        }
        if cy < minY {
            minY = cy
        }
        if cy > maxY {
            maxY = cy
        }

        // Check neighboring pixels
        for dx := -1; dx <= 1; dx++ {
            nx, ny := cx+dx, cy-1
            if nx >= 0 && nx < binaryMap.Bounds().Dx() && ny >= 0 && ny < binaryMap.Bounds().Dy() {
                if binaryMap.GrayAt(nx, ny).Y > 0 && labelMatrix[ny][nx] == 0 {
                    stack = append(stack, image.Point{nx, ny})
                }
            }
            ny = cy + 1
            if nx >= 0 && nx < binaryMap.Bounds().Dx() && ny >= 0 && ny < binaryMap.Bounds().Dy() {
                if binaryMap.GrayAt(nx, ny).Y > 0 && labelMatrix[ny][nx] == 0 {
                    stack = append(stack, image.Point{nx, ny})
                }
            }
        }
        nx := cx - 1
        ny := cy
        if nx >= 0 && nx < binaryMap.Bounds().Dx() && ny >= 0 && ny < binaryMap.Bounds().Dy() {
            if binaryMap.GrayAt(nx, ny).Y > 0 && labelMatrix[ny][nx] == 0 {
                stack = append(stack, image.Point{nx, ny})
            }
        }
        nx = cx + 1
        if nx >= 0 && nx < binaryMap.Bounds().Dx() && ny >= 0 && ny < binaryMap.Bounds().Dy() {
            if binaryMap.GrayAt(nx, ny).Y > 0 && labelMatrix[ny][nx] == 0 {
                stack = append(stack, image.Point{nx, ny})
            }
        }
    }

    return image.Rect(minX, minY, maxX+1, maxY+1)
}

// ByArea implements the sort.Interface for []image.Rectangle based on area.
type ByArea []image.Rectangle

func (a ByArea) Len() int           { return len(a) }
func (a ByArea) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByArea) Less(i, j int) bool { return a[i].Dx()*a[i].Dy() > a[j].Dx()*a[j].Dy() }

// CenterPoint returns the center point of the given rectangle.
func CenterPoint(rect image.Rectangle) image.Point {
	return image.Point{
		X: (rect.Min.X + rect.Max.X) / 2,
		Y: (rect.Min.Y + rect.Max.Y) / 2,
	}
}

func calculateCentroid(img *image.Gray, bbox image.Rectangle) image.Point {
    totalX, totalY, count := 0, 0, 0
    for y := bbox.Min.Y; y < bbox.Max.Y; y++ {
        for x := bbox.Min.X; x < bbox.Max.X; x++ {
            if img.GrayAt(x, y).Y == 255 { // white pixel
                totalX += x
                totalY += y
                count++
            }
        }
    }

    if count > 0 {
        return image.Point{totalX / count, totalY / count}
    }

    return image.Point{0, 0}
}