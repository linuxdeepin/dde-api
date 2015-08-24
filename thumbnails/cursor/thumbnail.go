package cursor

import (
	"os"
	"path"
	"pkg.deepin.io/dde/api/thumbnails/loader"
	"pkg.deepin.io/lib/graphic"
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

	err = graphic.ThumbnailImage(tmp, dest, width, height, graphic.FormatPng)
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
	var (
		x   = defaultPointX
		y   = defaultPointY
		tmp = bg
	)
	for _, file := range files {
		err := graphic.CompositeImage(tmp, file, dest, x, y, graphic.FormatPng)
		if err != nil {
			return err
		}
		tmp = dest
		x += defaultPointX + defaultIconSize
	}
	return nil
}
