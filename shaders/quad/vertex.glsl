#version 410

in vec2 pos;
in vec2 uv;

out vec2 TexCoords;

void main() {
  TexCoords = uv;
  gl_Position = vec3(pos,1.0);
}