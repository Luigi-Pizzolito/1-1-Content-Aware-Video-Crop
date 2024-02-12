package main

import (
	"github.com/AlexEidt/Vidio"
	"github.com/AlexEidt/aio"
	"image"
	"path/filepath"
	"time"
	"fmt"
)

var (
    fps             float64
)

// ---------------- Step Function ----------------

func initalResize() {
	fmt.Printf("Initial resize to %dx%dpx\n", screenWidth, squareSize)
    err := resizeTiny(inputVideo, filepath.Join(tempDirTinyVid,getBasenameWithoutExt(inputVideo)+".mp4"), screenWidth, squareSize)
    if err != nil {
        fmt.Println("Error:", err)
    }
    fmt.Println("Resize completed, Running algorithm...")
}

func readVideo() {
	video, _ := vidio.NewVideo(filepath.Join(tempDirTinyVid,getBasenameWithoutExt(inputVideo)+".mp4"))

    vbuffer := image.NewRGBA(image.Rect(0, 0, video.Width(), video.Height()))
	video.SetFrameBuffer(vbuffer.Pix)

	audio, _ := aio.NewAudio(filepath.Join(tempDirTinyVid,getBasenameWithoutExt(inputVideo)+".mp4"), nil)
	player, _ := aio.NewPlayer(audio.Channels(), audio.SampleRate(), audio.Format())
	defer player.Close()

	fps = video.FPS()
	maxFrameDuration := time.Second / time.Duration(fps)

	if realTime {
		go func(){
			for audio.Read() {
				player.Play(audio.Buffer())
			}
		}()
	} else {
		player.Close()
	}

    for video.Read() {
        startTime := time.Now()

		if drawUI {
			inputFrameRW.Lock()
			inputFrame = vbuffer
			inputFrameRW.Unlock()
		}
		procFrameQueue1 <- vbuffer

		elapsedTime := time.Since(startTime)
		remainingTime := maxFrameDuration - elapsedTime
		if remainingTime > 0 && realTime {
            // 2% speed up
			time.Sleep(time.Duration(int64(float64(remainingTime) * 0.98)))
		}

    }

    video.Close()
}

// ---------------- Helper Functions ----------------

func countFrames(videoFile string) (int, error) {
	// Open the video file
	video, err := vidio.NewVideo(videoFile)
	if err != nil {
		return 0, err
	}
	defer video.Close()

	return video.Frames(), nil
}