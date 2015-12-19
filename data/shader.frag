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
	texColor.a *= u_alpha;

	vec4 grayColor = vec4(
		vec3(texColor.r * 0.21 + texColor.g * 0.72 + texColor.b * 0.07),
		texColor.a);

	vec4 mixedColor = mix(texColor, grayColor, u_grayscale);

	fragColor = vec4(mixedColor.rgb * lighting, mixedColor.a);
}