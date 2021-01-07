#version 410

in vec2 PixelCoordinate;
uniform sampler2D TextAtlas;
uniform float CharWidthPixels;
uniform float CharHeightPixels;
uniform int Characters[10];
uniform int QuadsPerLine = 10;
uniform float HalfPixel = 1.0/1024.0;
out vec4 FragColor;

vec2 getAtlasOrigin(int atlasIndex)
{
  float i = mod(atlasIndex, QuadsPerLine);
  float j = floor(atlasIndex/QuadsPerLine);
  float quadSize = 1.0/QuadsPerLine;
  return vec2(i*quadSize, j*quadSize);
}

void main()
{
  int charIndex = int(floor(PixelCoordinate.x/CharWidthPixels));
  int atlasIndex = Characters[charIndex];
  if (atlasIndex >= 0) {
    vec2 atlasOrigin = getAtlasOrigin(atlasIndex);
    float pixelsFromLeft = mod(PixelCoordinate.x, CharWidthPixels);
    float atlasScale = (1.0 / QuadsPerLine) - 4*HalfPixel;
    vec2 atlasOffset = vec2(pixelsFromLeft/CharWidthPixels, PixelCoordinate.y/CharHeightPixels) * atlasScale;
    vec2 atlasCoords = atlasOrigin + atlasOffset;
    FragColor = texture(TextAtlas, atlasCoords);
  } else {
    FragColor = vec4(0);
  }
}
