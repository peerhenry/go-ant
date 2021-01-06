package ant

import (
	"image"
	"log"

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

// todo: refactor to uniform store
func LoadImageFileToUniform(filePath string, uniformName string, programHandle uint32) {
	log.Println("Reading texture atlas")
	i := ReadImage(filePath)
	switch i.(type) {
	case *image.Gray:
		panic("image was Gray instead of NRGBA")
	case *image.Gray16:
		panic("image was Gray16 instead of NRGBA")
	case *image.CMYK:
		panic("image was CMYK instead of NRGBA")
	case *image.Alpha:
		panic("image was Alpha instead of NRGBA")
	case *image.Alpha16:
		panic("image was Alpha16 instead of NRGBA")
	case *image.Paletted:
		panic("image was Paletted instead of NRGBA")
	case *image.RGBA64:
		panic("image was RGBA64 instead of NRGBA")
	case *image.RGBA:
		panic("image was RGBA instead of NRGBA")
	case *image.NRGBA64:
		panic("image was NRGBA64 instead of NRGBA")
	case *image.NRGBA:
		if nrgba, ok := i.(*image.NRGBA); ok {
			log.Println("image size:", nrgba.Bounds().Dx(), "x", nrgba.Bounds().Dy())
			loadTexture(nrgba)
			gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
			gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
			// set uniform
			texUniformLocation := gl.GetUniformLocation(programHandle, gl.Str(uniformName+"\x00"))
			gl.Uniform1i(texUniformLocation, 0)
		} else {
			panic("Could not extract NRGBA from image...")
		}
	default:
		panic("image type unknown")
	}
}
