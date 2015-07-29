//Cursor theme thumbnail generator
package cursor

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
		mime.MimeTypeCursor,
	}
}

// GenThumbnail generate cursor theme thumbnail
// src: the uri of cursor theme index.theme
func GenThumbnail(src, bg string, width, height int) (string, error) {
	ty, err := mime.Query(src)
	if err != nil {
		return "", err
	}

	if ty != mime.MimeTypeCursor {
		return "", fmt.Errorf("Not supported type: %v", ty)
	}

	src = dutils.DecodeURI(src)
	dir := path.Dir(src)
	dest, err := GetThumbnailDest(path.Join(dir, "cursors", "left_ptr"),
		width, height)
	if err != nil {
		return "", err
	}
	if dutils.IsFileExist(dest) {
		return dest, nil
	}
	return doGenThumbnail(dir, dest, dutils.DecodeURI(bg), width, height)
}
