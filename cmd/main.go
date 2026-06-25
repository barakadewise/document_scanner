package main

import (
	"fmt"

	cam "scanner.go/internal/camera"
	scan "scanner.go/internal/scanner"
)

func main() {
	println("Initializing  GoDocument scanner....")

	//call camera intilization see if it works
	cam.Webcam()

	imageScaner := scan.ImageScanner{}
	result, err := imageScaner.Scan("/home/dewise/Downloads/test.jpeg")
	if err != nil {
		panic(err)
	}

	fmt.Println("Original file provided ", result.OriginalPath)
	fmt.Println("New file scaned is ", result.ScannedPath)
}
