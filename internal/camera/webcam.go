package camera

import (
	"fmt"
	"log"

	"gocv.io/x/gocv"
)

func Webcam() {
	println("Camera initilization....")
	webcam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		log.Fatalf("cannot open camera: %v", err)
	}

	defer webcam.Close()

	if !webcam.IsOpened() {
		log.Fatal("camera failed to open")
	}

	fmt.Println("Camera opened successfully")

	img := gocv.NewMat()
	defer img.Close()

	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Println("cannot read camera")
			break
		}

		if img.Empty() {
			continue
		}

		fmt.Println("Frame received:", img.Cols(), "x", img.Rows())
		break
	}
}
