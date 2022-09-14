// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package themes

// #cgo pkg-config: x11 xcursor xfixes gtk+-3.0
// #cgo CFLAGS: -W -Wall -fPIC -fstack-protector-all
// #include <stdlib.h>
// #include "cursor.h"
import "C"
import (
	"unsafe"
)

func setGtkCursor(name string) {
	C.init_gtk()
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.set_gtk_cursor(cName)
}

func setQtCursor(name string) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.set_qt_cursor(cName)
}
