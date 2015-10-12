package cursor

import (
	"os"
	"path"

	"pkg.deepin.io/dde/api/thumbnails/loader"
	dutils "pkg.deepin.io/lib/utils"
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
)

func doGenThumbnail(src, bg string, width, height int, force, theme bool) (string, error) {
	src = dutils.DecodeURI(src)
	dest, err := getThumbDest(src, width, height, theme)
	if err != nil {
		return "", err
	}

	if !force && dutils.IsFileExist(dest) {
		return dest, nil
	}

	dir := path.Dir(src)
	tmp := loader.GetTmpImage()
	err = loader.CompositeIcons(getCursorIcons(dir), bg, tmp,
		defaultIconSize, defaultWidth, defaultHeight)
	os.RemoveAll(xcur2pngCache)
	if err != nil {
		return "", err
	}

	defer os.Remove(tmp)
	if !theme {
		err = loader.ThumbnailImage(tmp, dest, width, height)
	} else {
		err = loader.ScaleImage(tmp, dest, width, height)
	}
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
