package ant

import (
	"bufio"
	"fmt"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func readImageBytes(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	bytes := make([]byte, size)
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)
}

func LoadTexture(byteArray *[]byte, width int32, height int32) *uint32 {
	var textureId uint32
	gl.GenTextures(1, &textureId)
	gl.BindTexture(gl.TEXTURE_2D, textureId)
	// width, height = img.Bounds().Dx(), img.Bounds().Dy()
	gl.TexImage2D(
		gl.TEXTURE_2D, // target
		0,             // samples
		gl.RGBA,       // internalFormat
		width,
		height,
		0,                 // border
		gl.RGBA,           // format
		gl.UNSIGNED_BYTE,  // xtype
		gl.Ptr(byteArray), // pixels (unsafe pointer)
	)
	return &textureId
}
