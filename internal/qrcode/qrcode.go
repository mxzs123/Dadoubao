package qrcode

import (
	"image"

	"github.com/skip2/go-qrcode"
)

// Generate 生成二维码图片
func Generate(url string, size int) (image.Image, error) {
	qr, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		return nil, err
	}
	return qr.Image(size), nil
}

// GeneratePNG 生成二维码 PNG 字节
func GeneratePNG(url string, size int) ([]byte, error) {
	return qrcode.Encode(url, qrcode.Medium, size)
}

