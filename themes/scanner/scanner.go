/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
