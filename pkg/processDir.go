package lib

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func ProcessDirectory(directoryPath string, directoryLimit int, percent int, exportLocation string) error {

	log.Println("directoryPath:", directoryPath)
	log.Println("directoryLimit:", directoryLimit)
	log.Println("percent:", percent)
	log.Println("exportLocation:", exportLocation)

	images, err := os.ReadDir(directoryPath)
	if err != nil {
		return err
	}

	for _, img := range images {
		if !img.IsDir() && filepath.Ext(img.Name()) == ".iso" {
			inputPath := filepath.Join(directoryPath, img.Name())
			fmt.Println("Processing image:", inputPath)

			if err := ProcessImage(inputPath, directoryLimit, percent, exportLocation); err != nil {
				return err
			}

		}
	}

	return nil
}
