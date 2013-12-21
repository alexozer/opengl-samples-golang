package cpuoffset

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/go-gl/glh"
	"math"
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
	window.SetFramebufferSizeCallback(reshape)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	gl.Init()

	posBuffer = genVertexBuffer(vertices)
	shaderProgram = glh.NewProgram(vertShader, fragShader)

	for !window.ShouldClose() {
		updateVertexBuffer(posBuffer, vertices)
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
	posBuffer     gl.Buffer
	shaderProgram gl.Program
)

// All float32's have a size of 4 bytes
const float32_size = 4

func genVertexBuffer(verts []float32) gl.Buffer {
	buffer := gl.GenBuffer()
	buffer.Bind(gl.ARRAY_BUFFER)
	defer buffer.Unbind(gl.ARRAY_BUFFER)

	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*float32_size, verts, gl.STREAM_DRAW)

	return buffer
}

func updateVertexBuffer(buf gl.Buffer, verts []float32) {
	newverts := make([]float32, len(verts))
	copy(newverts, verts)

	const speed = 3.5
	const size = 0.5
	angle := math.Mod(glfw.GetTime(), speed) / speed * 2.0 * math.Pi
	xOffset, yOffset := size*float32(math.Cos(angle)), size*float32(math.Sin(angle))


	// Intentionally modify the color values as well, because it looks awesome.

	//for i := 0; i < len(newverts)/2; i += 4 {
	for i := 0; i < len(newverts); i += 4 {
		newverts[i] += xOffset
		newverts[i+1] += yOffset
	}

	buf.Bind(gl.ARRAY_BUFFER)
	defer buf.Unbind(gl.ARRAY_BUFFER)

	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(newverts)*float32_size, newverts)
}

func display() {
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	shaderProgram.Use()
	defer gl.ProgramUnuse()

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
