package vertcalcoffset

import (
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

func displayDual() {
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	shaderProgram.Use()
	defer gl.ProgramUnuse()

	posBuffer.Bind(gl.ARRAY_BUFFER)
	defer posBuffer.Unbind(gl.ARRAY_BUFFER)

	timeUniform.Uniform1f(float32(glfw.GetTime()))
	loopDurUniform.Uniform1f(loopDuration)

	positionAttrib := gl.AttribLocation(shaderProgram.GetAttribLocation("position"))
	positionAttrib.AttribPointer(4, gl.FLOAT, false, 0, uintptr(0))
	positionAttrib.EnableArray()
	defer positionAttrib.DisableArray()

	colorAttrib := gl.AttribLocation(shaderProgram.GetAttribLocation("color"))
	colorAttrib.AttribPointer(4, gl.FLOAT, false, 0, uintptr((len(vertices)*float32_size)/2))
	colorAttrib.EnableArray()
	defer colorAttrib.DisableArray()

	// Draw the first triangle
	gl.DrawArrays(gl.TRIANGLES, 0, len(vertices)/2/float32_size)

	// Draw the second triangle
	timeUniform.Uniform1f(float32(glfw.GetTime() + (0.5 * loopDuration)))
	gl.DrawArrays(gl.TRIANGLES, 0, len(vertices)/2/float32_size)
}
