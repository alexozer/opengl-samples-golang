package hierarchy

import (
	"fmt"
	"math"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/go-gl/glh"

	"github.com/ungerik/go3d/mat4"
)

func onError(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Press {
		return
	}

	switch key {
	case glfw.KeyEscape:
		w.SetShouldClose(true)
	case glfw.KeyA:
		adjBase(true)
	case glfw.KeyD:
		adjBase(false)
	case glfw.KeyW:
		adjUpperArm(false)
	case glfw.KeyS:
		adjUpperArm(true)
	case glfw.KeyR:
		adjLowerArm(false)
	case glfw.KeyF:
		adjLowerArm(true)
	case glfw.KeyT:
		adjWristPitch(false)
	case glfw.KeyG:
		adjWristPitch(true)
	case glfw.KeyZ:
		adjWristRoll(true)
	case glfw.KeyC:
		adjWristRoll(false)
	case glfw.KeyQ:
		adjFingerOpen(true)
	case glfw.KeyE:
		adjFingerOpen(false)
	}
}

func onResize(w *glfw.Window, width, height int) {
	cameraToClipMatrix[0] = frustumScale / (float32(width) / float32(height))
	cameraToClipMatrix[5] = frustumScale

	shaderProgram.Use()
	cameraToClipMatrixUnif.UniformMatrix4f(false, &cameraToClipMatrix)
	gl.ProgramUnuse()

	gl.Viewport(0, 0, width, height)
}

func initGlfw() {
	glfw.SetErrorCallback(onError)

	if !glfw.Init() {
		panic("Can't init glfw!")
	}

	window, err := glfw.CreateWindow(640, 640, "Graphics", nil, nil)
	if err != nil {
		panic(err)
	}
	window.SetKeyCallback(onKey)
	window.SetSizeCallback(onResize)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	glWindow = window
}

func initGl() {
	gl.Init()

	posBuffer = genVertexBuffer(vertexData)
	indexBuffer = genIndexBuffer(indexData)
	shaderProgram = glh.NewProgram(vertShader, fragShader)
	initAttribs(shaderProgram)
	vao = initVAO()

	modelToCameraMatrixUnif = shaderProgram.GetUniformLocation("modelToCameraMatrix")
	cameraToClipMatrixUnif = shaderProgram.GetUniformLocation("cameraToClipMatrix")

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CW)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthMask(true)
	gl.DepthFunc(gl.LEQUAL)
	gl.DepthRange(0, 1)
}

func Run() {
	initGlfw()
	defer glfw.Terminate()

	initGl()

	shaderProgram.Use()
	cameraToClipMatrixUnif.UniformMatrix4f(false, &cameraToClipMatrix)
	gl.ProgramUnuse()

	for !glWindow.ShouldClose() {
		display()
		glWindow.SwapBuffers()
		glfw.PollEvents()
	}
}

var (
	glWindow *glfw.Window

	posBuffer, indexBuffer gl.Buffer

	vao gl.VertexArray

	shaderProgram               gl.Program
	positionAttrib, colorAttrib gl.AttribLocation
	modelToCameraMatrixUnif     gl.UniformLocation
	cameraToClipMatrixUnif      gl.UniformLocation
)

var cameraToClipMatrix = [16]float32{
	0:  frustumScale,
	5:  frustumScale,
	10: (zFar + zNear) / (zNear - zFar),
	14: (2 * zFar * zNear) / (zNear - zFar),
	11: -1,
}

const float32_size = 4 // All float32's have a size of 4 bytes
const uint16_size = 2  // All uint16's have a size of 2 bytes

const zNear, zFar = 1.0, 100.0

var frustumScale float32 = calcFrustumScale(45)

/*
 * The greater the angle, the greater the field of view captured,
 * which makes the objects seem smaller. Therefore, the larger the
 * angle, the smaller the scale.
 */
func calcFrustumScale(angleDeg float64) float32 {
	angleRad := angleDeg * math.Pi / 180.0
	return float32(1.0 / math.Tan(angleRad/2.0))
}

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

func initVAO() (vao gl.VertexArray) {
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

	return
}

func flattenMatrix(mat *mat4.T) *[16]float32 {
	// This is more efficient than a loop... right?
	return &[16]float32{
		mat[0][0], mat[0][1], mat[0][2], mat[0][3],
		mat[1][0], mat[1][1], mat[1][2], mat[1][3],
		mat[2][0], mat[2][1], mat[2][2], mat[2][3],
		mat[3][0], mat[3][1], mat[3][2], mat[3][3],
	}
}

func display() {
	gl.ClearColor(0, 0, 0, 0)
	gl.ClearDepth(1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	shaderProgram.Use()
	defer gl.ProgramUnuse()

	vao.Bind()
	defer gl.VertexArray(0).Bind()

	draw()
}
