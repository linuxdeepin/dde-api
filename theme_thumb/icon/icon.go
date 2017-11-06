package icon

/*
#cgo pkg-config: gtk+-3.0 gdk-pixbuf-2.0
#include <gtk/gtk.h>
#include <gdk-pixbuf/gdk-pixbuf.h>
#include <glib/gprintf.h>
#include <stdlib.h>

static gboolean choose_icon(const char *theme_name, const char **icon_names, int size, const char* dest_filename) {
	if (!gtk_init_check(NULL, NULL)) {
		g_warning("Init gtk environment failed");
		return FALSE;
	}
	GtkIconTheme *icon_theme = gtk_icon_theme_new();
	gtk_icon_theme_set_custom_theme(icon_theme, theme_name);

	GtkIconInfo* icon_info = gtk_icon_theme_choose_icon(icon_theme, icon_names, size, 0);
	if (icon_info == NULL ) {
		g_printf("gtk_icon_theme_choose_icon failed icon_theme: %s, size: %d\n", icon_theme, size);
		return FALSE;
	}

	GError *err = NULL;
	GdkPixbuf* pixbuf = gtk_icon_info_load_icon(icon_info, &err);
	if (err != NULL ) {
		g_printf("err msg: %s\n", err->message);
		g_error_free(err);

		g_object_unref(icon_info);
		g_object_unref(icon_theme);
		return FALSE;
	}

	int pb_width = gdk_pixbuf_get_width(pixbuf);
	if (pb_width != size ) {
		// need scale
		GdkPixbuf* tmp = gdk_pixbuf_scale_simple(pixbuf, size, size, GDK_INTERP_BILINEAR);
		g_object_unref(pixbuf);
		pixbuf = tmp;
	}

	gboolean ok = gdk_pixbuf_save(pixbuf, dest_filename, "png", &err, NULL);
	if (!ok) {
		g_printf("err msg: %s\n", err->message);
		g_error_free(err);

		g_object_unref(pixbuf);
		g_object_unref(icon_info);
		g_object_unref(icon_theme);
		return FALSE;
	}

	g_object_unref(pixbuf);
	g_object_unref(icon_info);
	g_object_unref(icon_theme);
	return TRUE;
}
*/
import "C"
import (
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"unsafe"

	"pkg.deepin.io/dde/api/theme_thumb/common"
)

const (
	Version      = 0
	baseIconSize = 48
	basePadding  = 4
)

// default present icons
var presentIcons = [][]string{
	// file manager:
	{"dde-file-manager", "system-file-manager"},
	// music player:
	{"deepin-music", "banshee", "amarok", "deadbeef", "clementine", "rhythmbox"},
	// image viewer:
	{"deepin-image-viewer", "eog", "gthumb", "gwenview", "gpicview", "showfoto", "phototonic"},
	// media/video player:
	{"deepin-movie", "media-player", "totem", "smplayer", "vlc", "dragonplayer", "kmplayer"},
	// web browser:
	{"google-chrome", "firefox", "chromium", "opear", "internet-web-browser", "web-browser", "browser"},
	// system settings:
	{"preferences-system"},
	// text editor:
	{"accessories-text-editor", "text-editor", "gedit", "kedit", "xfce-edit"},
	// terminal:
	{"deepin-terminal", "utilities-terminal", "terminal", "gnome-terminal", "xfce-terminal", "terminator", "openterm"},
}

func Gen(theme string, width, height int, scaleFactor float64, out string) error {
	iconSize := int(baseIconSize * scaleFactor)
	padding := int(basePadding * scaleFactor)
	width = int(float64(width) * scaleFactor)
	height = int(float64(height) * scaleFactor)

	images := getIcons(theme, iconSize)
	ret := common.CompositeIcons(images, width, height, iconSize, padding)
	return common.SavePngFile(ret, out)
}

// wrap for choose_icon
func chooseIcon(theme string, iconNames []string, size int) (string, error) {
	f, err := ioutil.TempFile(os.TempDir(), "theme-thumb-icon-")
	if err != nil {
		return "", err
	}
	destFilename := f.Name()
	f.Close()

	cTheme := C.CString(theme)
	defer C.free(unsafe.Pointer(cTheme))

	cDestFilename := C.CString(destFilename)
	defer C.free(unsafe.Pointer(cDestFilename))

	cArr := cStrv(iconNames)
	cNames := (**C.char)(unsafe.Pointer(&cArr[0]))
	// TODO error handle
	ok := C.choose_icon(cTheme, cNames, C.int(size), cDestFilename)

	// free cArr
	for i := range cArr {
		C.free(unsafe.Pointer(cArr[i]))
	}

	if ok == 0 {
		// fail
		return "", errors.New("choose icon failed")
	}
	return destFilename, nil
}

func ChooseIcon(theme string, iconNames []string, size int) (string, error) {
	return chooseIcon(theme, iconNames, size)
}

func loadIcon(theme string, iconNames []string, size int) (image.Image, error) {
	filename, err := chooseIcon(theme, iconNames, size)
	if err != nil {
		return nil, err
	}
	defer os.Remove(filename)
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// return NUL-Terminated slice of C String
func cStrv(strv []string) []*C.char {
	cArr := make([]*C.char, len(strv)+1)
	for i, str := range strv {
		cArr[i] = C.CString(str)
	}
	return cArr
}

func getIcons(theme string, size int) (images []image.Image) {
	for _, iconNames := range presentIcons {
		img, err := loadIcon(theme, iconNames, size)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to load icon %v", iconNames)
			continue
		}

		images = append(images, img)
		if len(images) == 6 {
			break
		}
	}
	return
}
