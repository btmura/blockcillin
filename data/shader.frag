#version 330 core

uniform sampler2D u_texture;
uniform float u_grayscale;
uniform float u_brightness;
uniform float u_alpha;

uniform vec3 u_mixColor;
uniform float u_mixAmount;

in vec2 texCoord;
in vec3 lighting;

out vec4 fragColor;

void main(void) {
	vec4 color = texture2D(u_texture, texCoord);
	color.rgb += u_brightness;
	color.a *= u_alpha;

	color = vec4(mix(color.rgb, u_mixColor, u_mixAmount), color.a);

	vec3 grayColor = vec3(color.r * 0.21 + color.g * 0.72 + color.b * 0.07);
	color = vec4(mix(color.rgb, grayColor, u_grayscale), color.a);

	fragColor = vec4(color.rgb * lighting, color.a);
}