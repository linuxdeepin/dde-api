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
// #include "icon.h"
import "C"

import (
	"unsafe"
)

func ChooseIcon(theme string, names []string) string {
	cTheme := C.CString(theme)
	defer C.free(unsafe.Pointer(cTheme))

	cArr := StrvInC(names)
	cNames := (**C.char)(unsafe.Pointer(&cArr[0]))
	cFile := C.choose_icon(cTheme, cNames, C.int(defaultIconSize))

	// free cArr
	for i := range cArr {
		C.free(unsafe.Pointer(cArr[i]))
	}
	defer C.free(unsafe.Pointer(cFile))

	return C.GoString(cFile)
}

// return NUL-Terminated slice of C String
func StrvInC(strv []string) []*C.char {
	cArr := make([]*C.char, len(strv)+1)
	for i, str := range strv {
		cArr[i] = C.CString(str)
	}
	return cArr
}
