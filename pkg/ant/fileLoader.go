package ant

import (
	"bufio"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ReadFile(path string) string {
	dat, err := ioutil.ReadFile(path)
	check(err)
	return string(dat)
}

func ReadFileAsByteBuffer(path string) *bufio.Reader {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	bytes := make([]byte, size)
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)
	check(err)
	return buffer
}

func ReadImage(path string) image.Image {
	log.Println("os.Open", path)
	file, err := os.Open(path)
	defer file.Close()
	check(err)
	log.Println("decoding file")
	image, _, err := image.Decode(file)
	check(err)
	return image
}
