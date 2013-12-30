package hierarchy

import (
	"github.com/ungerik/go3d/mat4"
	"github.com/ungerik/go3d/vec3"
	"github.com/ungerik/go3d/vec4"
)

type matrixStack struct {
	mat4.T
	stack []mat4.T
}

func NewMatrixStack() *matrixStack {
	return &matrixStack{
		mat4.Ident,
		[]mat4.T{}, // Start with len() == 0, let append() extend
	}
}

// Override mat4.T.Translate()
func (ms *matrixStack) Translate(v *vec3.T) {
	transMat := mat4.Ident
	transMat.SetTranslation(v)
	ms.AssignMul(&transMat)
}

func (ms *matrixStack) RotateX(angle float32) {
	rotationMat := mat4.Ident
	rotationMat.AssignXRotation(angle)
	ms.AssignMul(&rotationMat)
}

func (ms *matrixStack) RotateY(angle float32) {
	rotationMat := mat4.Ident
	rotationMat.AssignYRotation(angle)
	ms.AssignMul(&rotationMat)
}

func (ms *matrixStack) RotateZ(angle float32) {
	rotationMat := mat4.Ident
	rotationMat.AssignZRotation(angle)
	ms.AssignMul(&rotationMat)
}

func (ms *matrixStack) ScaleVec3(s *vec3.T) {
	scaleVec := vec4.T{s[0], s[1], s[2], 1}
	scaleMat := mat4.Ident

	scaleMat.SetScaling(&scaleVec)

	ms.AssignMul(&scaleMat)
}

// Apparently, go3d.AssignMul() doesn't multiply matrices
// in the usual fashion. This function will do for now.
func (ms *matrixStack) AssignMul(a *mat4.T) {
	product := mat4.Zero
	for col := 0; col < a.Cols(); col++ {
		for row := 0; row < a.Rows(); row++ {
			for vecPos := 0; vecPos < a.Rows(); vecPos++ {
				product[col][row] += ms.T[vecPos][row] * a[col][vecPos]
			}
		}
	}
	ms.T = product
}

func (ms *matrixStack) push() {
	ms.stack = append(ms.stack, ms.T)
}

func (ms *matrixStack) pop() {
	currLength := len(ms.stack)

	if currLength > 0 {
		ms.T = ms.stack[currLength-1]
		ms.stack = ms.stack[:currLength-1]
	}
}

func (ms *matrixStack) mat4() *mat4.T {
	return &ms.T
}
