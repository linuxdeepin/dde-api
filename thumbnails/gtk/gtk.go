// Gtk theme thumbnail generator
package gtk

import (
	"dbus/com/deepin/api/gtkthumbnailer"
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

var themeThumbDir = path.Join(os.Getenv("HOME"),
	".cache", "thumbnails", "appearance")

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

func ThumbnailForTheme(src, bg string, width, height int, force bool) (string, error) {
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

	return doGenThumbnail(path.Base(path.Dir(src)), bg, dest,
		width, height, force)
}

func genGtkThumbnail(src, bg string, width, height int, force bool) (string, error) {
	dest, err := getThumbDest(src, width, height, false)
	if err != nil {
		return "", err
	}

	return doGenThumbnail(path.Base(path.Dir(src)), bg, dest,
		width, height, force)
}

func getThumbDest(src string, width, height int, theme bool) (string, error) {
	dest, err := GetThumbnailDest(src, width, height)
	if err != nil {
		return "", err
	}

	if theme {
		dest = path.Join(themeThumbDir, "gtk-"+path.Base(dest))
	}
	return dest, nil
}

func doGenThumbnail(name, bg, dest string, width, height int, force bool) (string, error) {
	if !force && dutils.IsFileExist(dest) {
		return dest, nil
	}

	thumbnailer, err := gtkthumbnailer.NewGtkThumbnailer(
		"com.deepin.api.GtkThumbnailer",
		"/com/deepin/api/GtkThumbnailer",
	)
	if err != nil {
		return "", err
	}

	err = thumbnailer.Thumbnail(name, bg, dest, int32(width), int32(height), force)
	if err != nil {
		return "", err
	}
	return dest, nil
}
