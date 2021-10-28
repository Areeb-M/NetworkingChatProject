package main

import (
	"image"
	"math/rand"
)

func funtionA() int {
	return 10000 + rand.Intn(99999-10000)
}

func Engine(parameters RasterParameters) image.RGBA {
	raster := image.NewRGBA(image.Rect(0, 0, parameters.rasterWidth, parameters.rasterHeight))

	var dx float64 = (parameters.xMax - parameters.xMin) / float64(parameters.rasterWidth)
	var di float64 = (parameters.iMax - parameters.iMin) / float64(parameters.rasterHeight)
	var c complex128

	colorBlend := ColorBlend{rainbowPalette2[:], parameters.maxIterations}

	for rx := 0; rx < parameters.rasterWidth; rx++ {
		for ry := 0; ry < parameters.rasterHeight; ry++ {
			c = complex(dx*float64(rx)+parameters.xMin, di*float64(ry)+parameters.iMin)

			depth := parameters.fractalFunction(c, parameters.divergenceThreshold, parameters.maxIterations)
			raster.Set(rx, ry, colorBlend.Blend(depth))
		}
	}

	return *raster
}
