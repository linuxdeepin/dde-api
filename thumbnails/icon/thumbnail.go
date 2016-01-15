package icon

import (
	"os"
	"path"
	"strings"

	"pkg.deepin.io/dde/api/thumbnails/images"
	"pkg.deepin.io/dde/api/thumbnails/loader"
	dutils "pkg.deepin.io/lib/utils"
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
)

func doGenThumbnail(src, bg, dest string, width, height int, force, theme bool) (string, error) {
	if !force && dutils.IsFileExist(dest) {
		return dest, nil
	}

	src = dutils.DecodeURI(src)
	bg = dutils.DecodeURI(bg)
	dir := path.Dir(src)
	tmp := loader.GetTmpImage()
	err := loader.CompositeIcons(getIconFiles(path.Base(dir)), bg, tmp,
		defaultIconSize, defaultWidth, defaultHeight)
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
		tmp, err := images.GenThumbnail(file, defaultIconSize, defaultIconSize, true)
		if err != nil {
			return nil
		}
		ret = append(ret, tmp)
	}

	return ret
}
