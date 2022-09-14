// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package images

// #cgo pkg-config: gdk-3.0 librsvg-2.0
// #cgo CFLAGS: -W -Wall -fPIC -fstack-protector-all
// #include <stdlib.h>
// #include "convert.h"
import "C"

import (
	"fmt"
	"unsafe"
)

func svgToPng(src, dest string) error {
	cSrc := C.CString(src)
	defer C.free(unsafe.Pointer(cSrc))
	cDest := C.CString(dest)
	defer C.free(unsafe.Pointer(cDest))

	ret := C.svg_to_png(cSrc, cDest)
	if ret != 0 {
		return fmt.Errorf("Convert svg to png failed")
	}
	return nil
}
