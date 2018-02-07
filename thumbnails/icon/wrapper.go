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
