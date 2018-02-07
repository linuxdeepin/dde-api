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

package cursor

import (
	"io/ioutil"
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
	cacheDir, err := ioutil.TempDir("", "xcur2png")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(cacheDir)

	cursorIcons := getCursorIcons(dir, cacheDir)
	err = loader.CompositeIcons(cursorIcons, bg, tmp,
		defaultIconSize, defaultWidth, defaultHeight, defaultPadding)
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

func getCursorIcons(dir, cacheDir string) []string {
	var files []string
	for _, cursors := range presentCursors {
		for _, cursor := range cursors {
			tmp, err := XCursorToPng(path.Join(dir, "cursors", cursor), cacheDir)
			if err == nil {
				files = append(files, tmp)
				break
			}
		}
	}
	return files
}
