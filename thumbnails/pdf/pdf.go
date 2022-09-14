// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

// PDF thumbnail generator
package pdf

import (
	"fmt"

	. "github.com/linuxdeepin/dde-api/thumbnails/loader"
	"github.com/linuxdeepin/go-lib/mime"
	dutils "github.com/linuxdeepin/go-lib/utils"
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
