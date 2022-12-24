package reference

import (
	"ever-parse/internal/util"
	"fmt"
	"image"
	"image/png"
	"os"
)

func imgStuff(path string) (img image.Image, didCrop bool) {
	imgFile, err := os.OpenFile(path, os.O_RDONLY, 0644)
	util.CheckWithoutPanic(err, "imagefile", path)

	srcImg, err := png.Decode(imgFile)
	util.CheckWithoutPanic(err, "srcImg image", path)

	type subber interface {
		SubImage(r image.Rectangle) image.Image
	}

	subs, ok := srcImg.(subber)

	biggest := util.Ternary(srcImg.Bounds().Dx() > srcImg.Bounds().Dy(), srcImg.Bounds().Dx(), srcImg.Bounds().Dy())

	if ok && biggest > 4000 {
		firstY := traverseFirst(srcImg, xA)
		firstX := traverseFirst(srcImg, yA)
		lastY := traverseLast(srcImg, xA)
		lastX := traverseLast(srcImg, yA)

		rect := image.Rect(firstX.X, firstY.Y, lastX.X, lastY.Y)

		fmt.Printf("Cropping [%s] down to [%+v]\n", path, rect)
		return subs.SubImage(rect), true
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

	dy = util.Ternary(dy > dx, dx, dy)

	size := dy * dy
	counter := 0
	ax := bool(axis)

	for counter < size {
		x := util.Ternary(ax, counter%dy, counter/dy)
		y := util.Ternary(ax, counter/dy, counter%dy)

		alpha := GetAlphaAt(img, x, y)
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

		alpha := GetAlphaAt(img, x, y)
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

func GetAlphaAt(img image.Image, x int, y int) uint32 {
	_, _, _, a := img.At(x, y).RGBA()
	return a
}