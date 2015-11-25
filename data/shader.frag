#version 330 core

uniform sampler2D u_texture;
uniform float u_grayscale;
uniform float u_brightness;
uniform float u_alpha;

in vec2 texCoord;
in vec3 lighting;

out vec4 fragColor;

void main(void) {
	vec4 texColor = texture2D(u_texture, texCoord);
	texColor.rgb += u_brightness;

	vec3 grayColor = vec3(texColor.r * 0.21 + texColor.g * 0.72 + texColor.b * 0.07);

	vec3 color = mix(texColor.rgb, grayColor, u_grayscale);
	fragColor = vec4(color.rgb * lighting, u_alpha);
}