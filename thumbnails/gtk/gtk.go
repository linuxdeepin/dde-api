// Gtk theme thumbnail generator
package gtk

import (
	"fmt"
	"os"
	"path"
	. "pkg.deepin.io/dde/api/thumbnails/loader"
	"pkg.deepin.io/lib/mime"
	dutils "pkg.deepin.io/lib/utils"
)

func init() {
	for _, ty := range SupportedTypes() {
		Register(ty, genGtkThumbnail)
	}
}

func SupportedTypes() []string {
	return []string{
		mime.MimeTypeGtk,
	}
}

func GenThumbnail(src, bg string, width, height int, force bool) (string, error) {
	if width <= 0 || height <= 0 {
		return "", fmt.Errorf("Invalid width or height")
	}

	ty, err := mime.Query(src)
	if err != nil {
		return "", err
	}
	if ty != mime.MimeTypeGtk {
		return "", fmt.Errorf("Unspported mime: %s", ty)
	}

	return genGtkThumbnail(src, bg, width, height, force)
}

func genGtkThumbnail(src, bg string, width, height int, force bool) (string, error) {
	dest, err := GetThumbnailDest(src, width, height)
	if err != nil {
		return "", err
	}
	if !force && dutils.IsFileExist(dest) {
		return dest, nil
	}

	if len(bg) == 0 {
		bg, err = GetBackground(width, height)
		if err != nil {
			return "", err
		}
		defer os.Remove(bg)
	} else {
		dutils.DecodeURI(bg)
	}

	return doGenThumbnail(path.Base(path.Dir(src)), dest, bg,
		width, height)
}
