#version 410

layout(location = 0)in vec3 vert;

void main() {
  gl_Position = vec4(vert, 0.0);
}