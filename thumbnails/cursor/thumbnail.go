package cursor

import (
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"os"
	"path"
	"pkg.deepin.io/dde/api/thumbnails/images"
	"pkg.deepin.io/dde/api/thumbnails/loader"
)

const (
	presentCursorLeftPtr   = "left_ptr"
	presentCursorLeftWatch = "left_ptr_watch"
	presentCursorQuestion  = "question_arrow"
)

const (
	defaultWidth    = 128
	defaultHeight   = 72
	defaultIconSize = 24
	defaultPointX   = (defaultWidth - defaultIconSize*3) / 4
	defaultPointY   = (defaultHeight - defaultIconSize) / 2
)

func doGenThumbnail(dir, dest, bg string, width, height int) (string, error) {
	tmp := loader.GetTmpImage()
	err := compositeImages(bg, tmp, getCursorIcons(dir))
	os.RemoveAll(xcur2pngCache)
	if err != nil {
		return "", err
	}

	err = images.Scale(tmp, dest, width, height)
	if err != nil {
		return "", err
	}

	return dest, nil
}

func getCursorIcons(dir string) []string {
	presents := []string{
		presentCursorLeftPtr,
		presentCursorLeftWatch,
		presentCursorQuestion,
	}

	var files []string
	for _, name := range presents {
		tmp, err := XCursorToPng(path.Join(dir, "cursors", name))
		if err != nil {
			return nil
		}
		files = append(files, tmp)
	}
	return files
}

func compositeImages(bg, dest string, files []string) error {
	var dst image.Image
	if len(bg) != 0 {
		img, err := imaging.Open(bg)
		if err != nil {
			return err
		}
		dst = imaging.Fit(img, defaultWidth, defaultHeight,
			imaging.Lanczos)
	} else {
		dst = imaging.New(defaultWidth, defaultHeight,
			color.NRGBA{1, 1, 1, 1})
	}

	var (
		x = defaultPointX
		y = defaultPointY
	)
	for _, file := range files {
		img, err := imaging.Open(file)
		if err != nil {
			return err
		}

		dst = imaging.Paste(dst, img, image.Pt(x, y))
		x += defaultPointX + defaultIconSize
	}
	return imaging.Save(dst, dest)
}
