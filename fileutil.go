package main

import (
	"io"
	"os"
	"time"

	"github.com/otiai10/copy"
)

func copyFile(sourcePath, destPath string) error {
    fileInfo, err := os.Stat(sourcePath)
    if err != nil {
        return err
    }

    if fileInfo.IsDir() {
        err := copy.Copy(sourcePath, destPath)
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

    return nil
}


func moveFile(src, dst string) error {
    fileInfo, err := os.Stat(src)
    if err != nil {
        return err
    }

    if fileInfo.IsDir() {
        err := copy.Copy(src, dst)
        if err != nil {
            return err
        }
        err = os.Remove(src)
        if err != nil {
            return err
        }
    } else {
        sourceFile, err := os.Open(src)
        if err != nil {
            return err
        }
        defer sourceFile.Close()

        destFile, err := os.Create(dst)
        if err != nil {
            return err
        }
        defer destFile.Close()

        _, err = io.Copy(destFile, sourceFile)
        if err != nil {
            return err
        }

        err = os.Remove(src)
        if err != nil {
            return err
        }
    }

    return nil
}

func getLastModified(path interface{}) (string, error) {
    var (
        filePath string
        fileInfo os.FileInfo
        err error
    )

    switch v := path.(type) {
    case string:
        filePath = v
        fileInfo, err = os.Stat(filePath)
    case os.DirEntry:
        fileInfo, err = v.Info()
        filePath = fileInfo.Name()
    default:
        return "", err
    }

    if err != nil {
        return "", err
    }

    createTime := fileInfo.ModTime().Format(time.ANSIC)
    return createTime, nil
}

func sizeFile(size int64) (float64, string) {
    fileSize := float64(size)
    defaultSize := "Bytes"

    if fileSize > 1000000000 {
        calc := fileSize / (1024 * 1024 * 1024)
        size := "GB"
        return calc, size
    } else if fileSize > 1000000 {
        calc := fileSize / (1024 * 1024)
        size := "MB"
        return calc, size
    } else if fileSize > 1000 {
        calc := fileSize / 1024
        size := "KB"
        return calc, size
    }

    return fileSize, defaultSize
}

func displaySingleFileInfo(filepath string) (string, float64, string, string) {
    file, _ := os.Stat(filepath)
    fileType := ""
    formatted, size := sizeFile(file.Size())
    creationTime, _ := getLastModified(filepath)

    if file.IsDir() {
        fileType = "Directory"
        return fileType, formatted, size, creationTime
    } else {
        fileType = "File"
        return fileType, formatted, size, creationTime
    }
}