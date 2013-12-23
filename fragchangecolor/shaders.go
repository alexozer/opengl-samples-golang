package fragchangecolor

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

		uniform float loopDuration;
		uniform float time;

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
		}
		`,
	}

	fragShader = glh.Shader{
		gl.FRAGMENT_SHADER,
		`#version 100

		precision mediump float;

		uniform float loopDuration;
		uniform float time;

		const vec4 firstColor = vec4(1.0f, 1.0f, 1.0f, 1.0f); 
		const vec4 secondColor = vec4(0.0f, 0.0f, 1.0f, 1.0f); 

		void main()
		{
			float currLerp = mod(time, loopDuration) / loopDuration;
			gl_FragColor = mix(firstColor, secondColor, currLerp);
		}
		`,
	}
)
