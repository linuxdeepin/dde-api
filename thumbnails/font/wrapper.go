package font

// #cgo pkg-config: cairo-ft glib-2.0
// #cgo LDFLAGS: -lm
// #include <stdlib.h>
// #include "thumbnail.h"
import "C"

import (
	"fmt"
	"pkg.deepin.io/dde/api/thumbnails/images"
	. "pkg.deepin.io/dde/api/thumbnails/loader"
	"unsafe"
)

func doGenThumbnail(file, dest string, width, height int) (string, error) {
	cFile := C.CString(file)
	defer C.free(unsafe.Pointer(cFile))
	tmp := GetTmpImage()
	cTmp := C.CString(tmp)
	defer C.free(unsafe.Pointer(cTmp))
	ret := C.gen_thumbnail(cFile, cTmp, C.int(getThumbSize(width, height)))
	if ret == -1 {
		return "", fmt.Errorf("Gen thumbnail for '%s' failed", file)
	}

	err := images.Scale(tmp, dest, width, height)
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
