package fragposition

import (
	"github.com/go-gl/glh"
	"github.com/go-gl/gl"
)

var (
	vertShader = glh.Shader{
		gl.VERTEX_SHADER,
		`#version 100
		precision mediump float;
		
		attribute vec4 position;
		void main()
		{
			gl_Position = position;
		}
		`,
	}

	fragShader = glh.Shader{
		gl.FRAGMENT_SHADER,
		`#version 100
		precision mediump float;

		uniform int height;

		void main()
		{
			float lerpValue = gl_FragCoord.y / float(height);
			gl_FragColor = mix(vec4(1.0f, 1.0f, 0.64f, 1.0f), 
					vec4(0.25f, 1.0f-lerpValue*0.5, lerpValue, 0.0f), 
					lerpValue);
		}
		`,
	}
)

