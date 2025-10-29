package lib

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/hooklift/iso9660"
)

var (
	dirLimit       int
	imgDirs        = map[string][]string{}
	img            string
	exportLocation string
)

func ProcessImage(imagePath string, directoryLimit int, exportLoc string) error {
	img = imagePath
	dirLimit = directoryLimit
	exportLocation = exportLoc

	if err := setup(); err != nil {
		return err
	}

	if err := readImage(img); err != nil {
		return err
	}

	if err := analyzeDirectories(); err != nil {
		return err
	}

	if err := exportDirectories(img); err != nil {
		return err
	}

	return nil
}

func setup() error {
	imageName := filepath.Base(img)
	ext := filepath.Ext(imageName)
	imageName = imageName[0 : len(imageName)-len(ext)]
	exportLocation = filepath.Join(exportLocation, imageName)
	if err := os.MkdirAll(exportLocation, os.ModePerm); err != nil {
		return err
	}
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
				if err := sampleDirectory(dir); err != nil {
					return err
				}
			}
		} else {
			delete(imgDirs, dir)
		}
	}
	return nil
}

func sampleDirectory(dir string) error {
	files := imgDirs[dir]
	selectedFiles := []string{}
	for i := dirLimit; i > 0; i-- {
		j := rand.Intn(len(files)-0) + 0
		selectedFiles = append(selectedFiles, files[j])
		files = append(files[:j], files[j+1:]...)
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
			destDir := filepath.Join(exportLocation, filepath.Clean(f.Name()))
			if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
				return err
			}
		} else {
			if isIncluded(f.Name()) {
				dir, filename := filepath.Split(f.Name())
				df := filepath.Join(exportLocation, filepath.Clean(dir), filename)
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
