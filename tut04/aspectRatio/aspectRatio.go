package aspectRatio

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/go-gl/glh"
)

func onError(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	}
}

func reshape(w *glfw.Window, width, height int) {
	perspectiveMatrix[0] = frustumScale / (float32(width) / float32(height))
	perspectiveMatrix[5] = frustumScale

	shaderProgram.Use()
	perspectiveMatrixLoc.UniformMatrix4f(false, &perspectiveMatrix)
	gl.ProgramUnuse()

	gl.Viewport(0, 0, width, height)
}

func Run() {
	glfw.SetErrorCallback(onError)

	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(640, 480, "Graphics", nil, nil)
	if err != nil {
		panic(err)
	}
	window.SetKeyCallback(onKey)
	window.SetSizeCallback(reshape)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	gl.Init()

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CW)

	posBuffer = genVertexBuffer(vertexData)
	shaderProgram = glh.NewProgram(vertShader, fragShader)

	offsetLocation := shaderProgram.GetUniformLocation("offset")
	perspectiveMatrixLoc = shaderProgram.GetUniformLocation("perspectiveMatrix")

	shaderProgram.Use()
	offsetLocation.Uniform2f(xOffset, yOffset)
	perspectiveMatrixLoc.UniformMatrix4f(false, &perspectiveMatrix)
	gl.ProgramUnuse()

	for !window.ShouldClose() {
		display()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

var (
	posBuffer            gl.Buffer
	shaderProgram        gl.Program
	perspectiveMatrixLoc gl.UniformLocation

	perspectiveMatrix = [16]float32{
		0:  frustumScale,
		5:  frustumScale,
		10: (zFar + zNear) / (zNear - zFar),
		14: (2 * zFar * zNear) / (zNear - zFar),
		11: -1,
	}
)

// All float32's have a size of 4 bytes
const float32_size = 4
const xOffset, yOffset float32 = 1.5, 0.5
const zNear, zFar float32 = 1, 3
const frustumScale float32 = 1

func genVertexBuffer(verts []float32) gl.Buffer {
	buffer := gl.GenBuffer()
	buffer.Bind(gl.ARRAY_BUFFER)
	defer buffer.Unbind(gl.ARRAY_BUFFER)

	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*float32_size, verts, gl.STATIC_DRAW)

	return buffer
}

func display() {
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	posBuffer.Bind(gl.ARRAY_BUFFER)
	defer posBuffer.Unbind(gl.ARRAY_BUFFER)

	shaderProgram.Use()
	defer gl.ProgramUnuse()

	positionAttrib := gl.AttribLocation(shaderProgram.GetAttribLocation("position"))
	positionAttrib.AttribPointer(4, gl.FLOAT, false, 0, uintptr(0))
	positionAttrib.EnableArray()
	defer positionAttrib.DisableArray()

	colorAttrib := gl.AttribLocation(shaderProgram.GetAttribLocation("color"))
	colorAttrib.AttribPointer(4, gl.FLOAT, false, 0, uintptr((len(vertexData)*float32_size)/2))
	colorAttrib.EnableArray()
	defer colorAttrib.DisableArray()

	gl.DrawArrays(gl.TRIANGLES, 0, len(vertexData)/2/float32_size)
}
