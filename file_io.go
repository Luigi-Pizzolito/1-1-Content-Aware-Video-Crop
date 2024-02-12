package main

import (
	"path/filepath"
	"io/ioutil"
	"os"
	"fmt"
)

func convertToAbsolutePath(path string) (string, error) {
	// Check if the path is already absolute
	if filepath.IsAbs(path) {
		return path, nil
	}

	// Convert relative path to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

func getBasenameWithoutExt(absPath string) string {
	// Get the base name of the file
	base := filepath.Base(absPath)

	// Remove the file extension
	ext := filepath.Ext(base)
	basename := base[:len(base)-len(ext)]

	return basename
}

func cleanPathsAndTempDir() {
	var err error
    inputVideo, err = convertToAbsolutePath(inputVideo)
    if err != nil {
        panic("Couldn't find absolute path of input video")
    }
    

    outputDir, err = convertToAbsolutePath(outputDir)
    if err != nil {
        panic("Couldn't find absolute path of output directory")
    }

	// Create a temporary directory for frames
    tempDirFrames, err = ioutil.TempDir("", "tinyframes")
    if err != nil {
        panic(err)
    }

	// Create a temporary directory for tiny video
    tempDirTinyVid, err = ioutil.TempDir("", "tinyvideo")
    if err != nil {
        panic(err)
    }
}

func rmTempDirs() {
	fmt.Println("Cleaned up.")
	os.RemoveAll(tempDirFrames)
    os.RemoveAll(tempDirTinyVid)
}