#version 410

in vec2 PixelCoordinate;
out vec4 FragColor;

void main()
{
  FragColor = vec4(1,min(PixelCoordinate.x, 0.0),0,1);
}
