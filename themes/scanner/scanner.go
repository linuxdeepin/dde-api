// Theme scanner
package scanner

import (
	"fmt"
	"os"
	"path"
	"pkg.deepin.io/lib/mime"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	themeTypeGtk    = "gtk"
	themeTypeIcon   = "icon"
	themeTypeCursor = "cursor"
)

// uri: ex "file:///usr/share/themes"
func ListGtkTheme(uri string) ([]string, error) {
	return doListTheme(uri, themeTypeGtk, IsGtkTheme)
}

// uri: ex "file:///usr/share/icons"
func ListIconTheme(uri string) ([]string, error) {
	return doListTheme(uri, themeTypeIcon, IsIconTheme)
}

// uri: ex "file:///usr/share/icons"
func ListCursorTheme(uri string) ([]string, error) {
	return doListTheme(uri, themeTypeCursor, IsCursorTheme)
}

func IsGtkTheme(uri string) bool {
	ty, err := mime.Query(uri)
	if ty == mime.MimeTypeGtk {
		return true
	}
	fmt.Printf("%s not gtk theme: %v\n", uri, err)
	return false
}

func IsIconTheme(uri string) bool {
	ty, err := mime.Query(uri)
	if ty == mime.MimeTypeIcon {
		return true
	}
	fmt.Printf("%s not icon theme: %v\n", uri, err)
	return false
}

func IsCursorTheme(uri string) bool {
	ty, err := mime.Query(uri)
	if ty == mime.MimeTypeCursor {
		return true
	}
	fmt.Printf("%s not cursor theme: %v\n", uri, err)
	return false
}

func doListTheme(uri string, ty string, filter func(string) bool) ([]string, error) {
	dir := dutils.DecodeURI(uri)
	subDirs, err := listSubDir(dir)
	if err != nil {
		return nil, err
	}

	var themes []string
	for _, subDir := range subDirs {
		var tmp string
		if ty == themeTypeCursor {
			tmp = path.Join(subDir, "cursor.theme")
		} else {
			tmp = path.Join(subDir, "index.theme")
		}
		if !filter(tmp) {
			continue
		}
		themes = append(themes, subDir)
	}
	return themes, nil
}

func listSubDir(dir string) ([]string, error) {
	if !dutils.IsDir(dir) {
		return nil, fmt.Errorf("'%s' not a dir", dir)
	}

	fr, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer fr.Close()

	names, err := fr.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	var subDirs []string
	for _, name := range names {
		tmp := path.Join(dir, name)
		if !dutils.IsDir(tmp) {
			continue
		}

		subDirs = append(subDirs, tmp)
	}
	return subDirs, nil

}
