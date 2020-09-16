package main

import (
	"fmt"
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

func main() {
	// if len(os.Args) < 3 {
	// 	fmt.Println("How to run:\n\tfacedetect [camera ID] [classifier XML file]")
	// 	return
	// }

	// parse args
	deviceID := 0
	xmlFile := "data/face.xml"

	// open webcam
	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()

	// open display window
	window := gocv.NewWindow("Face Detect")
	defer window.Close()

	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	// color for the rect when faces detected
	c := color.RGBA{255, 255, 255, 0}

	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load(xmlFile) {
		fmt.Printf("Error reading cascade file: %v\n", xmlFile)
		return
	}

	fmt.Printf("start reading camera device: %v\n", deviceID)
	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("cannot read device %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// detect faces
		rects := classifier.DetectMultiScale(img)
		// fmt.Printf("found %d faces\n", len(rects))

		// draw a rectangle around each face on the original image,
		// along with text identifying as "Human"
		for _, r := range rects {
			gocv.Rectangle(&img, r, c, 3)

			txt1 := fmt.Sprintf("( %d, %d )", r.Min.X, r.Min.Y)
			txt2 := fmt.Sprintf("( %d, %d )", r.Max.X, r.Max.Y)

			size1 := gocv.GetTextSize(txt1, gocv.FontHersheyPlain, 1.2, 2)
			size2 := gocv.GetTextSize(txt2, gocv.FontHersheyPlain, 1.2, 2)

			pt1 := image.Pt(r.Min.X-(size1.X/2), r.Min.Y-size2.Y)
			pt2 := image.Pt(r.Max.X-(size2.X/2), r.Max.Y+size2.Y+10)

			gocv.PutText(&img, txt1, pt1, gocv.FontHersheyPlain, 1.2, c, 2)
			gocv.PutText(&img, txt2, pt2, gocv.FontHersheyPlain, 1.2, c, 2)
		}

		// show the image in the window, and wait 1 millisecond
		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
