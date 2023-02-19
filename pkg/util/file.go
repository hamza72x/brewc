package util

import (
	"io"
	"os"
)

// DoesFileExist returns true if the file exists
// and is not a directory.
func DoesFileExist(path string) bool {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

// DoesDirExist returns true if the directory exists.
func DoesDirExist(path string) bool {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

// CreateDirIfNotExists creates the given directory if it does not exist.
func CreateDirIfNotExists(path string) error {
	if DoesDirExist(path) {
		return nil
	}

	return os.MkdirAll(path, os.ModePerm)
}

// WriteFile writes the given io.Reader to the given path.
// and it overwrites the file if it already exists.
func WriteFile(path string, reader io.Reader) error {
	var file, err = os.Create(path)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, reader)

	if err != nil {
		return err
	}

	return nil
}
