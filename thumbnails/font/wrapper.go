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
