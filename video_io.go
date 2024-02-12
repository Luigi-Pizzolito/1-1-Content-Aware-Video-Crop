package main

import (
	"image"
	"image/draw"
	"image/jpeg"
	"path/filepath"
	"strings"
	"strconv"
	"io/ioutil"
	"bufio"
	"github.com/schollz/progressbar/v3"
	"os"
	"os/exec"
	"fmt"
)

// CalculateWidthHeight calculates the width of a 16:9 video given the height in pixels,
// ensuring both width and height are divisible by 2.
func CalculateWidthHeight(height int) (width, newHeight int) {
	// Calculate width based on the 16:9 aspect ratio formula: width = (height * 16) / 9
	width = (height * 16) / 9

	// Ensure width is divisible by 2
	if width%2 != 0 {
		width++
	}

	// Ensure height is divisible by 2
	if height%2 != 0 {
		height++
	}

	return width, height
}

// copyAndWrite copies the image.RGBA and passes it to a video writer.
//! takes 5.2% of run time
//! does not block requests if pipe would have a buffer, hm
func copyAndWrite(img *image.RGBA) {
	// Copy the image.RGBA to a new image.RGBA
	copiedImg := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	draw.Draw(copiedImg, copiedImg.Bounds(), img, img.Bounds().Min, draw.Src)

	// Write the frame to the image
    abspath := filepath.Join(tempDirFrames, fmt.Sprintf("%d.jpg", cframes))

    f, _ := os.Create(abspath)
	jpeg.Encode(f, copiedImg, nil)
	f.Close()

    cframes++
    bar.Add(1)
}

func joinOutputImgs() {
    abspath, err := filepath.Abs(tempDirFrames)
    if err != nil {
		fmt.Println("Error:", err)
	}
    err = convertImagesToMP4(abspath, inputVideo, filepath.Join(outputDir,getBasenameWithoutExt(inputVideo)+".mp4"), fps)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func convertImagesToMP4(inputFolder, audioFile, outputFile string, fps float64) error {
	cmd := exec.Command("ffmpeg",
		"-r", strconv.FormatFloat(fps, 'f', -1, 64),
        "-framerate", strconv.FormatFloat(fps, 'f', -1, 64),
		// "-pattern_type", "glob",
		"-i", inputFolder+"/%d.jpg",
		"-i", audioFile,
        "-map", "0:v:0",
        "-map", "1:a:0",
        // "-c:v", "copy",
		"-c:v", "mjpeg",
		"-q:v", "6",
        // "-vf", "setsar=1,crop=aspect_ratio=1",//"cropdetect=limit=0:round=2:reset=0,scale="+strconv.Itoa(squareSize)+":"+strconv.Itoa(squareSize),
        // "-aspect", "1:1",
		"-c:a", "copy",
		outputFile,"-y", // Overwrite output file without asking
	)
    fmt.Println(cmd.String())

	// Create a progress bar
	files, err := ioutil.ReadDir(inputFolder)
    if err != nil {
		return err
	}
	totalFrames := len(files)
	barp := progressbar.Default(int64(totalFrames))

	// Create a progress bar
	cmdReader, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

    // Start the command
	err = cmd.Start()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(cmdReader)
    scanner.Split(bufio.ScanWords)
    getF := false
    for scanner.Scan() {
        line := scanner.Text()
        // fmt.Println(line)
        if getF {
            getF = false
            cf, err := strconv.Atoi(line)
            if err != nil {
                continue
            }
            barp.Set(cf)
        }
        if strings.Contains(line, "frame=") {
            getF = true
        }
    }

	// Wait for the command to finish
	err = cmd.Wait()
	if err != nil {
		return err
	}

	barp.Clear()
    barp.Close()

	//! Must clear output folder image sequence after, or it might get appended to other videos being processed
	err = deleteFilesInFolder(inputFolder)
	if err != nil {
		return err
	}

	return nil
}

func resizeTiny(inputFile, outputFile string, width, height int) error {
	cmd := exec.Command("ffmpeg",
		"-i", inputFile,
		"-vf", "scale="+strconv.Itoa(width)+":"+strconv.Itoa(height),
		outputFile,
		"-y", // Overwrite output file without asking
	)

    fmt.Println(cmd.String())

	// Create a progress bar
	cmdReader, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

    // Start the command
	err = cmd.Start()
	if err != nil {
		return err
	}

    totalFrames, err := countFrames(inputFile)
    if err != nil {
		return err
	}

	scanner := bufio.NewScanner(cmdReader)
    scanner.Split(bufio.ScanWords)
    getF := false
    barp := progressbar.Default(int64(totalFrames))
    for scanner.Scan() {
        line := scanner.Text()
        // fmt.Println(line)
        if getF {
            getF = false
            cf, err := strconv.Atoi(line)
            if err != nil {
                continue
            }
            barp.Set(cf)
        }
        if strings.Contains(line, "frame=") {
            getF = true
        }
    }

	// Wait for the command to finish
	err = cmd.Wait()
	if err != nil {
		return err
	}

	barp.Clear()
    barp.Close()
	return nil
}