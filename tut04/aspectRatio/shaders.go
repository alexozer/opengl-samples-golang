package aspectRatio

import (
	"github.com/go-gl/gl"
	"github.com/go-gl/glh"
)

var (
	vertShader = glh.Shader{
		gl.VERTEX_SHADER,
		`#version 100
		precision mediump float;

		attribute vec4 position;
		attribute vec4 color;

		uniform vec2 offset;
		uniform mat4 perspectiveMatrix;

		varying vec4 theColor;

		void main()
		{
			vec4 cameraPos = position + vec4(offset.x, offset.y, 0.0, 0.0);

			gl_Position = perspectiveMatrix * cameraPos;
			theColor = color;
		}
		`,
	}

	fragShader = glh.Shader{
		gl.FRAGMENT_SHADER,
		`#version 100
		precision mediump float;

		varying vec4 theColor;

		void main()
		{
			gl_FragColor = theColor;
		}
		`,
	}
)
