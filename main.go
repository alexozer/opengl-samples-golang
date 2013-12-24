package main

import (
	"github.com/alexozer/opengl-samples-golang/tut02/fragPosition"
	"github.com/alexozer/opengl-samples-golang/tut02/vertexColors"
	"github.com/alexozer/opengl-samples-golang/tut03/cpuOffset"
	"github.com/alexozer/opengl-samples-golang/tut03/fragChangeColor"
	"github.com/alexozer/opengl-samples-golang/tut03/vertCalcOffset"
	"github.com/alexozer/opengl-samples-golang/tut03/vertOffset"
)

func main() {
	fragPosition.Run()
	vertexColors.Run()
	cpuOffset.Run()
	vertOffset.Run()
	vertCalcOffset.Run()
	fragChangeColor.Run()
}
