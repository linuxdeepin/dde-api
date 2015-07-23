//Icon theme thumbnail generator
package icon

import (
	"fmt"
	"path"
	. "pkg.deepin.io/dde/api/thumbnails/loader"
	"pkg.deepin.io/lib/mime"
	dutils "pkg.deepin.io/lib/utils"
)

func init() {
	for _, ty := range SupportedTypes() {
		Register(ty, GenThumbnail)
	}
}

func SupportedTypes() []string {
	return []string{
		mime.MimeTypeIcon,
	}
}

// GenThumbnail generate icon theme thumbnail
// src: the uri of icon theme index.theme
func GenThumbnail(src, bg string, width, height int) (string, error) {
	ty, err := mime.Query(src)
	if err != nil {
		return "", err
	}

	if ty != mime.MimeTypeIcon {
		return "", fmt.Errorf("Not supported type: %v", ty)
	}

	src = dutils.DecodeURI(src)
	dest, err := GetThumbnailDest(src, width, height)
	if err != nil {
		return "", err
	}
	if dutils.IsFileExist(dest) {
		return dest, nil
	}
	return doGenThumbnail(path.Dir(src), dest, dutils.DecodeURI(bg), width, height)
}
