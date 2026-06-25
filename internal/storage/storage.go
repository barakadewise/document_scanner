package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gocv.io/x/gocv"
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
func StoreScannedFile(img gocv.Mat) (string, error) {

	targetDir := "storage/scanned"

	err := os.MkdirAll(targetDir, 0755)

	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf(
		"scan_%d.jpeg",
		time.Now().Unix(),
	)

	fullPath := filepath.Join(targetDir, filename)

	if ok := gocv.IMWrite(fullPath, img); !ok {
		return "", fmt.Errorf(
			"failed writing scanned image",
		)
	}

	return fullPath, nil
}
