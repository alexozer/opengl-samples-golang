package vertcalcoffset

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

		uniform float loopDuration;
		uniform float time;

		varying vec4 theColor;

		void main()
		{
			float loopPos = mod(time, loopDuration) / loopDuration;
			float fullCircle = 2.0f * 3.14159f;

			vec4 totalOffset = vec4(
				0.5f * cos(loopPos * fullCircle),
				0.5f * sin(loopPos * fullCircle),
				0.0f,
				0.0f);

			gl_Position = position + totalOffset;
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
