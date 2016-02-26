/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

// #cgo pkg-config: libdeepin-metacity-private
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

func initGtkEnv() error {
	if C.try_init() != 0 {
		return fmt.Errorf("Init gtk environment failed")
	}
	return nil
}

func doGenThumbnail(name, bg, dest string, w, h int, force bool) error {
	if !force && dutils.IsFileExist(dest) {
		return nil
	}

	bg = dutils.DecodeURI(bg)
	if len(bg) == 0 {
		tmp, err := loader.GetBackground(w, h)
		if err != nil {
			return err
		}
		bg = tmp
		defer os.Remove(bg)
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
		return fmt.Errorf("MetaTheme load failed for '%s'", name)
	}
	return nil
}
