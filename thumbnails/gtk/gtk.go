// Gtk theme thumbnail generator
package gtk

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
		mime.MimeTypeGtk,
	}
}

func GenThumbnail(src, bg string, width, height int) (string, error) {
	ty, err := mime.Query(src)
	if err != nil {
		return "", err
	}
	if ty != mime.MimeTypeGtk {
		return "", fmt.Errorf("Unspported mime: %s", ty)
	}

	dest, err := GetThumbnailDest(src, width, height)
	if err != nil {
		return "", err
	}
	if dutils.IsFileExist(dest) {
		return dest, nil
	}

	if len(bg) == 0 {
		bg, err = GetBackground(width, height)
		if err != nil {
			return "", err
		}
	} else {
		dutils.DecodeURI(bg)
	}

	return doGenThumbnail(path.Base(path.Dir(src)), dest, bg,
		width, height)
}
