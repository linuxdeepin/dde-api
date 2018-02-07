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

package thumbnails

import (
	"fmt"
	_ "pkg.deepin.io/dde/api/thumbnails/cursor"
	_ "pkg.deepin.io/dde/api/thumbnails/font"
	_ "pkg.deepin.io/dde/api/thumbnails/gtk"
	_ "pkg.deepin.io/dde/api/thumbnails/icon"
	_ "pkg.deepin.io/dde/api/thumbnails/images"
	"pkg.deepin.io/dde/api/thumbnails/loader"
	_ "pkg.deepin.io/dde/api/thumbnails/pdf"
	_ "pkg.deepin.io/dde/api/thumbnails/text"
	"pkg.deepin.io/lib/mime"
)

func GenThumbnail(uri string, size int) (string, error) {
	if size < 0 {
		return "", fmt.Errorf("Invalid size: '%v'", size)
	}

	ty, err := mime.Query(uri)
	if err != nil {
		return "", err
	}

	size = correctSize(size)
	return GenThumbnailWithMime(uri, ty, size)
}

func GenThumbnailWithMime(uri, ty string, size int) (string, error) {
	if size < 0 {
		return "", fmt.Errorf("Invalid size: '%v'", size)
	}

	handler, err := loader.GetHandler(ty)
	if err != nil {
		return "", err
	}

	size = correctSize(size)
	return handler(uri, "", size, size, false)
}

func correctSize(size int) int {
	if size < loader.SizeFlagNormal {
		return loader.SizeFlagSmall
	} else if size >= loader.SizeFlagLarge {
		return loader.SizeFlagLarge
	} else {
		return loader.SizeFlagNormal
	}
}
