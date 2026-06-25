package storage

import (
	"scanner.go/helpers"
)

func LoadFile(path string) (string, error) {
	data, err := helpers.GetFilePath(path)
	if err != nil {
		return "", err
	}
	return data, nil

}

// store the scanned copy of the file
func StoreScannedFile(sourceDir string) (string, error) {
	targetdir := "storage/scanned"

	destinaltion, err := helpers.StoreFile(sourceDir, targetdir)
	if err != nil {
		return "", err
	}
	return destinaltion, nil

}
