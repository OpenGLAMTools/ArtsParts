package artsparts

import (
	"testing"

	"github.com/disintegration/imaging"
)

func TestImageToBaseString(t *testing.T) {
	img, err := imaging.Open("test/test.jpg")
	if err != nil {
		t.Error("Error opening test image")
	}
	s, err := ImageToBaseString(img)
	if err != nil {
		t.Error("imageToBaseString returns an error")
	}
	if s == "" {
		t.Error("String is empty and should return a value")
	}
}
