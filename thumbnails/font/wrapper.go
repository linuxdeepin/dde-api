// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package font

// #cgo pkg-config: cairo-ft glib-2.0
// #cgo CFLAGS: -W -Wall -fPIC -fstack-protector-all
// #cgo LDFLAGS: -lm
// #include <stdlib.h>
// #include "thumbnail.h"
import "C"

import (
	"fmt"
	"os"
	"unsafe"

	. "github.com/linuxdeepin/dde-api/thumbnails/loader"
)

func doGenThumbnail(file, dest string, width, height int) (string, error) {
	cFile := C.CString(file)
	defer C.free(unsafe.Pointer(cFile))
	tmp := GetTmpImage()
	cTmp := C.CString(tmp)
	defer C.free(unsafe.Pointer(cTmp))
	ret := C.font_thumbnail(cFile, cTmp, C.int(getThumbSize(width, height)))
	if ret == -1 {
		return "", fmt.Errorf("Gen thumbnail for '%s' failed", file)
	}

	defer os.Remove(tmp)
	err := ThumbnailImage(tmp, dest, width, height)
	if err != nil {
		return "", err
	}

	return dest, nil
}

func getThumbSize(width, height int) int {
	if width > height {
		return width
	}
	return height
}
