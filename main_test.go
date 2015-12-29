package main

import (
	"image/color"
	"testing"
)

func TestBlendColors(t *testing.T) {
	t.Parallel()
	c1 := color.RGBA{255, 255, 0, 1}
	c2 := color.RGBA{255, 0, 0, 1}
	c3 := BlendColors(c1, c2, 0.5)

	if c3.R != 255 || c3.G != 180 || c3.B != 0 || c3.A != 1 {
		t.Error("Color blending did not work correctly")
	}
}
