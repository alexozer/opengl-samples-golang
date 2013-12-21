package fragposition

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/go-gl/glh"
)

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v.n", err, desc)
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	}
}

func reshape(w *glfw.Window, width, height int) {
	gl.Viewport(0, 0, width, height)
}

func Run() {
	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(640, 480, "Graphics", nil, nil)
	if err != nil {
		panic(err)
	}
	window.SetKeyCallback(keyCallback)
	window.SetFramebufferSizeCallback(reshape)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	gl.Init()

	posBuffer = genVertexBuffer(vertices)
	shaderProgram = glh.NewProgram(vertShader, fragShader)

	for !window.ShouldClose() {
		display()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

var vertices = []float32{
	0.75, 0.75, 0, 1,
	0.75, -0.75, 0, 1,
	-0.75, -0.75, 0, 1,
}

var (
	posBuffer     gl.Buffer
	shaderProgram gl.Program
)

// All float32's have a size of 4 bytes
const float32_size = 4

func genVertexBuffer(verts []float32) gl.Buffer {
	buffer := gl.GenBuffer()
	buffer.Bind(gl.ARRAY_BUFFER)

	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*float32_size, verts, gl.STATIC_DRAW)

	buffer.Unbind(gl.ARRAY_BUFFER)
	return buffer
}

func display() {
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	shaderProgram.Use()

	posBuffer.Bind(gl.ARRAY_BUFFER)
	attribLoc := gl.AttribLocation(shaderProgram.GetAttribLocation("position"))
	attribLoc.EnableArray()

	attribLoc.AttribPointer(4, gl.FLOAT, false, 0, uintptr(0))
	gl.DrawArrays(gl.TRIANGLES, 0, len(vertices)/4)

	attribLoc.DisableArray()
	gl.ProgramUnuse()
}