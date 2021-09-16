package icon

/*
#cgo pkg-config: gtk+-3.0 gdk-pixbuf-2.0
#include <gtk/gtk.h>
#include <gdk-pixbuf/gdk-pixbuf.h>
#include <glib/gprintf.h>
#include <stdlib.h>

static char* choose_icon(const char *theme_name, const char **icon_names, int size) {
	if (!gtk_init_check(NULL, NULL)) {
		g_warning("Init gtk environment failed");
		return FALSE;
	}
	GtkIconTheme *icon_theme = gtk_icon_theme_new();
	gtk_icon_theme_set_custom_theme(icon_theme, theme_name);

	GtkIconInfo* icon_info = gtk_icon_theme_choose_icon(icon_theme, icon_names, size, 0);
	if (icon_info == NULL ) {
		g_printf("gtk_icon_theme_choose_icon failed theme_name: %s, size: %d\n", theme_name, size);
		g_object_unref(icon_theme);
		return NULL;
	}
	const gchar* filename = gtk_icon_info_get_filename(icon_info);
	gchar* filename_dup = g_strdup(filename);

	g_object_unref(icon_info);
	g_object_unref(icon_theme);

	return filename_dup;
}

*/
import "C"
import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"unsafe"

	"github.com/nfnt/resize"
	"pkg.deepin.io/dde/api/theme_thumb/common"
)

const (
	Version      = 1
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
	{"google-chrome", "firefox", "chromium", "opera", "internet-web-browser", "browser"},
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
func chooseIcon(theme string, iconNames []string, size int) string {
	cTheme := C.CString(theme)
	cArr := cStrv(iconNames)
	cNames := (**C.char)(unsafe.Pointer(&cArr[0]))

	cFilename := C.choose_icon(cTheme, cNames, C.int(size))
	filename := C.GoString(cFilename)

	C.free(unsafe.Pointer(cFilename))
	C.free(unsafe.Pointer(cTheme))
	// free cArr
	for i := range cArr {
		C.free(unsafe.Pointer(cArr[i]))
	}
	return filename
}

func ChooseIcon(theme string, iconNames []string, size int) string {
	return chooseIcon(theme, iconNames, size)
}

func loadIcon(theme string, iconNames []string, size int) (img image.Image, err error) {
	filename := chooseIcon(theme, iconNames, size)
	if filename == "" {
		return nil, errors.New("failed to choose icon")
	}
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == ".svg" {
		img, err = loadSvg(filename, size)
	} else {
		img, err = loadOther(filename)
	}

	if err != nil {
		return nil, err
	}
	imgWidth := img.Bounds().Dx()
	if imgWidth != size {
		img = resize.Resize(uint(size), 0, img, resize.Bilinear)
	}
	return img, nil
}

func loadOther(filename string) (image.Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()
	img, _, err := image.Decode(f)
	return img, err
}

func loadSvg(filename string, size int) (img image.Image, err error) {
	sizeStr := strconv.Itoa(size)
	cmd := exec.Command("rsvg-convert", "-f", "png", "-w", sizeStr, "-h", sizeStr, filename)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		return
	}
	defer func() {
		err = cmd.Wait()
	}()
	return png.Decode(stdout)
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
			fmt.Fprintf(os.Stderr, "failed to load icon %s %v: %v\n", theme, iconNames, err)
			continue
		}

		images = append(images, img)
		if len(images) == 6 {
			break
		}
	}
	return
}
