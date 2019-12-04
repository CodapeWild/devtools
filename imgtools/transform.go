package imgtools

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

type Transform struct {
	img    image.Image
	format string
}

type TransformSetting func(*Transform) error

func SetResourceFromFile(filePath string) TransformSetting {
	return func(trans *Transform) error {
		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		if trans.img, trans.format, err = image.Decode(f); err != nil {
			return err
		}

		return nil
	}
}

func SetResourceFromUrl(url string) TransformSetting {
	return func(trans *Transform) error {
		return nil
	}
}

func SetResourceFromBuffer(buf []byte) TransformSetting {
	return func(trans *Transform) error {
		var err error
		if trans.img, trans.format, err = image.Decode(bytes.NewBuffer(buf)); err != nil {
			return err
		}

		return nil
	}
}

func NewTransform(opt ...TransformSetting) (*Transform, error) {
	trans := &Transform{}
	for _, v := range opt {
		if err := v(trans); err != nil {
			return nil, err
		}
	}

	return trans, nil
}

func (this *Transform) Zoom(x, y float64) *Transform {
	return this
}

func (this *Transform) Decolorize() *Transform {
	return this
}
