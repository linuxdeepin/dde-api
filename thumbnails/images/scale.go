package images

import (
	"pkg.deepin.io/lib/graphic"
)

func Scale(src, dest string, width, height int) error {
	return graphic.ThumbnailImage(src, dest, width, height,
		graphic.FormatPng)
}
