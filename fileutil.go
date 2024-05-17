package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/otiai10/copy"
)

func copyFiles(sourcePaths, destPaths []string) error {
	if len(sourcePaths) != len(destPaths) {
		return errors.New("sourcePaths and destPaths slices must have the same length")
	}

	for i, sourcePath := range sourcePaths {
		destPath := destPaths[i]

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			err := os.MkdirAll(destPath, 0755)
			if err != nil {
				return err
			}

			err = copy.Copy(sourcePath, destPath)
			if err != nil {
				return err
			}

		} else {
			sourceFile, err := os.Open(sourcePath)
			if err != nil {
				return err
			}
			defer sourceFile.Close()

			destFile, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer destFile.Close()

			_, err = io.Copy(destFile, sourceFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func moveFiles(sourcePaths, destPaths []string) error {
	if len(sourcePaths) != len(destPaths) {
		return errors.New("sourcePaths and destPaths slices must have the same length")

	}

	for i, sourcePath := range sourcePaths {
		destPath := destPaths[i]

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			err := os.MkdirAll(destPath, 0755)
			if err != nil {
				return err
			}

			err = copy.Copy(sourcePath, destPath)
			if err != nil {
				return err
			}

			err = os.RemoveAll(sourcePath)
			if err != nil {
				return err
			}
		} else {
			sourceFile, err := os.Open(sourcePath)
			if err != nil {
				return err
			}
			defer sourceFile.Close()

			destFile, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer destFile.Close()

			_, err = io.Copy(destFile, sourceFile)
			if err != nil {
				return err
			}

			err = os.Remove(sourcePath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getLastModified(path interface{}) (string, error) {
	var (
		filePath string
		fileInfo os.FileInfo
		err      error
	)

	switch v := path.(type) {
	case string:
		filePath = v
		fileInfo, err = os.Stat(filePath)
	case os.DirEntry:
		fileInfo, err = v.Info()
	default:
		return "", err
	}

	if err != nil {
		return "", err
	}

	createTime := fileInfo.ModTime().Format(time.ANSIC)
	return createTime, nil
}

func Size(size int64) (float64, string) {
	fileSize := float64(size)
	unit := "Bytes"

	if fileSize > 1_000_000_000 {
		return fileSize / (1024 * 1024 * 1024), "GB"
	} else if fileSize > 1_000_000 {
		return fileSize / (1024 * 1024), "MB"
	} else if fileSize > 1_000 {
		return fileSize / 1024, "KB"
	}

	return fileSize, unit
}

func calcSize(path string) (int64, error) {
	var totalSize int64

	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	if !fileInfo.IsDir() {
		return fileInfo.Size(), nil
	}

	var files []os.DirEntry
	files, err = os.ReadDir(path)
	if err != nil {
		return 0, err
	}

	for _, file := range files {
		filePath := filepath.Join(path, file.Name())
		if file.IsDir() {
			size, err := calcSize(filePath)
			if err != nil {
				return 0, err
			}
			totalSize += size
		} else {
			info, err := file.Info()
			if err != nil {
				return 0, err
			}
			totalSize += info.Size()
		}
	}
	return totalSize, nil
}


func displaySingleFileInfo(filepath string) (string, float64, string, string, os.FileMode, error) {
	file, err := os.Stat(filepath)
	if err != nil {
		return "", 0, "", "", 0, err
	}
	var fileType string
	if file.IsDir() {
		fileType = "Directory"
	} else {
		fileType = "File"
	}
	size, err := calcSize(filepath)
		if err != nil {
		return "", 0, "", "", 0, err
	}
	formattedSize, sizeUnit := Size(size)
	creationTime, err := getLastModified(filepath)
	if err != nil {
		return "", 0, "", "", 0, err
	}
	mode := file.Mode().Perm()

	return fileType, formattedSize, sizeUnit, creationTime, mode, nil
}
