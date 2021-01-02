package ant

import (
	"image"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func loadTexture(rgba *image.NRGBA) uint32 {
	var textureId uint32
	gl.GenTextures(1, &textureId)
	gl.BindTexture(gl.TEXTURE_2D, textureId)
	gl.TexImage2D(
		gl.TEXTURE_2D, // target
		0,             // samples
		gl.RGBA,       // internalFormat
		int32(rgba.Bounds().Dx()),
		int32(rgba.Bounds().Dy()),
		0,                // border
		gl.RGBA,          // format
		gl.UNSIGNED_BYTE, // xtype
		gl.Ptr(rgba.Pix), // pixels (unsafe pointer)
	)
	return textureId
}
