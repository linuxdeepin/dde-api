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
