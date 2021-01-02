#version 410

layout (location = 0) in vec2 pos;
layout (location = 1) in vec2 uv;

out vec2 TexCoords;

void main() {
  TexCoords = uv;
  gl_Position = vec4(pos, 0.0, 1.0);
}