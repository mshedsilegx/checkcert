package main

import (
	"os"
	"path/filepath"
	"strings"
)

// getFilesFromDir recursively scans a directory for files matching the specified extension.
// It returns a slice of file paths or an error if the directory cannot be read.
func getFilesFromDir(dir, ext string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ext) {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
