package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Fetcher struct {
	Destination string
}

// NewFetcher creates a new scriplets fetcher
func NewFetcher() (*Fetcher, error) {
	tempDir, err := os.MkdirTemp("", "scriplets-*")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		return nil, err
	}

	return &Fetcher{
		Destination: tempDir,
	}, nil
}

// CopyScriplets from args to a temp directory on the filesystem
// TODO: currently this only works with scriplets on the file system
func (f *Fetcher) CopyScriplets(paths []string) error {
	// Copy each file to the temp directory
	for _, filePath := range paths {
		err := copyFileToDir(filePath, f.Destination)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyFileToDir(srcPath, destDir string) error {
	// Open the source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create the destination file path
	destPath := filepath.Join(destDir, filepath.Base(srcPath))

	// Create the destination file
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy the contents
	_, err = io.Copy(destFile, srcFile)
	return err
}
