package service

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

type JPEGEncoder struct{}

func (j JPEGEncoder) Encode(w io.Writer, i image.Image) error {
	return jpeg.Encode(w, i, &jpeg.Options{Quality: 20})
}

type PNGEncoder struct{}

func (p PNGEncoder) Encode(w io.Writer, i image.Image) error {
	return png.Encode(w, i)
}
