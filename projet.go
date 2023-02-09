package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"time"
)

type IMAGE struct {
	x int
	y int
}

type Result struct {
	x     int
	y     int
	color color.RGBA
}

func invertion(img *image.RGBA, x, y int) color.RGBA {
	oldColor := img.At(x, y)
	r, g, b, a := oldColor.RGBA()
	return color.RGBA{
		R: 255 - uint8(r),
		G: 255 - uint8(g),
		B: 255 - uint8(b),
		A: uint8(a),
	}
}

func work(jobChan <-chan IMAGE, resultChan chan<- Result, img *image.RGBA) {
	for image := range jobChan {
		result := invertion(img, image.x, image.y)
		resultChan <- Result{image.x, image.y, result}
		time.Sleep(1)
	}
}

func processImage(img image.Image, jobChan chan<- IMAGE) {
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			jobChan <- IMAGE{x, y}

		}
	}
}

func main() {

	f, err := os.Open("tata.png")
	if err != nil {
		panic(err)
	}
	img, _ := png.Decode(f)
	start := time.Now()
	defer f.Close()
	dst := image.NewRGBA(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	draw.Draw(dst, dst.Bounds(), img, img.Bounds().Min, draw.Src)

	// Etape 5 : utilisation d'un pool de travailleurs
	numWorkers := 10 // nombre de travailleurs

	jobChan := make(chan IMAGE, numWorkers*2)
	resultChan := make(chan Result, numWorkers*2)

	for i := 0; i < numWorkers; i++ {
		go work(jobChan, resultChan, dst)
	}

	go processImage(img, jobChan)

	for compt := 0; compt < img.Bounds().Dx()*img.Bounds().Dy(); compt++ {
		res := <-resultChan
		dst.Set(res.x, res.y, res.color)
	}

	elapsed := time.Since(start)
	f_2, err := os.Create("sorti3.png")
	if err != nil {
		panic(err)
	}
	defer f_2.Close()
	err = png.Encode(f_2, dst)
	if err != nil {
		panic(err)
	}

	fmt.Println("Temps d'exécution:", elapsed)
}
