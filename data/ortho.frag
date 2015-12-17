#version 330 core

uniform sampler2D u_texture;

in vec2 texCoord;

out vec4 fragColor;

void main(void) {
	fragColor = texture2D(u_texture, texCoord);
}