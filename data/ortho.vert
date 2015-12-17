#version 330 core

uniform mat4 u_projectionMatrix;

layout (location = 0) in vec4 i_position;
layout (location = 2) in vec2 i_texCoord;

out vec2 texCoord;

void main(void) {
	gl_Position = u_projectionMatrix * vec4(i_position.xy, 0.0, 1.0);
	texCoord = i_texCoord;
}