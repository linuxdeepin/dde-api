// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

// Text thumbnail generator
package text

import (
	"fmt"

	. "github.com/linuxdeepin/dde-api/thumbnails/loader"
	"github.com/linuxdeepin/go-lib/mime"
	dutils "github.com/linuxdeepin/go-lib/utils"
)

const (
	TextTypeText = "text/plain"
	TextTypeC    = "text/x-c"
	TextTypeCpp  = "text/x-cpp"
	TextTypeGo   = "text/x-go"
)

func init() {
	for _, ty := range SupportedTypes() {
		Register(ty, genTextThumbnail)
	}
}

func SupportedTypes() []string {
	return []string{
		TextTypeText,
		TextTypeC,
		TextTypeCpp,
		TextTypeGo,
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

	return genTextThumbnail(src, "", width, height, force)
}

func genTextThumbnail(src, bg string, width, height int, force bool) (string, error) {
	dest, err := GetThumbnailDest(src, width, height)
	if err != nil {
		return "", err
	}
	if !force && dutils.IsFileExist(dest) {
		return dest, nil
	}

	return doGenThumbnail(dutils.DecodeURI(src), dest, width, height)
}
