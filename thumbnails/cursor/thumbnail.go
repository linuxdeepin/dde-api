/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package cursor

import (
	"os"
	"path"
	"pkg.deepin.io/dde/api/thumbnails/loader"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	defaultWidth    = 320
	defaultHeight   = 70
	defaultIconSize = 24
	defaultPadding  = 12
)

func doGenThumbnail(src, bg, dest string, width, height int, force, theme bool) (string, error) {
	if !force && dutils.IsFileExist(dest) {
		return dest, nil
	}

	src = dutils.DecodeURI(src)
	bg = dutils.DecodeURI(bg)
	dir := path.Dir(src)
	tmp := loader.GetTmpImage()
	cursorIcons := getCursorIcons(dir)
	err := loader.CompositeIcons(cursorIcons, bg, tmp,
		defaultIconSize, defaultWidth, defaultHeight, defaultPadding)
	os.RemoveAll(xcur2pngCache)
	if err != nil {
		return "", err
	}

	defer os.Remove(tmp)
	if !theme {
		err = loader.ThumbnailImage(tmp, dest, width, height)
	} else {
		err = loader.ScaleImage(tmp, dest, width, height)
	}
	if err != nil {
		return "", err
	}

	return dest, nil
}

var presentCursors = [][]string{
	{"left_ptr"},
	{"left_ptr_watch"},
	{"x-cursor", "X_cursor"},
	{"hand2", "hand1"},
	{"grab", "grabbing", "closedhand"},
	{"move"},
	{"sb_v_double_arrow"},
	{"sb_h_double_arrow"},
	{"watch", "wait"},
}

func getCursorIcons(dir string) []string {
	var files []string
	for _, cursors := range presentCursors {
		for _, cursor := range cursors {
			tmp, err := XCursorToPng(path.Join(dir, "cursors", cursor))
			if err == nil {
				files = append(files, tmp)
				break
			}
		}
	}
	return files
}
