package lib

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/hooklift/iso9660"
	"github.com/nyudlts/bytemath"
)

var (
	imageExtensions  = []string{".jpg", ".jpeg", ".png", ".gif", ".tiff", ".tif", ".cr2", ".cr3", ".nef", ".arw", ".dng", ".orf", ".pef", ".rw2", ".3fr"}
	imgDirs          map[string][]string
	imageFilename    string
	numFilesSelected int
	numFilesRemoved  int
	exportLoc        string
	removedLoc       string
)

func isImageFile(ext string) bool {
	for _, imageExt := range imageExtensions {
		if ext == imageExt {
			return true
		}
	}
	return false
}

func ProcessImage(imagePath string, exportLocation string, removedLocation string) error {
	imageFilename = filepath.Base(imagePath)
	exportLoc = exportLocation
	removedLoc = removedLocation

	if err := setup(imagePath); err != nil {
		return err
	}

	if err := readImage(imagePath); err != nil {
		return err
	}

	if err := analyzeDirectories(); err != nil {
		return err
	}

	if err := exportDirectories(imagePath); err != nil {
		return err
	}

	log.Printf("[info] processing complete on %s: %d files selected, %d files removed.\n", filepath.Base(imagePath), numFilesSelected, numFilesRemoved)

	return nil
}

func setup(imagePath string) error {
	imageName := filepath.Base(imagePath)
	ext := filepath.Ext(imageName)
	imageName = imageName[0 : len(imageName)-len(ext)]
	exportLoc = filepath.Join(exportLoc, imageName)
	removedLoc = filepath.Join(removedLoc, imageName)
	log.Println("[info] creating export directory at:", exportLoc)
	if err := os.MkdirAll(exportLoc, os.ModePerm); err != nil {
		return err
	}

	if exportRemoved {
		log.Println("[info] creating removed directory at:", removedLoc)
		if err := os.MkdirAll(removedLoc, os.ModePerm); err != nil {
			return err
		}
	}

	imgDirs = map[string][]string{}
	return nil
}

func readImage(imagePath string) error {

	img, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	iso, err := iso9660.NewReader(img)
	if err != nil {
		return err
	}

	for {
		f, err := iso.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("failed to read next file: %v\n", err)
			return err
		}

		if f.IsDir() {
			imgDirs[f.Name()] = []string{}
		} else {
			ext := filepath.Ext(f.Name())
			if ext == ".db" || ext == ".df" || f.Name() == ".dss" || f.Name() == ".DS_Store" {
				continue
			}
			path, name := filepath.Split(f.Name())
			imgDirs[path] = append(imgDirs[path], name)
		}

	}
	return nil

}

func analyzeDirectories() error {
	for dir, files := range imgDirs {
		if len(files) > 0 {
			if len(files) >= dirLimit {

				imgCount, err := numImageFiles(files)
				if err != nil {
					return err
				}

				if imgCount >= dirLimit {

					if err := sampleDirectory(dir); err != nil {
						return err
					}
				}
			}
		} else {
			delete(imgDirs, dir)
		}
	}
	return nil
}

func numImageFiles(files []string) (int, error) {
	imageFileCount := 0
	for _, file := range files {
		ext := filepath.Ext(file)
		if isImageFile(ext) {
			imageFileCount++
		}
	}
	return imageFileCount, nil
}

func getImageFiles(files []string) ([]string, []string) {
	imageFiles := []string{}
	nonImageFiles := []string{}
	for _, file := range files {
		ext := filepath.Ext(file)
		if isImageFile(ext) {
			imageFiles = append(imageFiles, file)
		} else {
			nonImageFiles = append(nonImageFiles, file)
		}
	}
	return imageFiles, nonImageFiles
}

func sampleDirectory(dir string) error {

	files := imgDirs[dir]
	imageFiles, nonImageFiles := getImageFiles(files)
	numImageFiles := len(imageFiles)
	sampleSize := (numImageFiles * percentage) / 100

	log.Printf("[info] sampling %d out of %d image files from directory: %s\n", sampleSize, numImageFiles, filepath.Join(imageFilename, dir))

	selectedFiles := []string{}
	selectedFiles = append(selectedFiles, nonImageFiles...)

	for i := sampleSize; i > 0; i-- {
		j := rand.Intn(len(imageFiles))
		selectedFiles = append(selectedFiles, imageFiles[j])
		imageFiles = append(imageFiles[:j], imageFiles[j+1:]...)
	}

	imgDirs[dir] = selectedFiles
	return nil
}

func exportDirectories(imagePath string) error {
	img, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	defer img.Close()
	iso, err := iso9660.NewReader(img)
	if err != nil {
		return err
	}

	for {
		f, err := iso.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("failed to read next file: %v\n", err)
			return err
		}

		if f.IsDir() {
			exportDir := filepath.Join(exportLoc, filepath.Clean(f.Name()))
			if err := os.MkdirAll(exportDir, os.ModePerm); err != nil {
				return err
			}

			if exportRemoved {
				removedDir := filepath.Join(removedLoc, filepath.Clean(f.Name()))
				if err := os.MkdirAll(removedDir, os.ModePerm); err != nil {
					return err
				}
			}
		} else {

			if isIncluded(f.Name()) {
				dir, filename := filepath.Split(f.Name())
				df := filepath.Join(exportLoc, filepath.Clean(dir), filename)
				fReader := f.Sys().(io.Reader)
				ff, err := os.OpenFile(df, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
				if err != nil {
					return err
				}
				defer ff.Close()
				if _, err := io.Copy(ff, fReader); err != nil {
					return err
				}
				numFilesSelected++
				totalNumFilesRetained++
				normalizedPath, fName := filepath.Split(strings.ReplaceAll(f.Name(), "\\", "/"))
				ext := filepath.Ext(fName)
				size := fmt.Sprintf("%d", f.Size())
				humanSize := bytemath.ConvertBytesToHumanReadable(f.Size())
				retainedWriter.Write([]string{imageFilename, normalizedPath, fName, ext, size, humanSize})
				totalSizeRetained += f.Size()

			} else {
				if exportRemoved {
					dir, filename := filepath.Split(f.Name())
					df := filepath.Join(removedLoc, filepath.Clean(dir), filename)
					fReader := f.Sys().(io.Reader)
					ff, err := os.OpenFile(df, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
					if err != nil {
						return err
					}
					defer ff.Close()
					if _, err := io.Copy(ff, fReader); err != nil {
						return err
					}
				}
				numFilesRemoved++
				normalizedPath, fName := filepath.Split(strings.ReplaceAll(f.Name(), "\\", "/"))
				ext := filepath.Ext(fName)
				size := fmt.Sprintf("%d", f.Size())
				humanSize := bytemath.ConvertBytesToHumanReadable(f.Size())
				removedWriter.Write([]string{imageFilename, normalizedPath, fName, ext, size, humanSize})
				totalSizeRemoved += f.Size()
				totalNumFilesRemoved++
			}
		}

	}
	return nil
}

func isIncluded(filename string) bool {
	dir, name := filepath.Split(filename)
	files, exists := imgDirs[dir]
	if !exists {
		return false
	}
	for _, f := range files {
		if f == name {
			return true
		}
	}
	return false
}
