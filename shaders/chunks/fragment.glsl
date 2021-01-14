#version 410

in vec3 LightIntensity;
in vec2 TexCoords;
uniform sampler2D Tex;
layout( location = 0 ) out vec4 FragColor;

void main()
{
  vec4 texColor = texture(Tex, TexCoords);
  vec3 color = LightIntensity * vec3(texColor);
  FragColor = vec4(color, 1.0);
}
