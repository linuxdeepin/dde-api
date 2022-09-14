// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package cursor

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/linuxdeepin/dde-api/thumbnails/loader"
	dutils "github.com/linuxdeepin/go-lib/utils"
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
