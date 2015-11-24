#version 330 core

uniform sampler2D u_texture;
uniform float u_flash;
uniform float u_alpha;

in vec2 texCoord;
in vec3 lighting;

out vec4 fragColor;

void main(void) {
	vec4 texColor = texture2D(u_texture, texCoord);
	texColor.rgb += u_flash;
	fragColor = vec4(texColor.rgb * lighting, u_alpha);
}