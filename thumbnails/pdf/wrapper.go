package pdf

// #cgo pkg-config: poppler-glib cairo
// #include <stdlib.h>
// #include "thumbnail.h"
import "C"

import (
	"fmt"
	"pkg.deepin.io/dde/api/thumbnails/images"
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

	err := images.Scale(tmp, dest, width, height)
	if err != nil {
		return "", err
	}

	return dest, nil
}
