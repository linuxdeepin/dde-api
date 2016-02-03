/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

//Cursor theme thumbnail generator
package cursor

import (
	"fmt"
	"os"
	"path"

	. "pkg.deepin.io/dde/api/thumbnails/loader"
	"pkg.deepin.io/lib/mime"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	sysThemeThumbDir = "/var/cache/thumbnails/appearance"
)

var themeThumbDir = path.Join(os.Getenv("HOME"),
	".cache", "thumbnails", "appearance")

func init() {
	for _, ty := range SupportedTypes() {
		Register(ty, GenThumbnail)
	}
}

func SupportedTypes() []string {
	return []string{
		mime.MimeTypeCursor,
	}
}

// GenThumbnail generate cursor theme thumbnail
// src: the uri of cursor theme index.theme
func GenThumbnail(src, bg string, width, height int, force bool) (string, error) {
	if width <= 0 || height <= 0 {
		return "", fmt.Errorf("Invalid width or height")
	}

	ty, err := mime.Query(src)
	if err != nil {
		return "", err
	}

	if ty != mime.MimeTypeCursor {
		return "", fmt.Errorf("Not supported type: %v", ty)
	}

	return genCursorThumbnail(src, bg, width, height, force)
}

// ThumbnailForTheme generate thumbnail for deepin appearance preview
func ThumbnailForTheme(src, bg string, width, height int, force bool) (string, error) {
	if width <= 0 || height <= 0 {
		return "", fmt.Errorf("Invalid width or height")
	}

	dest, err := getThumbDest(src, width, height, true)
	if err != nil {
		return "", err
	}

	thumb := path.Join(sysThemeThumbDir, path.Base(dest))
	if !force && dutils.IsFileExist(thumb) {
		return thumb, nil
	}

	return doGenThumbnail(src, bg, dest, width, height, force, true)
}

func genCursorThumbnail(src, bg string, width, height int, force bool) (string, error) {
	dest, err := getThumbDest(src, width, height, false)
	if err != nil {
		return "", err
	}

	return doGenThumbnail(src, bg, dest, width, height, force, false)
}

func getThumbDest(src string, width, height int, theme bool) (string, error) {
	var (
		dest string
		err  error
	)
	if dutils.IsFileExist(src) {
		dest, err = GetThumbnailDest(src, width, height)
	} else {
		dest, err = GetThumbnailDest(path.Join(path.Dir(dutils.DecodeURI(src)),
			"cursors", "left_ptr"), width, height)
		dest = path.Join(path.Dir(dest), "cursor-"+path.Base(dest))
	}
	if err != nil {
		return "", err
	}

	if theme {
		dest = path.Join(themeThumbDir, "cursor-"+path.Base(dest))
	}
	return dest, nil
}
