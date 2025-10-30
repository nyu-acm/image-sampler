package lib

import (
	"fmt"
	"os"
	"path/filepath"
)

func ProcessDirectory(directoryPath string, directoryLimit int, exportLocation string) error {
	// Implementation for processing a directory of ISO images

	images, err := os.ReadDir(directoryPath)
	if err != nil {
		return err
	}

	for _, img := range images {
		if !img.IsDir() && filepath.Ext(img.Name()) == ".iso" {
			inputPath := filepath.Join(directoryPath, img.Name())
			fmt.Println("Processing image:", inputPath)

			if err := ProcessImage(inputPath, directoryLimit, exportLocation); err != nil {
				return err
			}

		}
	}

	return nil
}
