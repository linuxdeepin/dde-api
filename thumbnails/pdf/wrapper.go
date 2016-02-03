/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package pdf

// #cgo pkg-config: poppler-glib cairo
// #include <stdlib.h>
// #include "thumbnail.h"
import "C"

import (
	"fmt"
	"os"
	. "pkg.deepin.io/dde/api/thumbnails/loader"
	"unsafe"
)

func doGenThumbnail(uri, dest string, width, height int) (string, error) {
	tmp := GetTmpImage()
	cTmp := C.CString(tmp)
	defer C.free(unsafe.Pointer(cTmp))
	cUri := C.CString(uri)
	defer C.free(unsafe.Pointer(cUri))
	ret := C.pdf_thumbnail(cUri, cTmp)
	if ret == -1 {
		return "", fmt.Errorf("Gen thumbnail failed")
	}

	defer os.Remove(tmp)
	err := ThumbnailImage(tmp, dest, width, height)
	if err != nil {
		return "", err
	}

	return dest, nil
}
