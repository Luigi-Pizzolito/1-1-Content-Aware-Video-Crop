package main

import (
	"image"
	"image/color"
)

var (
	// ma				*MovingAverage
	pid             *PIDController
	pidz            *PIDController
)

// ---------------- Step Function ----------------

func calcVideo() {
    // pid = NewPIDController(0.1, 0.000000000005, 0.00000002, 0) // nice oscillation at 0.1, 0.000000000005, 0.00000002
    //                     P    I             D            Initial Setpoint
    pid = NewPIDController(0.1, 0.001*10e-10, 0.002*10e10, 0)
    pidz = NewPIDController(0.2, 0.003*10e-10, 0.002*10e10, 0)

	for {
        select {
        case mask, ok := <-procFrameQueue3:
            if !ok {
                // Channel is closed
                return
            }

            // Process the received image
			// bbox := findBoundingBox(mask.Mask)
            bbox := mask.Bbox
            var centroid image.Point
            // if mask.Subj > 1 {
                // Multi subject, center is more balanced
                // centroid = CenterPoint(bbox)
            // } else {
                // Single subject, centroid is more balanced
                centroid = calculateCentroid(mask.Mask, bbox)
            // }
            
			panAdj := calculatePanAdjustment(centroid,mask.Image.Bounds().Dx())
			bboxArea := bbox.Dx()*bbox.Dy()
            frameArea := mask.Image.Bounds().Dx()*mask.Image.Bounds().Dy()
            perCovered := float64(float64(bboxArea)/float64(frameArea))*100
            bboxC := color.RGBA{0,0,255,255}
            if perCovered < 0.1 {
				// bounding box too small (< .1%), do not pan
				panAdj = 0
                bboxC = color.RGBA{255,0,255,255}
			} else if perCovered > 75 {
                // bounding box too big (> 80%), do not pan
                panAdj = 0
                bboxC = color.RGBA{255,0,255,255}
            }
            // fmt.Println(bboxArea, float64(float64(bboxArea)/float64(frameArea)))
			cropRect := calculateCropRect(centroid,mask.Image.Bounds().Dx(),mask.Image.Bounds().Dy(),panAdj)

            //? Experimental zoom out if width of bbox > cropRect
            zoom := false
            padding := int(cropRect.Dx()/8) // padding width of 1/10th of the crop rectangle width
            perOut := float64(float64(bbox.Dx()) / float64(cropRect.Dx())) * 100
            difOut := minBoxHDist(bbox, cropRect) //bbox.Dx()-cropRect.Dx()
            difOut = clamp(difOut+padding, 0, cropRect.Dx()*5)
            zoompx := 0
            targetCropRect := cropRect
            zoompx = difOut
            pidz.SetSetpoint(float64(zoompx))
            smoothedZoomAdj := pidz.Update()
            pidz.SetDelta(pidz.delta+smoothedZoomAdj)
            // if perOut > 90 {
                // fmt.Println(zoompx, int(pidz.delta), int(smoothedZoomAdj))
                targetCropRect = ExpandRectangle(cropRect,clamp(zoompx, 0, cropRect.Dx()*5))
                cropRect = ExpandRectangle(cropRect, clamp(int(pidz.delta), 0, cropRect.Dx()*5))
                zoom = true
            // }
            // zoom = true //! test
            // fmt.Println(perOut, difOut)
            //----

            
			var calcImg *image.RGBA
			if drawUI {
				calcImg = convertGrayToRGBA(mask.Mask)
				
				pt := 2
				drawRectangleOutline(calcImg, bbox, bboxC, pt)
				drawPoint(calcImg, centroid, color.RGBA{255,0,0,255}, pt*3)
				
				drawRectangleOutline(calcImg, cropRect, color.RGBA{0,255,0,255}, pt)
				RectCtrX := cropRect.Min.X + (cropRect.Max.X - cropRect.Min.X) / 2
				drawLine(calcImg, RectCtrX, 0, RectCtrX, mask.Image.Bounds().Dy(), color.RGBA{0,255,0,255},1)
				drawLine(calcImg, centroid.X, 0, centroid.X, mask.Image.Bounds().Dy(), color.RGBA{255,0,0,255}, 1)

				if zoom {
					drawRectangleOutline(calcImg, targetCropRect, color.RGBA{255,255,0,255}, 1)
					// fmt.Println(posF(perOut-100))
					calcImg = ZoomOutRGBA(calcImg, posF(perOut-100)) //clampF(100-posF(perOut)
				}
			}

            // fmt.Printf("Pan: %d, Zoom: %.2f\n", int(pid.delta), 100-posF(perOut-100))
            ppan = int(pid.delta)
            pzoom = 100-posF(perOut-100)
        
            if drawUI {
				calcFrameRW.Lock()
				calcFrame = calcImg
				calcFrameRW.Unlock()
			}
			procFrameQueue4 <- &ImageCropMask{CBox:cropRect,Image:mask.Image, zoom:zoom}
        }
    }
}

// ---------------- Helper Functions ----------------

func calculatePanAdjustment(centroid image.Point, frameWidth int) int {
    // Calculate how much you need to pan to keep the centroid in the center
	panAdj := centroid.X - frameWidth/2

	// smoothedPanAdj := ma.Add(panAdj)
    pid.SetSetpoint(float64(panAdj))
    smoothedPanAdj := pid.Update()
    pid.SetDelta(pid.delta+smoothedPanAdj)

    return int(pid.delta)
}

func calculateCropRect(centroid image.Point, frameWidth, cropWidth, panAdj int) image.Rectangle {
	frameCtr := frameWidth/2
	min := (cropWidth/2)-frameCtr
	max := -min
	panAdj = clamp(panAdj, min, max)

	cropRect := image.Rect((frameCtr+panAdj)-cropWidth/2, 0, (frameCtr+panAdj)+cropWidth/2, cropWidth)

	return cropRect
}

func minBoxHDist(bbox, cropbox image.Rectangle) int {
    leftD := ((bbox.Min.X-cropbox.Min.X))
    rightD := (-(bbox.Max.X-cropbox.Max.X))
    if leftD < rightD {
		return -leftD
	}
	return -rightD
}