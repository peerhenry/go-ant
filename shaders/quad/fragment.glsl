
#version 410

in vec2 TexCoords;
uniform Sampler2D Tex;
out vec4 FragColor;

void main() {
  vec4 texColor = texture(Tex, TexCoords);
  FragColor = texColor;
}