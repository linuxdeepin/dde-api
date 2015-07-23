// Image thumbnail generator
package images

import (
	"fmt"
	. "pkg.deepin.io/dde/api/thumbnails/loader"
	"pkg.deepin.io/lib/mime"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	ImageTypePng  string = "image/png"
	ImageTypeJpeg        = "image/jpeg"
	ImageTypeGif         = "image/gif"
	ImageTypeBmp         = "image/bmp"
	ImageTypeTiff        = "image/tiff"
	ImageTypeSvg         = "image/svg+xml"
)

func init() {
	for _, ty := range SupportedTypes() {
		Register(ty, GenThumbnail)
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

func GenThumbnail(src, bg string, width, height int) (string, error) {
	ty, err := mime.Query(src)
	if err != nil {
		return "", err
	}

	if !IsStrInList(ty, SupportedTypes()) {
		return "", fmt.Errorf("No supported type: %v", ty)
	}

	src = dutils.DecodeURI(src)
	dest, err := GetThumbnailDest(src, width, height)
	if err != nil {
		return "", err
	}
	if dutils.IsFileExist(dest) {
		return dest, nil
	}

	switch ty {
	case ImageTypeSvg:
		tmp := GetTmpImage()
		err := svgToPng(src, tmp)
		if err != nil {
			return "", err
		}
		src = tmp
	}

	err = Scale(src, dest, width, height)
	if err != nil {
		return "", err
	}
	return dest, nil
}
