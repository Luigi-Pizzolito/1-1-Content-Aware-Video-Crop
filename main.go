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
    playOnlyMode    bool

    bar             *progressbar.ProgressBar
)

func parseFlags() {
    // Define flags
    playOnlyModeF := flag.Bool("play", false, "Player only mode, do not write output")
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
			fmt.Println("Error getting current working directory: "+err.Error())
            selfExit()
            return
		}
	} else {
        // make output folder if it doesn't exist
        if !*playOnlyModeF {
            _, err := os.Stat(outputDirF)
            if os.IsNotExist(err) {
                // Create the directory and its parents if they don't exist
                err := os.MkdirAll(outputDirF, os.ModePerm)
                if err != nil {
                    fmt.Println("Error creating output directory: "+err.Error())
                    selfExit()
                    return
                }
            }
        }
    }
    if !*playOnlyModeF {
        var err error
        outputDirF, err = convertToAbsolutePath(outputDirF)
        if err != nil {
            fmt.Println("Error getting absolute path of output folder: "+err.Error())
            selfExit()
            return
        }
    }

    // Process folder input instead of single video
    if inputVideosF == "" {
        fmt.Println("Error: no input video/directory specified with -i.")
        selfExit()
        return
    }
    files, isFolder, err := getFilesInPath(inputVideosF)
	if err != nil {
		fmt.Println("Error:", err)
		selfExit()
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
    if !*drawUIF {
        // if not drawing UI, cannot be real-time proccesed
        *realTimeF = false
    }
    if *playOnlyModeF {
        // if player only mode, must draw ui and real-time
        *realTimeF = true
        *drawUIF = true
    }


	// Print parsed flags
    fmt.Println("1:1 Content-Aware Video Cropper V1.0")
	fmt.Println("Parameters:")
    if *playOnlyModeF {
        fmt.Println("\tMode: Player-Only")
        fmt.Printf("\tInput Video(s): %v\n", files)
        fmt.Printf("\tSquare Size: %d\n", *squareSizeF)
    } else {
        if *drawUIF {
            fmt.Println("\tMode: Show Algorithm")
            fmt.Printf("\tReal-time: %v\n", *realTimeF)
        } else {
            fmt.Println("\tMode: Proccess")
        }
        // fmt.Printf("\tDraw UI: %v\n", *drawUIF)
        fmt.Printf("\tInput Video(s): %v\n", files)
        fmt.Printf("\tOutput Directory: %s\n", outputDirF)
        fmt.Printf("\tSquare Size: %d\n", *squareSizeF)
    }
	

    playOnlyMode = *playOnlyModeF
    drawUI = *drawUIF
    realTime = *realTimeF
    squareSize = *squareSizeF
    outputDir = outputDirF
    for _, video := range files {
        video, err := convertToAbsolutePath(video)
        if err != nil {
            fmt.Println("Error getting absolute path of input video: "+err.Error())
            selfExit()
            return
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

func selfExit() {
    syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}