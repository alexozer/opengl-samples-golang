package overlapNoDepth

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

func onResize(w *glfw.Window, width, height int) {
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
	window.SetSizeCallback(onResize)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	gl.Init()

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CW)

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

	vao1, vao2 gl.VertexArray

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
	vao1 = gl.GenVertexArray()
	vao1.Bind()

	var colorDataOffset uintptr = float32_size * 3 * numVertices

	posBuffer.Bind(gl.ARRAY_BUFFER)
	positionAttrib.EnableArray()
	colorAttrib.EnableArray()
	positionAttrib.AttribPointer(3, gl.FLOAT, false, 0, uintptr(0))
	colorAttrib.AttribPointer(4, gl.FLOAT, false, 0, colorDataOffset)
	indexBuffer.Bind(gl.ELEMENT_ARRAY_BUFFER)

	// Unbind vao1
	gl.VertexArray(0).Bind()

	vao2 = gl.GenVertexArray()
	vao2.Bind()

	var posDataOffset uintptr = float32_size * 3 * (numVertices / 2)
	colorDataOffset += float32_size * 4 * (numVertices / 2)

	//Use the same buffer object previously bound to GL_ARRAY_BUFFER.
	positionAttrib.EnableArray()
	colorAttrib.EnableArray()
	positionAttrib.AttribPointer(3, gl.FLOAT, false, 0, posDataOffset)
	colorAttrib.AttribPointer(4, gl.FLOAT, false, 0, colorDataOffset)
	indexBuffer.Bind(gl.ELEMENT_ARRAY_BUFFER)

	// Unbind vao2
	gl.VertexArray(0).Bind()
}

func display() {
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	shaderProgram.Use()
	defer gl.ProgramUnuse()

	vao1.Bind()
	offsetUnif.Uniform3f(0, 0, 0)
	gl.DrawElements(gl.TRIANGLES, len(indexData), gl.UNSIGNED_SHORT, uintptr(0))

	vao2.Bind()
	defer gl.VertexArray(0).Bind()
	offsetUnif.Uniform3f(0, 0, -1)
	gl.DrawElements(gl.TRIANGLES, len(indexData), gl.UNSIGNED_SHORT, uintptr(0))
}
