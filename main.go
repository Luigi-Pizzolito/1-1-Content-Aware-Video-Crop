package main

import (
	"fmt"
    "os"
    "syscall"
    "os/signal"
    // "context"
    // "time"
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
    catchExit()

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

func catchExit() {
    gracefulShutdown := make(chan os.Signal, 1)
    signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <- gracefulShutdown
        fmt.Println("\nCaught EXIT")
        rmTempDirs()
        os.Exit(1)
    }()
}