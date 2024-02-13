package main

import (
	"image"
	"image/color"
	"image/draw"
	"sync"
	"time"
	"path/filepath"
	"github.com/schollz/progressbar/v3"
	"github.com/AlexEidt/Vidio"
)

var (
	inputFrame 		*image.RGBA
	inputFrameRW 	sync.Mutex

	procFrameQueue1	chan *image.RGBA

	grayFrame		*image.RGBA
	grayFrameRW		sync.Mutex

	procFrameQueue2	chan *ImageGrayMask

	binaryFrame		*image.RGBA
	binaryFrameRW	sync.Mutex

	procFrameQueue3	chan *ImageGrayBoundMask

	calcFrame		*image.RGBA
	calcFrameRW		sync.Mutex

	procFrameQueue4	chan *ImageCropMask

	croppedFrame	*image.RGBA
	croppedFrameRW	sync.Mutex


	tempDirTinyVid  string
    tempDirFrames   string
)

type ImageGrayMask struct {
    Mask	*image.Gray
    Image	*image.RGBA
}

type ImageGrayBoundMask struct {
    Mask    *image.Gray
    Image   *image.RGBA
    Bbox    image.Rectangle
    Subj    int
}

type ImageCropMask struct {
	CBox	image.Rectangle
	Image	*image.RGBA
    zoom    bool
}

func establishPipesAndImgs() {
	inputFrame = image.NewRGBA(image.Rect(0, 0, screenWidth, squareSize))
	draw.Draw(inputFrame, inputFrame.Bounds(), &image.Uniform{color.RGBA{255,0,255,255}}, image.ZP, draw.Src)
	
	procFrameQueue1 = make(chan *image.RGBA, 1000)

	grayFrame = image.NewRGBA(image.Rect(0, 0, screenWidth, squareSize))
	draw.Draw(grayFrame, grayFrame.Bounds(), &image.Uniform{color.RGBA{255,0,0,255}}, image.ZP, draw.Src)
	
	procFrameQueue2 = make(chan *ImageGrayMask, 500)

	binaryFrame = image.NewRGBA(image.Rect(0, 0, screenWidth, squareSize))
	draw.Draw(binaryFrame, binaryFrame.Bounds(), &image.Uniform{color.RGBA{0,255,0,255}}, image.ZP, draw.Src)

	procFrameQueue3 = make(chan *ImageGrayBoundMask, 4)

	calcFrame = image.NewRGBA(image.Rect(0, 0, screenWidth, squareSize))
	draw.Draw(calcFrame, calcFrame.Bounds(), &image.Uniform{color.RGBA{0,0,255,255}}, image.ZP, draw.Src)

	procFrameQueue4 = make(chan *ImageCropMask, 100)

	croppedFrame = image.NewRGBA(image.Rect(0, 0, squareSize, squareSize))
	draw.Draw(croppedFrame, croppedFrame.Bounds(), &image.Uniform{color.RGBA{0,255,255,255}}, image.ZP, draw.Src)
}

func startPipeline() {
	video, _ := vidio.NewVideo(filepath.Join(tempDirTinyVid,getBasenameWithoutExt(inputVideo)+".mp4"))
    tframes = video.Frames()
    cframes = 0
    if !playOnlyMode {
		bar = progressbar.Default(int64(tframes))
	}
	video.Close()

	go readVideo()
	go grayscaleVideo()
	go binaryVideo()
	go calcVideo()
	go cropVideo()
}

func checkPipelineDone() bool {
    if cframes == tframes {
		if !playOnlyMode {
			bar.Clear()
			bar.Close()
		}
		close(procFrameQueue1)
		close(procFrameQueue2)
		close(procFrameQueue3)
		close(procFrameQueue4)
        if !playOnlyMode {
			// Last step to make video file from image sequence
			joinOutputImgs()
		}
        return true
    }
    return false
}

func tickPipeline() {
    ticker := time.NewTicker(time.Millisecond * 16)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            if checkPipelineDone() {
                return
            }
        }
    }
}