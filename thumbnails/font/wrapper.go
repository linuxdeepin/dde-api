/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package font

// #cgo pkg-config: cairo-ft glib-2.0
// #cgo LDFLAGS: -lm
// #include <stdlib.h>
// #include "thumbnail.h"
import "C"

import (
	"fmt"
	"os"
	. "pkg.deepin.io/dde/api/thumbnails/loader"
	"unsafe"
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
