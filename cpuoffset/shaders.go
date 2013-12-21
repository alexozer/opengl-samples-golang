package cpuoffset

import (
	"github.com/go-gl/gl"
	"github.com/go-gl/glh"
)

var (
	vertShader = glh.Shader{
		gl.VERTEX_SHADER,
		`
		#version 100

		precision mediump float;

		attribute vec4 position;
		attribute vec4 color;

		varying vec4 theColor;

		void main()
		{
			gl_Position = position;
			theColor = color;
		}
		`,
	}

	fragShader = glh.Shader{
		gl.FRAGMENT_SHADER,
		`
		#version 100

		precision mediump float;

		varying vec4 theColor;

		void main()
		{
			gl_FragColor = theColor;
		}
		`,
	}	
)
