package scanner

import (
	"fmt"
	"image"
	"image/color"

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
	defer scannedFile.Close()

	newFilepath, err := storage.StoreScannedFile(scannedFile)

	if err != nil {
		return nil, err
	}

	return &ScanResults{
		OriginalPath: path,
		ScannedPath:  newFilepath,
	}, nil

}

func (s ImageScanner) ApplyScan(path string) (gocv.Mat, error) {

	img := gocv.IMRead(
		path,
		gocv.IMReadColor,
	)

	if img.Empty() {
		return gocv.Mat{}, fmt.Errorf("failed to read image")
	}

	defer img.Close()

	// grayscale
	gray := gocv.NewMat()
	defer gray.Close()

	gocv.CvtColor(
		img,
		&gray,
		gocv.ColorBGRToGray,
	)

	// blur
	blur := gocv.NewMat()
	defer blur.Close()

	gocv.GaussianBlur(
		gray,
		&blur,
		image.Pt(5, 5),
		0,
		0,
		gocv.BorderDefault,
	)

	// edges
	edges := gocv.NewMat()
	defer edges.Close()

	gocv.Canny(
		blur,
		&edges,
		75,
		200,
	)
	kernel := gocv.GetStructuringElement(
		gocv.MorphRect,
		image.Pt(5, 5),
	)

	defer kernel.Close()

	closed := gocv.NewMat()
	defer closed.Close()

	gocv.MorphologyEx(
		edges,
		&closed,
		gocv.MorphClose,
		kernel,
	)

	edges.Close()

	edges = closed.Clone()

	// find contours
	contours := gocv.FindContours(
		edges,
		gocv.RetrievalExternal,
		gocv.ChainApproxSimple,
	)
	debug := img.Clone()
	defer debug.Close()

	gocv.DrawContours(
		&debug,
		contours,
		-1,
		color.RGBA{255, 0, 0, 255},
		2,
	)

	var document []image.Point
	maxArea := float64(0)

	for i := 0; i < contours.Size(); i++ {

		contour := contours.At(i)

		area := gocv.ContourArea(contour)

		if area < 5000 {
			continue
		}

		perimeter := gocv.ArcLength(
			contour,
			true,
		)

		approx := gocv.ApproxPolyDP(
			contour,
			0.02*perimeter,
			true,
		)

		if approx.Size() == 4 {

			if area > maxArea {

				maxArea = area

				document = approx.ToPoints()
			}
		}

		approx.Close()
	}

	// no document detected
	if len(document) == 0 {
		return img.Clone(), nil
	}

	ordered := orderPoints(document)

	// source points
	src := gocv.NewPointVector()
	defer src.Close()

	for _, p := range ordered {

		src.Append(
			p,
		)
	}

	width := 800
	height := 1000

	dst := gocv.NewPointVector()
	defer dst.Close()

	dst.Append(image.Point{
		X: 0,
		Y: 0,
	})

	dst.Append(image.Point{
		X: width,
		Y: 0,
	})

	dst.Append(image.Point{
		X: width,
		Y: height,
	})

	dst.Append(image.Point{
		X: 0,
		Y: height,
	})

	matrix := gocv.GetPerspectiveTransform(
		src,
		dst,
	)

	defer matrix.Close()

	scanned := gocv.NewMat()
	defer scanned.Close()

	gocv.WarpPerspective(
		img,
		&scanned,
		matrix,
		image.Pt(width, height),
	)

	return scanned.Clone(), nil
}

func orderPoints(points []image.Point) []image.Point {

	ordered := make([]image.Point, 4)

	// TL + BR using sum
	minSum := points[0].X + points[0].Y
	maxSum := minSum

	for _, p := range points {

		sum := p.X + p.Y

		if sum < minSum {
			minSum = sum
			ordered[0] = p
		}

		if sum > maxSum {
			maxSum = sum
			ordered[2] = p
		}
	}

	// TR + BL using difference
	minDiff := points[0].X - points[0].Y
	maxDiff := minDiff

	for _, p := range points {

		diff := p.X - p.Y

		if diff < minDiff {
			minDiff = diff
			ordered[1] = p
		}

		if diff > maxDiff {
			maxDiff = diff
			ordered[3] = p
		}
	}

	return ordered
}
