#version 410

layout (location = 0) in vec3 VertexPosition;
layout (location = 1) in int NormalIndex;
layout (location = 2) in vec2 VertexUv;

out vec3 LightIntensity;
out vec2 TexCoords;

// normals
vec3 normals[] = vec3[](
    vec3(0,1,0), // north
    vec3(1,0,0), // east
    vec3(0,-1,0), // south
    vec3(-1,0,0), // west
    vec3(0,0,1), // top
    vec3(0,0,-1) // bottom
);

// vec3 North = vec3(0,1,0);
// vec3 South = vec3(0,-1,0);
// vec3 East = vec3(1,0,0);
// vec3 West = vec3(-1,0,0);
// vec3 Top = vec3(0,0,1);
// vec3 Bottom = vec3(0,0,-1);

// uniform vec4 LightPosition = vec4(3.0, 3.0, 30.0, 1.0);
uniform vec3 LightDirection = normalize(vec3(4.0, 10.0, -20.0));
uniform vec3 Kd = vec3(1.0);
uniform vec3 Ld = vec3(1.0);
// uniform vec4 LightPosition; // Light position in eye coords.
// uniform vec3 Kd;            // Diffuse reflectivity
// uniform vec3 Ld;            // Diffuse light intensity

// uniform mat4 ModelViewMatrix;
// uniform mat3 NormalMatrix; // inverse transpose of modelview matrix
uniform mat4 MVP;

void main()
{
    // vec3 tnorm = normalize( NormalMatrix * VertexNormal );
    // vec4 eyeCoords = ModelViewMatrix * vec4(VertexPosition,1.0);
    // vec3 s = normalize(vec3(LightPosition - eyeCoords));
    vec3 Normal = normals[NormalIndex];
    LightIntensity = Ld * Kd * (0.5*dot( -LightDirection, Normal ) + 0.5);
    TexCoords = VertexUv;
    gl_Position = MVP * vec4(VertexPosition, 1.0);
}