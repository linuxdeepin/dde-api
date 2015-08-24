package icon

import (
	"path"
	"pkg.deepin.io/dde/api/thumbnails/images"
	"pkg.deepin.io/dde/api/thumbnails/loader"
	"pkg.deepin.io/lib/graphic"
	"strings"
)

const (
	presentIconFolder      = "folder"
	presentIconTrash       = "user-trash"
	presentIconFullTrash   = "user-trash-full"
	presentIconFilemanager = "system-file-manager"
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
	err := compositeImages(bg, tmp, getIconFiles(path.Base(dir)))
	if err != nil {
		return "", err
	}

	err = graphic.ThumbnailImage(tmp, dest, width, height, graphic.FormatPng)
	if err != nil {
		return "", err
	}

	return dest, nil
}

func getIconFiles(theme string) []string {
	var files []string
	files = append(files, GetIconFile(theme, presentIconFolder))
	trash := GetIconFile(theme, presentIconTrash)
	if len(trash) == 0 {
		trash = GetIconFile(theme, presentIconFullTrash)
	}
	files = append(files, trash)
	files = append(files, GetIconFile(theme, presentIconFilemanager))
	return convertSvgFiles(files)
}

func convertSvgFiles(files []string) []string {
	var ret []string
	for _, file := range files {
		if !strings.HasSuffix(file, ".svg") {
			ret = append(ret, file)
			continue
		}
		tmp, err := images.GenThumbnail(file, defaultIconSize, defaultIconSize)
		if err != nil {
			return nil
		}
		ret = append(ret, tmp)
	}

	return ret
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
