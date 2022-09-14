// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package thumbnails

import (
	"fmt"

	_ "github.com/linuxdeepin/dde-api/thumbnails/cursor"
	_ "github.com/linuxdeepin/dde-api/thumbnails/font"
	_ "github.com/linuxdeepin/dde-api/thumbnails/gtk"
	_ "github.com/linuxdeepin/dde-api/thumbnails/icon"
	_ "github.com/linuxdeepin/dde-api/thumbnails/images"
	"github.com/linuxdeepin/dde-api/thumbnails/loader"
	_ "github.com/linuxdeepin/dde-api/thumbnails/pdf"
	_ "github.com/linuxdeepin/dde-api/thumbnails/text"
	"github.com/linuxdeepin/go-lib/mime"
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
