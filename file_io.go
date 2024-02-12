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

func getFilesInPath(path string) ([]string, bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, false, err
	}

	if fileInfo.IsDir() {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return nil, true, err
		}

		var fileNames []string
		for _, file := range files {
			fileNames = append(fileNames, file.Name())
		}

		return fileNames, true, nil
	}

	return nil, false, nil
}

func deleteFilesInFolder(folderPath string) error {
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			err := os.Remove(path)
			if err != nil {
				return err
			}
			// fmt.Printf("Deleted file: %s\n", path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}