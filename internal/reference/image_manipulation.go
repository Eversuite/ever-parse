package reference

import (
	"ever-parse/internal/util"
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
)

func cropImage(path string) (img image.Image, didCrop bool) {
	imgFile, err := os.OpenFile(path, os.O_RDONLY, 0644)
	util.CheckWithoutPanic(err, "imagefile", path)

	srcImg, err := png.Decode(imgFile)
	util.CheckWithoutPanic(err, "srcImg image", path)

	type cropableImage interface {
		SubImage(r image.Rectangle) image.Image
	}

	cropable, ok := srcImg.(cropableImage)
	longestBorder := math.Max(float64(srcImg.Bounds().Dx()), float64(srcImg.Bounds().Dy()))

	if ok && longestBorder > 4000 {
		// I'm pretty sure traverseFirst and traverseSecond can somehow be combined but at this point my brain is melted...
		firstY := traverseFirst(srcImg, xA)
		firstX := traverseFirst(srcImg, yA)
		lastY := traverseLast(srcImg, xA)
		lastX := traverseLast(srcImg, yA)
		rect := image.Rect(firstX.X, firstY.Y, lastX.X, lastY.Y)

		fmt.Printf("Cropped to [%+v] from [%+v] for [%s]\n", rect, srcImg.Bounds(), path)

		return cropable.SubImage(rect), true
	}

	return srcImg, false
}

type Axis bool

const (
	xA Axis = true
	yA Axis = false
)

func traverseFirst(img image.Image, axis Axis) image.Point {

	dy := img.Bounds().Dy()
	dx := img.Bounds().Dx()

	size := dx * dy
	counter := 0
	ax := bool(axis)

	for counter < size {
		x := util.Ternary(ax, counter%dy, counter/dx)
		y := util.Ternary(ax, counter/dy, counter%dx)

		_, _, _, alpha := img.At(x, y).RGBA()
		if alpha > 0 {
			return image.Point{
				X: x,
				Y: y,
			}
		}
		counter += 5
	}

	// Could also return an error but why should I
	return image.Point{
		X: -1,
		Y: -1,
	}
}

func traverseLast(img image.Image, axis Axis) image.Point {

	dy := img.Bounds().Dy()
	size := dy * dy
	counter := size
	ax := bool(axis)

	for counter > 0 {
		x := util.Ternary(ax, counter%dy, counter/dy)
		y := util.Ternary(ax, counter/dy, counter%dy)

		_, _, _, alpha := img.At(x, y).RGBA()
		if alpha > 0 {
			return image.Point{
				X: x,
				Y: y,
			}
		}
		counter -= 5
	}

	// Could also return an error but why should I
	return image.Point{
		X: -1,
		Y: -1,
	}
}