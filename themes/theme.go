/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

// Gtk/Icon/Cursor theme scanner.
package themes

import (
	"os"
	"path"
	"pkg.deepin.io/dde/api/themes/scanner"
)

// Check whether 'theme' in 'list'
func IsThemeInList(theme string, list []string) bool {
	name := path.Base(theme)
	for _, l := range list {
		if path.Base(l) == name {
			return true
		}
	}
	return false
}

// List gtk theme in system.
//
// Scan '/usr/share/themes' and '$HOME/.themes'
func ListGtkTheme() []string {
	return doListTheme(
		[]string{
			path.Join(os.Getenv("HOME"), ".local/share/themes"),
			path.Join(os.Getenv("HOME"), ".themes"),
		},
		[]string{"/usr/share/themes"},
		scanner.ListGtkTheme)
}

// List icon theme in system.
//
// Scan '/usr/share/icons' and '$HOME/.icons'
func ListIconTheme() []string {
	return doListTheme(
		[]string{
			path.Join(os.Getenv("HOME"), ".local/share/icons"),
			path.Join(os.Getenv("HOME"), ".icons"),
		},
		[]string{"/usr/share/icons"},
		scanner.ListIconTheme)
}

// List cursor theme in system.
//
// Scan '/usr/share/icons' and '$HOME/.icons'
func ListCursorTheme() []string {
	return doListTheme(
		[]string{
			path.Join(os.Getenv("HOME"), ".local/share/icons"),
			path.Join(os.Getenv("HOME"), ".icons"),
		},
		[]string{"/usr/share/icons"},
		scanner.ListCursorTheme)
}

func doListTheme(local, sys []string, scanner func(string) ([]string, error)) []string {
	list := scanThemeDirs(local, scanner)
	sysList := scanThemeDirs(sys, scanner)
	return mergeThemeList(list, sysList)
}

func scanThemeDirs(dirs []string, scanner func(string) ([]string, error)) []string {
	var list []string
	for _, d := range dirs {
		tmp, err := scanner(d)
		if err != nil {
			continue
		}
		list = append(list, tmp...)
	}
	return list
}

func mergeThemeList(src, target []string) []string {
	if len(target) == 0 {
		return src
	}

	for _, t := range target {
		if IsThemeInList(t, src) {
			continue
		}
		src = append(src, t)
	}
	return src
}
