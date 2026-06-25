package scanner

import (
	"fmt"

	"gocv.io/x/gocv"
	"scanner.go/helpers"
	"scanner.go/internal/storage"
)

// scan results
type ScanResults struct {
	OriginalPath string
	ScannedPath  string
}

// scanner interface
type Scanner interface {
	Scan(path string) (*ScanResults, error)
}

type ImageScanner struct {
}

type DocumentScanner struct{}

// imegae scanner implementation
func (s *ImageScanner) Scan(path string) (*ScanResults, error) {
	file, err := storage.LoadFile(path)
	if err != nil {
		return nil, err

	}

	mime, err := helpers.CheckFileMimeType(file)
	if err != nil {
		return nil, err
	}

	if mime != "image/jpeg" && mime != "image/png" {
		return nil, fmt.Errorf(
			"unsupported file type: %s",
			mime,
		)
	}

	//appy scan
	scannedFile, err := s.ApplyScan(file)

	if err != nil {
		return nil, err
	}

	newFilepath, err := storage.StoreScannedFile(scannedFile)
	if err != nil {
		return nil, err
	}

	return &ScanResults{
		OriginalPath: path,
		ScannedPath:  newFilepath,
	}, nil

}

func (s ImageScanner) ApplyScan(path string) (string, error) {
	//read image with gocv
	img := gocv.IMRead(path, gocv.IMReadColor)

	if img.Empty() {
		return "", fmt.Errorf("Failed to read image")
	}
	defer img.Close()

	//convert image to grascale firt

	return path, nil
}
