/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package icon

// #cgo pkg-config: gtk+-3.0
// #include <stdlib.h>
// #include "lookup.h"
import "C"

import (
	"unsafe"
)

func GetIconFile(theme, name string) string {
	cTheme := C.CString(theme)
	defer C.free(unsafe.Pointer(cTheme))
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cFile := C.lookup_icon(cTheme, cName, C.int(defaultIconSize))
	defer C.free(unsafe.Pointer(cFile))

	return C.GoString(cFile)
}
