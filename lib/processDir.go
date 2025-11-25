package lib

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nyudlts/bytemath"
)

var (
	removedWriter         *csv.Writer
	removedFilesOut       *os.File
	retainedWriter        *csv.Writer
	retainedFilesOut      *os.File
	dirLimit              int
	percentage            int
	exportRemoved         bool
	totalSizeRemoved      int64 = 0
	totalNumFilesRemoved  int   = 0
	totalSizeRetained     int64 = 0
	totalNumFilesRetained int   = 0
)

func ProcessDirectory(directoryPath string, directoryLimit int, percent int, exportLocation string, removedLocation string, xportRemoved bool) error {
	dirLimit = directoryLimit
	percentage = percent
	exportRemoved = xportRemoved
	//print some info
	log.Printf("[info] settings: directoryPath: %s, directoryLimit: %d, percent: %d, exportLocation: %s, removedLocation: %s, exportRemoved: %t\n", directoryPath, directoryLimit, percent, exportLocation, removedLocation, xportRemoved)

	exportRemoved = xportRemoved

	//create output file for removed files
	var err error
	removedFilesOut, err = os.Create("removedFiles.csv")
	if err != nil {
		return err
	}
	defer removedFilesOut.Close()
	removedWriter = csv.NewWriter(removedFilesOut)
	defer removedWriter.Flush()
	removedWriter.Write([]string{"image", "path", "filename", "extension", "size", "size human"})

	//create output file for retained files
	retainedFilesOut, err = os.Create("retainedFiles.csv")
	if err != nil {
		return err
	}
	defer retainedFilesOut.Close()
	retainedWriter = csv.NewWriter(retainedFilesOut)
	defer retainedWriter.Flush()
	retainedWriter.Write([]string{"image", "path", "filename", "extension", "size", "size human"})

	var ISOs = []string{}
	if err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(info.Name()) == ".iso" {
			ISOs = append(ISOs, path)
			log.Println("[info] processing image:", path)
			if err := ProcessImage(path, exportLocation, removedLocation); err != nil {
				log.Printf("[error] processing image %s: %v\n", path, err)
			}
		}
		return nil
	}); err != nil {
		return err
	}

	totalSizeRemovedHuman := bytemath.ConvertBytesToHumanReadable(totalSizeRemoved)
	totalHuman := fmt.Sprintf("%d (%s)", totalNumFilesRemoved, totalSizeRemovedHuman)
	totalsMsg := fmt.Sprintf("total files removed: %s", totalHuman)
	removedWriter.Write([]string{"", "", "", "totals:", totalHuman, ""})
	log.Printf("[info] %s", totalsMsg)

	totalSizeRetainedHuman := bytemath.ConvertBytesToHumanReadable(totalSizeRetained)
	totalHuman = fmt.Sprintf("%d (%s)", totalNumFilesRetained, totalSizeRetainedHuman)
	totalsMsg = fmt.Sprintf("total files retained: %s", totalHuman)
	retainedWriter.Write([]string{"", "", "", "totals:", totalHuman, ""})
	log.Printf("[info] %s", totalsMsg)

	return nil
}
