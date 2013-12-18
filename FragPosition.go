package main

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

func main() {
	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	defer glfw.Terminate()

	glfw.SwapInterval(1)

	window, err := glfw.CreateWindow(640, 480, "Graphics", nil, nil)
	if err != nil {
		panic(err)
	}
	window.SetKeyCallback(keyCallback)
	window.SetFramebufferSizeCallback(reshape)
	window.MakeContextCurrent()

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
	posBuffer gl.Buffer
	shaderProgram gl.Program
)

var (
	vertShader = glh.Shader{
		gl.VERTEX_SHADER,
		`
		#version 100
		attribute vec4 position;
		void main()
		{
			gl_Position = position;
		}`,
	}

	fragShader = glh.Shader{
		gl.FRAGMENT_SHADER,
		`
		#version 100

		precision mediump float;
		void main()
		{
			float lerpValue = (gl_FragCoord.x + gl_FragCoord.y) / 1500.0f;
			gl_FragColor = mix(vec4(1.0f, 1.0f, 0.64f, 1.0f), 
					vec4(0.25f, 1.0f-lerpValue*0.5, lerpValue, 0.0f), 
					lerpValue);
		}

		`,
	}
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
