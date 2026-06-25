package helpers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// check the file mimetype
func CheckFileMimeType(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)

	if err != nil {
		return "", err

	}

	// use http detector to determine the mime and return the tyep
	mime := http.DetectContentType(buffer)
	return mime, nil

}

// get the file path  selcted
func GetFilePath(path string) (string, error) {
	abspath, err := filepath.Abs(path)
	if err != nil {
		return "", err

	}
	_, err = os.Stat(abspath)
	if err != nil {
		return "", fmt.Errorf("File not found %s %w", abspath, err)
	}
	return abspath, nil
}

// Store the file
func StoreFile(sourceDir string, targertDir string) (string, error) {

	err := os.MkdirAll(
		targertDir,
		0755,
	)

	if err != nil {
		return "", err
	}

	//rename the file with its oringinal extension
	filename := fmt.Sprintf(
		"scan_%d%s",
		time.Now().Unix(),
		filepath.Ext(sourceDir),
	)

	//join file path
	targertDir = filepath.Join(
		targertDir,
		filename,
	)

	//read file data
	data, err := os.ReadFile(sourceDir)

	if err != nil {
		return "", err
	}

	//wirte file to storage
	err = os.WriteFile(
		targertDir,
		data,
		0644,
	)

	if err != nil {
		return "", err
	}

	return targertDir, nil
}
