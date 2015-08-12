package thumbnails

import (
	"fmt"
	_ "pkg.deepin.io/dde/api/thumbnails/cursor"
	_ "pkg.deepin.io/dde/api/thumbnails/font"
	_ "pkg.deepin.io/dde/api/thumbnails/gtk"
	_ "pkg.deepin.io/dde/api/thumbnails/icon"
	_ "pkg.deepin.io/dde/api/thumbnails/images"
	"pkg.deepin.io/dde/api/thumbnails/loader"
	_ "pkg.deepin.io/dde/api/thumbnails/pdf"
	_ "pkg.deepin.io/dde/api/thumbnails/text"
	"pkg.deepin.io/lib/mime"
)

func GenThumbnail(uri string, size int) (string, error) {
	if size < 0 {
		return "", fmt.Errorf("Invalid size: '%v'", size)
	}

	ty, err := mime.Query(uri)
	if err != nil {
		return "", err
	}

	return GenThumbnailWithMime(uri, ty, size)
}

func GenThumbnailWithMime(uri, ty string, size int) (string, error) {
	if size < 0 {
		return "", fmt.Errorf("Invalid size: '%v'", size)
	}

	handler, err := loader.GetHandler(ty)
	if err != nil {
		return "", err
	}

	return handler(uri, "", size, size)
}
