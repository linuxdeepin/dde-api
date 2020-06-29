package cursor

import (
	"image"
	"path/filepath"

	"pkg.deepin.io/dde/api/theme_thumb/common"
)

const (
	Version      = 1
	basePadding  = 12
	baseIconSize = 24
)

func Gen(descFile string, width, height int, scaleFactor float64, out string) error {
	dir := filepath.Join(filepath.Dir(descFile), "cursors")

	iconSize := int(baseIconSize * scaleFactor)
	padding := int(basePadding * scaleFactor)
	width = int(float64(width) * scaleFactor)
	height = int(float64(height) * scaleFactor)

	images := getCursorIcons(dir, iconSize)
	ret := common.CompositeIcons(images, width, height, iconSize, padding)
	return common.SavePngFile(ret, out)
}

var presentCursors = [][]string{
	{"left_ptr"},
	{"left_ptr_watch"},
	{"x-cursor", "X_cursor"},
	{"hand2", "hand1"},
	{"grab", "grabbing", "closedhand"},
	{"fleur", "move"},
	{"sb_v_double_arrow"},
	{"sb_h_double_arrow"},
	{"watch", "wait"},
}

func getCursorIcons(dir string, size int) (images []image.Image) {
	for _, cursors := range presentCursors {
		for _, cursor := range cursors {
			img, err := loadXCursor(filepath.Join(dir, cursor), size)
			if err == nil {
				images = append(images, img)
				break
			}
		}
	}
	return
}
