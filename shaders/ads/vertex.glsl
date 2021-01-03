#version 410

layout (location = 0) in vec3 VertexPosition;
layout (location = 1) in vec3 VertexNormal;
layout (location = 2) in vec2 VertexUv;

out vec3 LightIntensity;
out vec2 TexCoords;

uniform vec4 LightPosition = vec4(3.0, 3.0, 30.0, 1.0);
uniform vec3 Kd = vec3(1.0);
uniform vec3 Ld = vec3(1.0);
// uniform vec4 LightPosition; // Light position in eye coords.
// uniform vec3 Kd;            // Diffuse reflectivity
// uniform vec3 Ld;            // Diffuse light intensity

uniform mat4 ModelViewMatrix;
uniform mat3 NormalMatrix; // inverse transpose of modelview matrix
uniform mat4 ProjectionMatrix;
uniform mat4 MVP;

void main()
{
    vec3 tnorm = normalize( NormalMatrix * VertexNormal );
    vec4 eyeCoords = ModelViewMatrix * vec4(VertexPosition,1.0);
    vec3 s = normalize(vec3(LightPosition - eyeCoords));
    LightIntensity = Ld * Kd * max( dot( s, tnorm ), 0.2 );
    TexCoords = VertexUv;
    gl_Position = MVP * vec4(VertexPosition, 1.0);
}