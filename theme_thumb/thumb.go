package theme_thumb

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"pkg.deepin.io/dde/api/theme_thumb/cursor"
	"pkg.deepin.io/dde/api/theme_thumb/gtk"
	"pkg.deepin.io/dde/api/theme_thumb/icon"
	"pkg.deepin.io/lib/xdg/basedir"
)

var scaleFactor float64

const (
	width  = 320
	height = 70
)

var cacheDir = filepath.Join(basedir.GetUserCacheDir(), "deepin", "dde-api", "theme_thumb")

func getScaleDir() string {
	return fmt.Sprintf("X%.2f", scaleFactor)
}

func getTypeDir(type0 string, version int) string {
	return fmt.Sprintf("%s-v%d", type0, version)
}

func Init(scaleFactor0 float64) {
	scaleFactor = scaleFactor0
	removeUnusedScaleDirs()
	removeAllTypesOldVersionDirs()
}

func removeUnusedScaleDirs() {
	removeUnusedDirs(cacheDir+"/X*", getScaleDir())
}

func removeAllTypesOldVersionDirs() {
	scaleDir := getScaleDir()
	removeOldVersionDirs(scaleDir, "gtk", gtk.Version)
	removeOldVersionDirs(scaleDir, "cursor", cursor.Version)
	removeOldVersionDirs(scaleDir, "icon", icon.Version)
}

func removeOldVersionDirs(scaleDir, type0 string, version int) {
	pattern := filepath.Join(cacheDir, scaleDir, type0+"-v*")
	usedDir := getTypeDir(type0, version)
	removeUnusedDirs(pattern, usedDir)
}

func removeUnusedDirs(pattern string, usedDir string) {
	dirs, err := filepath.Glob(pattern)
	if err != nil {
		return
	}

	for _, dir := range dirs {
		name := filepath.Base(dir)
		if name != usedDir {
			fmt.Println("rm dir", dir)
			os.RemoveAll(dir)
		}
	}
}

func checkScaleFactor() error {
	if scaleFactor <= 0 {
		return errors.New("bad scale factor")
	}
	return nil
}

func GetCursor(id, descFile string) (string, error) {
	err := checkScaleFactor()
	if err != nil {
		return "", err
	}

	out := prepareOutputPath("cursor", id, cursor.Version)
	genNew, err := shouldGenerateNewCursor(descFile, out)
	if err != nil {
		return "", err
	}

	if !genNew {
		return out, nil
	}

	err = cursor.Gen(descFile, width, height, scaleFactor, out)
	if err != nil {
		return "", err
	}
	return out, nil
}

func GetGtk(id, descFile string) (string, error) {
	err := checkScaleFactor()
	if err != nil {
		return "", err
	}

	out := prepareOutputPath("gtk", id, gtk.Version)
	genNew, err := shouldGenerateNew(descFile, out)
	if err != nil {
		return "", err
	}

	if !genNew {
		return out, nil
	}

	err = gtk.Gen(id, width, height, scaleFactor, out)
	if err != nil {
		return "", err
	}
	return out, nil
}

func GetIcon(id, descFile string) (string, error) {
	err := checkScaleFactor()
	if err != nil {
		return "", err
	}

	out := prepareOutputPath("icon", id, icon.Version)
	genNew, err := shouldGenerateNew(descFile, out)
	if err != nil {
		return "", err
	}

	if !genNew {
		return out, nil
	}

	err = icon.Gen(id, width, height, scaleFactor, out)
	if err != nil {
		return "", err
	}
	return out, nil
}

func shouldGenerateNew(descFile, out string) (bool, error) {
	descFileInfo, err := os.Stat(descFile)
	if err != nil {
		return false, err
	}

	descFileCTime := getChangeTime(descFileInfo)

	thumbFileInfo, err := os.Stat(out)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, err
	}
	thumbFileCTime := getChangeTime(thumbFileInfo)

	if descFileCTime.After(thumbFileCTime) {
		return true, nil
	}
	return false, nil
}

func shouldGenerateNewCursor(descFile, out string) (bool, error) {
	dir := filepath.Dir(descFile)
	ptrFile := filepath.Join(dir, "cursors", "left_ptr")
	return shouldGenerateNew(ptrFile, out)
}

// getChangeTime get time when file status was last changed.
func getChangeTime(fileInfo os.FileInfo) time.Time {
	stat := fileInfo.Sys().(*syscall.Stat_t)
	return time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))
}

// ex. $HOME/.cache/deepin/dde-api/theme_thumb/X1.00/icon-v0/deepin.png
func prepareOutputPath(type0, id string, version int) string {
	scaleDir := getScaleDir()
	typeDir := getTypeDir(type0, version)
	dir := filepath.Join(cacheDir, scaleDir, typeDir)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return ""
	}
	return filepath.Join(dir, id+".png")
}
