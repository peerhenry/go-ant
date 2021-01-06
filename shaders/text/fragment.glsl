#version 410

in vec2 PixelCoordinate;
uniform sampler2D TextAtlas;
uniform float AtlasWidthPixels;
uniform float CharWidthPixels;
uniform int Characters[10];
uniform int QuadsPerLine = 10;
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
  vec4 col;
  if (atlasIndex >= 0) {
    vec2 atlasOrigin = getAtlasOrigin(atlasIndex);
    float pixelLeft = mod(PixelCoordinate.x, CharWidthPixels);
    vec2 atlasCoords = atlasOrigin + vec2(pixelLeft/AtlasWidthPixels, PixelCoordinate.y/AtlasWidthPixels);
    col = texture(TextAtlas, atlasCoords) + vec4(0.5,0,0,0);
  } else {
    col = vec4(0);
  }
  FragColor = col + vec4(0.2,0.2,0.2,1);
}
