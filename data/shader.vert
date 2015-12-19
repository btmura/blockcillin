#version 330 core

uniform mat4 u_projectionViewMatrix;
uniform mat4 u_modelMatrix;
uniform mat4 u_normalMatrix;

uniform vec3 u_ambientLightColor;
uniform vec3 u_directionalLightColor;
uniform vec3 u_directionalVector;

layout (location = 0) in vec4 i_position;
layout (location = 1) in vec4 i_normal;
layout (location = 2) in vec2 i_texCoord;

out vec2 texCoord;
out vec3 lighting;

void main(void) {
	gl_Position = u_projectionViewMatrix * u_modelMatrix * i_position;

	texCoord = i_texCoord;

	vec4 transformedNormal = u_normalMatrix * vec4(i_normal.xyz, 1.0);
	float directional = max(dot(transformedNormal.xyz, u_directionalVector), 0.0);
	lighting = u_ambientLightColor + (u_directionalLightColor * directional);
}