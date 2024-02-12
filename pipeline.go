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

	//! profiling depps
	"fmt"
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

//! Adding capacity/bufffer to channels
//! Idle before:
//! -- readVideo: 7.7%
//! -- grayVideo: 4.4%
//! -- binaVideo: 0.3%
//! -- calcVideo: 7.9%
//! -- cropVideo: 1.4%
//! --> total:	  21.7% pipeline Idle
//! Idle after:
//! -- readVideo: 7.8%
//! -- grayVideo: 3.9%
//! -- binaVideo: 0.0%
//! -- calcVideo: 7.6%
//! -- cropVideo: 0.4%
//! --> total:	  19.7% pipeline Idle
//! Idle after:
//! -- readVideo: 7.8%
//! -- grayVideo: 4.0%
//! -- binaVideo: 0.0%
//! -- calcVideo: 7.7%
//! -- cropVideo: 0.5%
//! --> total:	  20% pipeline Idle

//! channel capacity monitoring for profiling
func monitorChan[T any](ch chan T, name string) {
	chanMonitorInterval := time.Second
    for {
        // if len(ch) == cap(ch) {
            fmt.Printf("Channel: %s, Size: %d\n", name, len(ch))
        // }
        time.Sleep(chanMonitorInterval)
    }
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

	//! channel capacity monitoring for profiling
	// go monitorChan(procFrameQueue1, "read->gray")
	// go monitorChan(procFrameQueue2, "gray->bin")
	// go monitorChan(procFrameQueue3, "bin->calc")
	// go monitorChan(procFrameQueue4, "calc->crop")
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