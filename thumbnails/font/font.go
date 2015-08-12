package font

import (
	"fmt"
	. "pkg.deepin.io/dde/api/thumbnails/loader"
	"pkg.deepin.io/lib/mime"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	FontTypeTTF = "application/x-font-ttf"
	FontTypeOTF = "application/vnd.ms-opentype"
)

func init() {
	for _, ty := range SupportedTypes() {
		Register(ty, genFontThumbnail)
	}
}

func SupportedTypes() []string {
	return []string{
		FontTypeOTF,
		FontTypeTTF,
	}
}

func GenThumbnail(src string, width, height int) (string, error) {
	if width <= 0 || height <= 0 {
		return "", fmt.Errorf("Invalid width or height")
	}

	ty, err := mime.Query(src)
	if err != nil {
		return "", err
	}

	if !IsStrInList(ty, SupportedTypes()) {
		return "", fmt.Errorf("Not supported type: %v", ty)
	}

	return genFontThumbnail(src, "", width, height)
}

func genFontThumbnail(src, bg string, width, height int) (string, error) {
	dest, err := GetThumbnailDest(src, width, height)
	if err != nil {
		return "", err
	}
	if dutils.IsFileExist(dest) {
		return dest, nil
	}

	return doGenThumbnail(dutils.DecodeURI(src), dest, width, height)
}
