// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

// Theme scanner
package scanner

import (
	"fmt"
	"os"
	"path"
	"github.com/linuxdeepin/go-lib/mime"
	dutils "github.com/linuxdeepin/go-lib/utils"
)

const (
	ThemeTypeGtk    = "gtk"
	ThemeTypeIcon   = "icon"
	ThemeTypeCursor = "cursor"
)

// uri: ex "file:///usr/share/themes"
func ListGtkTheme(uri string) ([]string, error) {
	return doListTheme(uri, ThemeTypeGtk, IsGtkTheme)
}

// uri: ex "file:///usr/share/icons"
func ListIconTheme(uri string) ([]string, error) {
	return doListTheme(uri, ThemeTypeIcon, IsIconTheme)
}

// uri: ex "file:///usr/share/icons"
func ListCursorTheme(uri string) ([]string, error) {
	return doListTheme(uri, ThemeTypeCursor, IsCursorTheme)
}

func IsGtkTheme(uri string) bool {
	if len(uri) == 0 {
		return false
	}

	ty, _ := mime.Query(uri)
	return ty == mime.MimeTypeGtk
}

func IsIconTheme(uri string) bool {
	if len(uri) == 0 {
		return false
	}

	ty, _ := mime.Query(uri)
	return ty == mime.MimeTypeIcon
}

func IsCursorTheme(uri string) bool {
	if len(uri) == 0 {
		return false
	}

	ty, _ := mime.Query(uri)
	return ty == mime.MimeTypeCursor
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
		if ty == ThemeTypeCursor {
			tmp = path.Join(subDir, "cursor.theme")
		} else {
			tmp = path.Join(subDir, "index.theme")
		}
		if !filter(tmp) || isHidden(tmp, ty) {
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
	defer func() {
		fr.Close()
	}()

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

func isHidden(file, ty string) bool {
	kf, err := dutils.NewKeyFileFromFile(file)
	if err != nil {
		return false
	}
	defer kf.Free()

	var hidden bool = false
	switch ty {
	case ThemeTypeGtk:
		hidden, _ = kf.GetBoolean("Desktop Entry", "Hidden")
	case ThemeTypeIcon, ThemeTypeCursor:
		hidden, _ = kf.GetBoolean("Icon Theme", "Hidden")
	}
	return hidden
}
