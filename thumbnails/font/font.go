/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package font

import (
	"fmt"
	. "pkg.deepin.io/dde/api/thumbnails/loader"
	"pkg.deepin.io/lib/mime"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	FontTypeTTF = "application/x-font-ttf"
	FontTypeOTF = "application/vnd.ms-opentype"
)

func init() {
	for _, ty := range SupportedTypes() {
		Register(ty, genFontThumbnail)
	}
}

func SupportedTypes() []string {
	return []string{
		FontTypeOTF,
		FontTypeTTF,
	}
}

func GenThumbnail(src string, width, height int, force bool) (string, error) {
	if width <= 0 || height <= 0 {
		return "", fmt.Errorf("Invalid width or height")
	}

	ty, err := mime.Query(src)
	if err != nil {
		return "", err
	}

	if !IsStrInList(ty, SupportedTypes()) {
		return "", fmt.Errorf("Not supported type: %v", ty)
	}

	return genFontThumbnail(src, "", width, height, force)
}

func genFontThumbnail(src, bg string, width, height int, force bool) (string, error) {
	dest, err := GetThumbnailDest(src, width, height)
	if err != nil {
		return "", err
	}
	if !force && dutils.IsFileExist(dest) {
		return dest, nil
	}

	return doGenThumbnail(dutils.DecodeURI(src), dest, width, height)
}
