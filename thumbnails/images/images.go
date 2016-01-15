// Image thumbnail generator
package images

import (
	"fmt"
	"os"
	"path"

	. "pkg.deepin.io/dde/api/thumbnails/loader"
	"pkg.deepin.io/lib/mime"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	sysThemeThumbDir = "/var/cache/thumbnails/appearance"
)

const (
	ImageTypePng  string = "image/png"
	ImageTypeJpeg        = "image/jpeg"
	ImageTypeGif         = "image/gif"
	ImageTypeBmp         = "image/bmp"
	ImageTypeTiff        = "image/tiff"
	ImageTypeSvg         = "image/svg+xml"
)

var themeThumbDir = path.Join(os.Getenv("HOME"),
	".cache", "thumbnails", "appearance")

func init() {
	for _, ty := range SupportedTypes() {
		switch ty {
		case ImageTypeSvg:
			Register(ty, genSvgThumbnail)
		default:
			Register(ty, genImageThumbnail)
		}
	}
}

func SupportedTypes() []string {
	return []string{
		ImageTypePng,
		ImageTypeJpeg,
		ImageTypeGif,
		ImageTypeBmp,
		ImageTypeTiff,
		ImageTypeSvg,
	}
}

func GenThumbnail(src string, width, height int, force bool) (string, error) {
	if width <= 0 || height <= 0 {
		return "", fmt.Errorf("Invalid width or height")
	}

	ty, err := mime.Query(src)
	if err != nil {
		return "", err
	}

	if !IsStrInList(ty, SupportedTypes()) {
		return "", fmt.Errorf("No supported type: %v", ty)
	}

	switch ty {
	case ImageTypeSvg:
		return genSvgThumbnail(src, "", width, height, force)
	}

	return genImageThumbnail(src, "", width, height, force)
}

func ThumbnailForTheme(src string, width, height int, force bool) (string, error) {
	if width <= 0 || height <= 0 {
		return "", fmt.Errorf("Invalid width or height")
	}

	dest, err := getThumbDest(src, width, height, true)
	if err != nil {
		return "", err
	}

	thumb := path.Join(sysThemeThumbDir, path.Base(dest))
	if !force && dutils.IsFileExist(thumb) {
		return thumb, nil
	}

	return doGenThumbnail(src, "", dest, width, height, force, true)
}

func genSvgThumbnail(src, bg string, width, height int, force bool) (string, error) {
	tmp := GetTmpImage()
	err := svgToPng(src, tmp)
	if err != nil {
		return "", err
	}

	defer os.Remove(tmp)
	return genImageThumbnail(tmp, bg, width, height, force)
}

func genImageThumbnail(src, bg string, width, height int, force bool) (string, error) {
	dest, err := getThumbDest(src, width, height, false)
	if err != nil {
		return "", err
	}

	return doGenThumbnail(src, bg, dest, width, height, force, false)
}

func doGenThumbnail(src, bg, dest string, width, height int, force, theme bool) (string, error) {
	if !force && dutils.IsFileExist(dest) {
		return dest, nil
	}

	var err error
	src = dutils.DecodeURI(src)
	if !theme {
		err = ThumbnailImage(src, dest, width, height)
	} else {
		err = ScaleImage(src, dest, width, height)
	}
	if err != nil {
		return "", err
	}
	return dest, nil
}

func getThumbDest(src string, width, height int, theme bool) (string, error) {
	dest, err := GetThumbnailDest(src, width, height)
	if err != nil {
		return "", err
	}
	if theme {
		dest = path.Join(themeThumbDir, "bg-"+path.Base(dest))
	}
	return dest, nil
}
