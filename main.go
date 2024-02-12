package main

import (
	"fmt"
    "os"
    "syscall"
    "os/signal"
    "flag"
    "path/filepath"
    
	"github.com/schollz/progressbar/v3"
)

var (
    // input parameters
    squareSize      int
    inputVideo      string
    outputDir       string
    inputVideos     []string	
    realTime        bool
    drawUI          bool

    bar             *progressbar.ProgressBar
)

func parseFlags() {
    // Define flags
	drawUIF := flag.Bool("ui", false, "Whether to draw UI")
	realTimeF := flag.Bool("rt", false, "Whether to run in real-time")
	squareSizeF := flag.Int("s", 256, "Size of square output video")
	var outputDirF string
	flag.StringVar(&outputDirF, "o", "", "Output directory (default is current directory)")
    var inputVideosF string
	flag.StringVar(&inputVideosF, "i", "", "Input video or directory")

    // Parse flags()
    flag.Parse()

    // Set default output directory to current working directory if not provided
	if outputDirF == "" {
		var err error
		outputDirF, err = os.Getwd()
		if err != nil {
			panic("Error getting current working directory: "+err.Error())
		}
	} else {
        // make output folder if it doesn't exist
        _, err := os.Stat(outputDirF)
        if os.IsNotExist(err) {
            // Create the directory and its parents if they don't exist
            err := os.MkdirAll(outputDirF, os.ModePerm)
            if err != nil {
                panic("Error creating output directory: "+err.Error())
            }
        }
    }
    var err error
    outputDirF, err = convertToAbsolutePath(outputDirF)
    if err != nil {
        panic("Error getting absolute path of output folder: "+err.Error())
    }

    // Process folder input instead of single video
    files, isFolder, err := getFilesInPath(inputVideosF)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
    if !isFolder {
        files = []string{}
        files = append(files, inputVideosF)
    }
    for i, file := range files {
        if isFolder {
            files[i] = inputVideosF+string(filepath.Separator)+file
        }
    }

    // Override Logic
    if !drawUI {
        realTime = false
    }


	// Print parsed flags
    fmt.Println("1:1 Content-Aware Video Cropper V1.0")
	fmt.Println("Parameters:")
	fmt.Printf("\tDraw UI: %v\n", *drawUIF)
	fmt.Printf("\tReal-time: %v\n", *realTimeF)
	fmt.Printf("\tSquare Size: %d\n", *squareSizeF)
	fmt.Printf("\tOutput Directory: %s\n", outputDirF)
	fmt.Printf("\tInput Video(s): %v\n", files)

    drawUI = *drawUIF
    realTime = *realTimeF
    squareSize = *squareSizeF
    outputDir = outputDirF
    for _, video := range files {
        video, err := convertToAbsolutePath(video)
        if err != nil {
            panic("Error getting absolute path of input video: "+err.Error())
        }
        inputVideos = append(inputVideos, video)
    }
}

func main() {
    catchExit()

    parseFlags()
    
    screenWidth, screenHeight = CalculateWidthHeight(squareSize)
    if drawUI {
        setupUI()
    }

    cleanPathsAndTempDir()
    defer rmTempDirs()

    if drawUI {
        go runUI()
    }

    var i int
	for i, inputVideo = range inputVideos {
        establishPipesAndImgs()

        fmt.Printf("\nProcessing %d/%d: %s\n", i+1, len(inputVideos), getBasenameWithoutExt(inputVideo))

        initalResize()

        startPipeline()

        if drawUI {
            tickPipeline()
            if i == len(inputVideos)-1 {
                //last one
                termUI = true
            }
        } else {
            tickPipeline()
        }
    }
}

func catchExit() {
    gracefulShutdown := make(chan os.Signal, 1)
    signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <- gracefulShutdown
        fmt.Println("\n\nCaught EXIT")
        termUI = true
        rmTempDirs()
        os.Exit(1)
    }()
}