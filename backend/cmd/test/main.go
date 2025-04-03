package main

import (
	"fmt"
	"image"
	"os/exec"

	"gocv.io/x/gocv"
)

func main() {
	img := gocv.IMRead("/Users/sharithgodamanna/Desktop/Code/recept-scanner/data/large-receipt-image-dataset-SRD/1003-receipt.jpg", gocv.IMReadColor)
	if img.Empty() {
		panic("cannot read image")
	}
	defer img.Close()

	// Convert to grayscale
	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

	// Gaussian blur to reduce noise
	blurred := gocv.NewMat()
	defer blurred.Close()
	gocv.GaussianBlur(gray, &blurred, image.Pt(5, 5), 0, 0, gocv.BorderDefault)

	// Adaptive thresholding
	thresh := gocv.NewMat()
	defer thresh.Close()
	gocv.AdaptiveThreshold(blurred, &thresh, 255,
		gocv.AdaptiveThresholdMean,
		gocv.ThresholdBinaryInv, 15, 10)

	// Optional: save preprocessed image
	gocv.IMWrite("receipt_preprocessed.png", thresh)

	// Run Tesseract on the preprocessed image
	out, err := exec.Command("tesseract", "receipt_preprocessed.png", "stdout", "--psm", "6").Output()
	if err != nil {
		panic(fmt.Sprintf("tesseract error: %v", err))
	}

	fmt.Println("=== OCR Output ===")
	fmt.Println(string(out))
}
