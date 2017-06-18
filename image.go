package artsparts

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
)

func imageToBaseString(img image.Image) (string, error) {
	buf := &bytes.Buffer{}
	err := jpeg.Encode(buf, img, nil)
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	return encoded, err
}
