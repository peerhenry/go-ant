#version 410

layout (location = 0) in int index;

// uniform vec2 viewDimensionsPixels;
// uniform vec2 marginPixels;
// uniform vec2 lineDimensionsPixels;
uniform vec2 DimensionsPixels;
uniform vec2 Dimensions; // = 2*lineDimensionsPixels/viewDimensionsPixels - 1
uniform vec2 Margin; // = 2*marginPixels/viewDimensionsPixels - 1

uniform vec4 topleft = vec4(-1.0, 1.0, 0.0, 1.0);
uniform vec4 topright = vec4(1.0, 1.0, 0.0, 1.0);
uniform vec4 bottomleft = vec4(-1.0, -1.0, 0.0, 1.0);
uniform vec4 bottomright = vec4(1.0, -1.0, 0.0, 1.0);

out vec2 PixelCoordinate;

void main() {
  switch(index) {
    case 1:
      gl_Position = 0.001*vec4(Margin, 0.0, 1.0) + topleft;
      PixelCoordinate = vec2(0.0, 0.0);
      break;
    case 2:
      gl_Position = 0.001*vec4(Margin.x, Margin.y + Dimensions.y, 0.0, 1.0) + bottomleft;
      PixelCoordinate = vec2(0.0, DimensionsPixels.y);
      break;
    case 3:
      gl_Position = 0.001*vec4(Margin.x + Dimensions.y, Margin.y, 0.0, 1.0) + topright;
      PixelCoordinate = vec2(DimensionsPixels.x, 0.0);
      break;
    case 4:
      gl_Position = 0.001*vec4(Margin + Dimensions, 0.0, 1.0) + bottomright;
      PixelCoordinate = DimensionsPixels;
      break;
  }
}