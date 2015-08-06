package icon

import (
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"path"
	"pkg.deepin.io/dde/api/thumbnails/images"
	"pkg.deepin.io/dde/api/thumbnails/loader"
	"runtime/debug"
	"strings"
)

const (
	presentIconFolder      = "folder"
	presentIconTrash       = "user-trash"
	presentIconFullTrash   = "user-trash-full"
	presentIconFilemanager = "system-file-manager"
)

const (
	defaultWidth    = 192
	defaultHeight   = 108
	defaultIconSize = 48
	defaultPointX   = (defaultWidth - defaultIconSize*3) / 4
	defaultPointY   = (defaultHeight - defaultIconSize) / 2
)

func doGenThumbnail(dir, dest, bg string, width, height int) (string, error) {
	defer debug.FreeOSMemory()

	tmp := loader.GetTmpImage()
	err := compositeImages(bg, tmp, getIconFiles(path.Base(dir)))
	if err != nil {
		return "", err
	}

	err = images.Scale(tmp, dest, width, height)
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
