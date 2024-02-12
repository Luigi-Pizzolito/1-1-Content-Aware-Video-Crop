package main

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
)

var (
    // input parameters
    squareSize      int
    inputVideo      string
    outputDir       string	
    realTime        bool
    drawUI          bool

    bar             *progressbar.ProgressBar
)

func main() {
    inputVideo = "vid/mukuro.mp4"
    outputDir = "out"
    realTime = false
    squareSize = 256
    drawUI = false

    if !drawUI {
        realTime = false
    }
    
    screenWidth, screenHeight = CalculateWidthHeight(squareSize)
    if drawUI {
        setupUI()
    }

    cleanPathsAndTempDir()
    defer rmTempDirs()

	establishPipesAndImgs()

    fmt.Println("Processing",inputVideo)

    initalResize()

	startPipeline()

	if drawUI {
        runUI()
    } else {
        tickPipeline()
    }
}