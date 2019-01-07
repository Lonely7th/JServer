package util

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
)

// 读取图片
func LoadImage(path string) (img image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	return
}

// 保存Png图片
func SaveImage(path string, img image.Image) (err error) {
	// 需要保存的文件
	imgfile, err := os.Create(path)
	defer imgfile.Close()

	// 以PNG格式保存文件
	err = png.Encode(imgfile, img)
	if err != nil {
		log.Fatal(err)
	}
	return
}
