package vertOffset

import (
	"fmt"
	"math"

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

func onResize(w *glfw.Window, width, height int) {
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
	window.SetSizeCallback(onResize)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	gl.Init()

	posBuffer = genVertexBuffer(vertices)
	shaderProgram = glh.NewProgram(vertShader, fragShader)
	offsetLocation = shaderProgram.GetUniformLocation("offset")

	for !window.ShouldClose() {
		display()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

var vertices = []float32{
	0, 0.5, 0, 1,
	0.5, -0.366, 0, 1,
	-0.5, -0.366, 0, 1,
	1, 0, 0, 1,
	0, 1, 0, 1,
	0, 0, 1, 1,
}

var (
	posBuffer      gl.Buffer
	shaderProgram  gl.Program
	offsetLocation gl.UniformLocation
)

// All float32's have a size of 4 bytes
const float32_size = 4

func genVertexBuffer(verts []float32) gl.Buffer {
	buffer := gl.GenBuffer()
	buffer.Bind(gl.ARRAY_BUFFER)
	defer buffer.Unbind(gl.ARRAY_BUFFER)

	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*float32_size, verts, gl.STATIC_DRAW)

	return buffer
}

func calcOffsets() (xOffset, yOffset float32) {
	const speed = 3.5
	const size = 0.5
	angle := math.Mod(glfw.GetTime(), speed) / speed * 2.0 * math.Pi

	return size * float32(math.Cos(angle)), size * float32(math.Sin(angle))
}

func display() {
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	shaderProgram.Use()
	defer gl.ProgramUnuse()

	// The pleasure of passing a multivalue return as parameters...
	offsetLocation.Uniform2f(calcOffsets())

	posBuffer.Bind(gl.ARRAY_BUFFER)

	positionAttrib := gl.AttribLocation(shaderProgram.GetAttribLocation("position"))
	positionAttrib.AttribPointer(4, gl.FLOAT, false, 0, uintptr(0))
	positionAttrib.EnableArray()
	defer positionAttrib.DisableArray()

	colorAttrib := gl.AttribLocation(shaderProgram.GetAttribLocation("color"))
	colorAttrib.AttribPointer(4, gl.FLOAT, false, 0, uintptr((len(vertices)*float32_size)/2))
	colorAttrib.EnableArray()
	defer colorAttrib.DisableArray()

	gl.DrawArrays(gl.TRIANGLES, 0, len(vertices)/2/float32_size)
}
