/**
 * Copyright (c) 2011 ~ 2015 Deepin, Inc.
 *               2013 ~ 2015 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package gtk

// #cgo pkg-config: libmetacity-private
// #include <stdlib.h>
// #include "common.h"
import "C"

import (
	"fmt"
	"os"
	"path"
	"unsafe"

	"pkg.deepin.io/dde/api/thumbnails/loader"
	dutils "pkg.deepin.io/lib/utils"
)

func doGenThumbnail(name, dest, bg string, w, h int, force bool) (string, error) {
	if !force && dutils.IsFileExist(dest) {
		return dest, nil
	}

	bg = dutils.DecodeURI(bg)
	if len(bg) == 0 {
		tmp, err := loader.GetBackground(w, h)
		if err != nil {
			return "", err
		}
		bg = tmp
		defer os.Remove(bg)
	}

	if C.try_init() != 0 {
		return "", fmt.Errorf("Init gtk environment failed")
	}

	os.MkdirAll(path.Dir(dest), 0755)
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cDest := C.CString(dest)
	defer C.free(unsafe.Pointer(cDest))
	cBg := C.CString(bg)
	defer C.free(unsafe.Pointer(cBg))

	ret := C.gtk_thumbnail(cName, cDest, cBg, C.int(w), C.int(h))
	if ret == -1 {
		return "", fmt.Errorf("MetaTheme load failed for '%s'", name)
	}
	return dest, nil
}
