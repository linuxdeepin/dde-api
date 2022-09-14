// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

//Icon theme thumbnail generator
package icon

import (
	"fmt"
	"os"
	"path"

	. "github.com/linuxdeepin/dde-api/thumbnails/loader"
	"github.com/linuxdeepin/go-lib/mime"
	dutils "github.com/linuxdeepin/go-lib/utils"
)

const (
	sysThemeThumbDir = "/var/cache/thumbnails/appearance"
)

var themeThumbDir = path.Join(os.Getenv("HOME"),
	".cache", "thumbnails", "appearance")

func init() {
	for _, ty := range SupportedTypes() {
		Register(ty, genIconThumbnail)
	}
}

func SupportedTypes() []string {
	return []string{
		mime.MimeTypeIcon,
	}
}

// GenThumbnail generate icon theme thumbnail
// src: the uri of icon theme index.theme
func GenThumbnail(src, bg string, width, height int, force bool) (string, error) {
	if width <= 0 || height <= 0 {
		return "", fmt.Errorf("Invalid width or height")
	}

	ty, err := mime.Query(src)
	if err != nil {
		return "", err
	}

	if ty != mime.MimeTypeIcon {
		return "", fmt.Errorf("Not supported type: %v", ty)
	}

	return genIconThumbnail(src, bg, width, height, force)
}

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

func genIconThumbnail(src, bg string, width, height int, force bool) (string, error) {
	dest, err := getThumbDest(src, width, height, false)
	if err != nil {
		return "", err
	}

	return doGenThumbnail(src, bg, dest, width, height, force, false)
}

func getThumbDest(src string, width, height int, theme bool) (string, error) {
	dest, err := GetThumbnailDest(src, width, height)
	if err != nil {
		return "", err
	}

	if theme {
		dest = path.Join(themeThumbDir, "icon-"+path.Base(dest))
	}
	return dest, nil
}
