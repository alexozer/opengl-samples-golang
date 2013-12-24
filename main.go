package main

import (
	"github.com/alexozer/opengl-samples-golang/cpuOffset"
	"github.com/alexozer/opengl-samples-golang/fragChangeColor"
	"github.com/alexozer/opengl-samples-golang/fragPosition"
	"github.com/alexozer/opengl-samples-golang/vertCalcOffset"
	"github.com/alexozer/opengl-samples-golang/vertOffset"
	"github.com/alexozer/opengl-samples-golang/vertexColors"
)

func main() {
	fragPosition.Run()
	vertexColors.Run()
	cpuOffset.Run()
	vertOffset.Run()
	vertCalcOffset.Run()
	fragChangeColor.Run()
}
