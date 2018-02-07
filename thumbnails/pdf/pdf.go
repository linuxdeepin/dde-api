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

// PDF thumbnail generator
package pdf

import (
	"fmt"
	. "pkg.deepin.io/dde/api/thumbnails/loader"
	"pkg.deepin.io/lib/mime"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	PDFTypePDF = "application/pdf"
)

func init() {
	for _, ty := range SupportedTypes() {
		Register(ty, genPDFThumbnail)
	}
}

func SupportedTypes() []string {
	return []string{
		PDFTypePDF,
	}
}

func GenThumbnail(uri string, width, height int, force bool) (string, error) {
	if width <= 0 || height <= 0 {
		return "", fmt.Errorf("Invalid width or height")
	}

	ty, err := mime.Query(uri)
	if err != nil {
		return "", err
	}

	if !IsStrInList(ty, SupportedTypes()) {
		return "", fmt.Errorf("Not supported type: %s", ty)
	}

	return genPDFThumbnail(uri, "", width, height, force)
}

func genPDFThumbnail(uri, bg string, width, height int, force bool) (string, error) {
	dest, err := GetThumbnailDest(uri, width, height)
	if err != nil {
		return "", err
	}

	if !force && dutils.IsFileExist(dest) {
		return dest, nil
	}

	return doGenThumbnail(dutils.EncodeURI(uri, dutils.SCHEME_FILE),
		dest, width, height)
}
