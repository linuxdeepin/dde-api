package images

import (
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"runtime/debug"
)

func Scale(src, dest string, width, height int) error {
	img, err := imaging.Open(src)
	if err != nil {
		return err
	}
	defer debug.FreeOSMemory()

	//thumb := imaging.Thumbnail(img, width, height, imaging.Box)
	thumb := imaging.Fit(img, width, height, imaging.Lanczos)

	dst := imaging.New(width, height, color.NRGBA{0, 0, 0, 0})
	dst = imaging.Paste(dst, thumb, image.Pt(0, 0))

	return imaging.Save(dst, dest)
}
