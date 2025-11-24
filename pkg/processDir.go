package lib

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
)

var exportRemoved bool

func ProcessDirectory(directoryPath string, directoryLimit int, percent int, exportLocation string, removedLocation string, xportRemoved bool) error {

	log.Printf("[info] directoryPath: %s, directoryLimit: %d, percent: %d, exportLocation: %s, removedLocation: %s\n", directoryPath, directoryLimit, percent, exportLocation, removedLocation)
	images, err := os.ReadDir(directoryPath)
	if err != nil {
		return err
	}

	exportRemoved = xportRemoved

	removedFilesOut, err = os.Create("removedFiles.txt")
	if err != nil {
		return err
	}
	defer removedFilesOut.Close()
	removedWriter = bufio.NewWriter(removedFilesOut)
	defer removedWriter.Flush()

	for _, img := range images {
		if !img.IsDir() && filepath.Ext(img.Name()) == ".iso" {
			inputPath := filepath.Join(directoryPath, img.Name())
			log.Println("[info] processing image:", inputPath)

			if err := ProcessImage(inputPath, directoryLimit, percent, exportLocation, removedLocation); err != nil {
				return err
			}

		}
	}

	return nil
}
