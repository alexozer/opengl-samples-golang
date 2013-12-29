package scale

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/go-gl/glh"
	"github.com/ungerik/go3d/mat4"
	"github.com/ungerik/go3d/vec3"
	"github.com/ungerik/go3d/vec4"
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
	cameraToClipMatrix[0] = frustumScale / (float32(width) / float32(height))
	cameraToClipMatrix[5] = frustumScale

	shaderProgram.Use()
	cameraToClipMatrixUnif.UniformMatrix4f(false, &cameraToClipMatrix)
	gl.ProgramUnuse()

	gl.Viewport(0, 0, width, height)
}

func Run() {
	glfw.SetErrorCallback(onError)

	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(640, 640, "Graphics", nil, nil)
	if err != nil {
		panic(err)
	}
	window.SetKeyCallback(onKey)
	window.SetSizeCallback(reshape)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	gl.Init()

	posBuffer = genVertexBuffer(vertexData)
	indexBuffer = genIndexBuffer(indexData)
	shaderProgram = glh.NewProgram(vertShader, fragShader)
	initAttribs(shaderProgram)
	vao = initVAO()

	modelToCameraMatrixUnif = shaderProgram.GetUniformLocation("modelToCameraMatrix")
	cameraToClipMatrixUnif = shaderProgram.GetUniformLocation("cameraToClipMatrix")

	shaderProgram.Use()
	cameraToClipMatrixUnif.UniformMatrix4f(false, &cameraToClipMatrix)
	gl.ProgramUnuse()

	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CW)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthMask(true)
	gl.DepthFunc(gl.LEQUAL)
	gl.DepthRange(0, 1)

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

const zNear, zFar = 1.0, 61.0

var frustumScale float32 = calcFrustumScale(45)

var spikeys = []instance{
	{NullScale, vec3.T{0.0, 0.0, -45.0}},
	{StaticUniformScale, vec3.T{-10.0, -10.0, -45.0}},
	{StaticNonUniformScale, vec3.T{-10.0, 10.0, -45.0}},
	{DynamicUniformScale, vec3.T{10.0, 10.0, -45.0}},
	{DynamicNonUniformScale, vec3.T{10.0, -10.0, -45.0}},
}

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

// instance contains a position and a function that generates a
// scale based on the current time
type instance struct {
	scaleFunc func(elapsedTime float64) (scale *vec4.T)
	position  vec3.T
}

func (obj instance) constructMatrix(elapsedTime float64) *mat4.T {
	scale := obj.scaleFunc(elapsedTime)
	mat := mat4.Ident
	mat.SetScaling(scale).SetTranslation(&obj.position)
	return &mat
}

func calcLerpFactor(elapsedTime, loopDuration float64) float32 {
	value := math.Mod(elapsedTime, loopDuration) / loopDuration
	if value > 0.5 {
		value = 1.0 - value
	}

	return float32(value * 2.0)
}

func mix(val1, val2, lerp float32) float32 {
	antiLerp := 1.0 - lerp

	return val1*antiLerp + val2*lerp
}

func NullScale(elapsedTime float64) (scale *vec4.T) {
	return &vec4.T{1, 1, 1, 1}
}

func StaticUniformScale(elapsedTime float64) (scale *vec4.T) {
	return &vec4.T{4, 4, 4, 1}
}

func StaticNonUniformScale(elapsedTime float64) (scale *vec4.T) {
	return &vec4.T{0.5, 1.0, 10.0, 1}
}

func DynamicUniformScale(elapsedTime float64) (scale *vec4.T) {
	const loopDuration = 3.0

	var initScale, finalScale = vec4.T{1, 1, 1, 1}, vec4.T{4, 4, 4, 1}
	lerpVal := calcLerpFactor(elapsedTime, loopDuration)

	scaleVec := vec4.Interpolate(&initScale, &finalScale, lerpVal)
	return &scaleVec
}

func DynamicNonUniformScale(elapsedTime float64) (scale *vec4.T) {
	const xLoopDuration, zLoopDuration = 3.0, 5.0

	return &vec4.T{
		mix(1.0, 0.5, calcLerpFactor(elapsedTime, xLoopDuration)),
		1.0,
		mix(1.0, 10.0, calcLerpFactor(elapsedTime, zLoopDuration)),
		1,
	}
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

	currTime := glfw.GetTime()
	for _, spikey := range spikeys {
		transformMatrix := spikey.constructMatrix(currTime)
		modelToCameraMatrixUnif.UniformMatrix4f(false, flattenMatrix(transformMatrix))

		gl.DrawElements(gl.TRIANGLES, len(indexData), gl.UNSIGNED_SHORT, uintptr(0))
	}
}
