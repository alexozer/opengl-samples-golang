package hierarchy

import (
	"fmt"
	"github.com/go-gl/gl"
	"github.com/ungerik/go3d/vec3"
	math "github.com/barnex/fmath"
)

var (
	posBase = vec3.T{3.0, -5.0, -40.0}
	angBase float32 = -45.0
	posBaseLeft = vec3.T{2.0, 0.0, 0.0}
	posBaseRight = vec3.T{-2.0, 0.0, 0.0}
	scaleBaseZ float32 = 3.0
	angUpperArm float32 = -33.75
	sizeUpperArm float32 = 9.0
	posLowerArm = vec3.T{0.0, 0.0, 8.0}
	angLowerArm float32 = 146.25
	lenLowerArm float32 = 5.0
	widthLowerArm float32 = 1.5
	posWrist = vec3.T{0.0, 0.0, 5.0}
	angWristRoll float32 = 0.0
	angWristPitch float32 = 67.5
	lenWrist float32 = 2.0
	widthWrist float32 = 2.0
	posLeftFinger = vec3.T{1.0, 0.0, 1.0}
	posRightFinger = vec3.T{-1.0, 0.0, 1.0}
	angFingerOpen float32 = 180.0
	lenFinger float32 = 2.0
	widthFinger float32 = 0.5
	angLowerFinger float32 = 45.0
)

func degToRad(degrees float32) (radians float32) {
	return degrees * math.Pi / 180.0
}

func drawNode(matStack *matrixStack) {
	flatMat := flattenMatrix(matStack.mat4())
	modelToCameraMatrixUnif.UniformMatrix4f(false, flatMat)
	gl.DrawElements(gl.TRIANGLES, len(indexData), gl.UNSIGNED_SHORT, uintptr(0))
}

func draw() {
	matStack := NewMatrixStack()

	matStack.Translate(&posBase)
	matStack.RotateY(angBase)

	// Draw left base
	{
		matStack.push()
		matStack.Translate(&posBaseLeft)
		matStack.ScaleVec3(&vec3.T{1.0, 1.0, scaleBaseZ})
		drawNode(matStack)
		matStack.pop()
	}

	// Draw right base
	{
		matStack.push()
		matStack.Translate(&posBaseRight)
		matStack.ScaleVec3(&vec3.T{1.0, 1.0, scaleBaseZ})
		drawNode(matStack)
		matStack.pop()
	}

	drawUpperArm(matStack)
}

func drawUpperArm(matStack *matrixStack) {
	matStack.push()
	matStack.RotateX(angUpperArm)

	{
		matStack.push()
		matStack.Translate(&vec3.T{0, 0, (sizeUpperArm / 2.0) - 1.0})
		matStack.ScaleVec3(&vec3.T{1, 1, sizeUpperArm / 2.0})
		drawNode(matStack)
		matStack.pop()
	}

	drawLowerArm(matStack)

	matStack.pop()
}

func drawLowerArm(matStack *matrixStack) {
	matStack.push()
	matStack.Translate(&posLowerArm)
	matStack.RotateX(degToRad(angLowerArm))

	matStack.push()
	matStack.Translate(&vec3.T{0, 0, lenLowerArm / 2.0})
	matStack.ScaleVec3(&vec3.T{widthLowerArm / 2.0, widthLowerArm / 2.0, lenLowerArm / 2.0})
	drawNode(matStack)
	matStack.pop()

	drawWrist(matStack)

	matStack.pop()
}

func drawWrist(matStack *matrixStack) {
	matStack.push()
	matStack.Translate(&posWrist)
	matStack.RotateZ(degToRad(angWristRoll))
	matStack.RotateX(degToRad(angWristPitch))

	matStack.push()
	matStack.ScaleVec3(&vec3.T{widthWrist / 2.0, widthWrist / 2.0, lenWrist / 2.0})
	drawNode(matStack)
	matStack.pop()

	drawFingers(matStack)

	matStack.pop()
}

func drawFingers(matStack *matrixStack) {
		//Draw left finger
		matStack.push()
		matStack.Translate(&posLeftFinger)
		matStack.RotateY(degToRad(angFingerOpen))

		matStack.push()
		matStack.Translate(&vec3.T{0, 0, lenFinger / 2.0})
		matStack.ScaleVec3(&vec3.T{widthFinger / 2.0, widthFinger/ 2.0, lenFinger / 2.0})
		drawNode(matStack)
		matStack.pop()

		{
			//Draw left lower finger
			matStack.push()
			matStack.Translate(&vec3.T{0, 0, lenFinger})
			matStack.RotateY(degToRad(-angLowerFinger))

			matStack.push()
			matStack.Translate(&vec3.T{0, 0, lenFinger / 2.0})
			matStack.ScaleVec3(&vec3.T{widthFinger / 2.0, widthFinger/ 2.0, lenFinger / 2.0})
			drawNode(matStack)
			matStack.pop()

			matStack.pop()
		}

		matStack.pop()

		//Draw right finger
		matStack.push()
		matStack.Translate(&posRightFinger)
		matStack.RotateY(degToRad(-angFingerOpen))

		matStack.push()
		matStack.Translate(&vec3.T{0, 0, lenFinger / 2.0})
		matStack.ScaleVec3(&vec3.T{widthFinger / 2.0, widthFinger/ 2.0, lenFinger / 2.0})
		drawNode(matStack)
		matStack.pop()

		{
			//Draw right lower finger
			matStack.push()
			matStack.Translate(&vec3.T{0, 0, lenFinger})
			matStack.RotateY(degToRad(angLowerFinger))

			matStack.push()
			matStack.Translate(&vec3.T{0, 0, lenFinger / 2.0})
			matStack.ScaleVec3(&vec3.T{widthFinger / 2.0, widthFinger/ 2.0, lenFinger / 2.0})
			drawNode(matStack)
			matStack.pop()

			matStack.pop()
		}

		matStack.pop()
}

func adjBase(increment bool) {
	if increment {
		angBase += stdAngleInc
	} else {
		angBase -= stdAngleInc
	}

	angBase = math.Mod(degToRad(angBase), degToRad(360.0))
}

func adjUpperArm(increment bool) {
	if increment {
		angUpperArm += stdAngleInc
	} else {
		angUpperArm -= stdAngleInc
	}

	angUpperArm = clamp(degToRad(angUpperArm), degToRad(-90.0), 0.0)
}

func adjLowerArm(increment bool) {
	if increment {
		angLowerArm += stdAngleInc
	} else {
		angLowerArm -= stdAngleInc
	}

	angLowerArm = clamp(degToRad(angLowerArm), 0.0, degToRad(146.25))
}

func adjWristPitch(increment bool) {
	if increment {
		angWristPitch += stdAngleInc
	} else {
		angWristPitch -= stdAngleInc
	}

	angWristPitch = clamp(degToRad(angWristPitch), 0.0, degToRad(90.0))
}

func adjWristRoll(increment bool) {
	if increment {
		angWristRoll += stdAngleInc
	} else {
		angWristRoll -= stdAngleInc
	}

	angWristRoll = math.Mod(degToRad(angWristRoll), degToRad(360.0))
}

func adjFingerOpen(increment bool) {
	if increment {
		angFingerOpen += smallAngleInc
	} else {
		angFingerOpen -= smallAngleInc
	}

	angFingerOpen = clamp(degToRad(angFingerOpen), degToRad(9.0), degToRad(180.0))
}

const (
	stdAngleInc = 11.25
	smallAngleInc = 9.0
)

func WritePose() {
	fmt.Printf("angBase:\t%f\n", angBase)
	fmt.Printf("angUpperArm:\t%f\n", angUpperArm)
	fmt.Printf("angLowerArm:\t%f\n", angLowerArm)
	fmt.Printf("angWristPitch:\t%f\n", angWristPitch)
	fmt.Printf("angWristRoll:\t%f\n", angWristRoll)
	fmt.Printf("angFingerOpen:\t%f\n", angFingerOpen)
	fmt.Printf("\n")
}

func clamp(value, minValue, maxValue float32) float32 {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}

	return value
}
