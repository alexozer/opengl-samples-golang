package vertexClipping

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
	perspectiveMatrixUnif.UniformMatrix4f(false, &perspectiveMatrix)
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

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthMask(true)
	gl.DepthFunc(gl.LESS)
	gl.DepthRange(0, 1)

	posBuffer = genVertexBuffer(vertexData)
	indexBuffer = genIndexBuffer(indexData)
	shaderProgram = glh.NewProgram(vertShader, fragShader)
	initAttribs(shaderProgram)
	initVAOs()

	offsetUnif = shaderProgram.GetUniformLocation("offset")
	perspectiveMatrixUnif = shaderProgram.GetUniformLocation("perspectiveMatrix")

	shaderProgram.Use()
	perspectiveMatrixUnif.UniformMatrix4f(false, &perspectiveMatrix)
	gl.ProgramUnuse()

	for !window.ShouldClose() {
		display()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

var (
	posBuffer, indexBuffer gl.Buffer

	vao gl.VertexArray

	shaderProgram               gl.Program
	positionAttrib, colorAttrib gl.AttribLocation
	perspectiveMatrixUnif       gl.UniformLocation
	offsetUnif                  gl.UniformLocation
)

var perspectiveMatrix = [16]float32{
	0:  frustumScale,
	5:  frustumScale,
	10: (zFar + zNear) / (zNear - zFar),
	14: (2 * zFar * zNear) / (zNear - zFar),
	11: -1,
}

// All float32's have a size of 4 bytes
const float32_size = 4

// All uint16's have a size of 2 bytes
const uint16_size = 2
const numVertices = 36

const (
	zNear, zFar  float32 = 1, 3
	frustumScale float32 = 1
)

func genVertexBuffer(verts []float32) gl.Buffer {
	buffer := gl.GenBuffer()
	buffer.Bind(gl.ARRAY_BUFFER)
	defer buffer.Unbind(gl.ARRAY_BUFFER)

	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*float32_size, verts, gl.STATIC_DRAW)

	return buffer
}

func genIndexBuffer(indices []uint16) gl.Buffer {
	buffer := gl.GenBuffer()
	buffer.Bind(gl.ELEMENT_ARRAY_BUFFER)
	defer buffer.Unbind(gl.ELEMENT_ARRAY_BUFFER)

	indexBufSize := len(indexData) * uint16_size
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indexBufSize, indexData, gl.STATIC_DRAW)

	return buffer
}

func initAttribs(program gl.Program) {
	positionAttrib = program.GetAttribLocation("position")
	colorAttrib = program.GetAttribLocation("color")
}

func initVAOs() {
	vao = gl.GenVertexArray()
	vao.Bind()

	var colorDataOffset uintptr = float32_size * 3 * numVertices

	posBuffer.Bind(gl.ARRAY_BUFFER)
	positionAttrib.EnableArray()
	colorAttrib.EnableArray()
	positionAttrib.AttribPointer(3, gl.FLOAT, false, 0, uintptr(0))
	colorAttrib.AttribPointer(4, gl.FLOAT, false, 0, colorDataOffset)
	indexBuffer.Bind(gl.ELEMENT_ARRAY_BUFFER)

	// Unbind vao
	gl.VertexArray(0).Bind()
}

func display() {
	gl.ClearColor(0, 0, 0, 0)
	gl.ClearDepth(1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	shaderProgram.Use()
	defer gl.ProgramUnuse()

	vao.Bind()
	defer gl.VertexArray(0).Bind()
	offsetUnif.Uniform3f(0, 0, 0.5)
	gl.DrawElements(gl.TRIANGLES, len(indexData), gl.UNSIGNED_SHORT, uintptr(0))

	offsetUnif.Uniform3f(0, 0, -1)
	gl.DrawElementsBaseVertex(gl.TRIANGLES, len(indexData), gl.UNSIGNED_SHORT, uintptr(0), numVertices/2)
}
