package shaderPerspective

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
		uniform float zNear;
		uniform float zFar;
		uniform float frustumScale;

		varying vec4 theColor;

		void main()
		{
			vec4 cameraPos = position + vec4(offset.x, offset.y, 0.0, 0.0);
			vec4 clipPos;

			clipPos.xy = cameraPos.xy * frustumScale;

			clipPos.z = cameraPos.z * (zNear + zFar) / (zNear - zFar);
			clipPos.z += 2.0 * zNear * zFar / (zNear - zFar);

			clipPos.w = -cameraPos.z;

			gl_Position = clipPos;
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
