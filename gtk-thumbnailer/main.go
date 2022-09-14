// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

// #cgo pkg-config: gtk+-3.0
// #cgo CFLAGS: -W -Wall -fPIC -fstack-protector-all
// #include <stdlib.h>
// void gtk_thumbnail(char *theme, char *dest, int width, int min_height);
import "C"

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"unsafe"
)

var (
	force  = flag.Bool("force", false, "Force to generate thumbnail")
	theme  = flag.String("theme", "", "The theme name")
	dest   = flag.String("dest", "", "The destination of thumbnail file")
	width  = flag.Int("width", 0, "The thumbnail width")
	height = flag.Int("height", 0, "The thumbnail min height")
)

func main() {
	flag.Parse()
	if flag.Parsed() {

		if *theme == "" || *dest == "" || *width == 0 || *height == 0 {
			flag.Usage()
			os.Exit(1)
		}

		err := doGenThumbnail(*theme, *dest, *width, *height, *force)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	}
}

func doGenThumbnail(name, dest string, width, height int, force bool) error {
	if _, err := os.Stat(dest); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		// file dest not exist
	} else {
		// file dest exist
		if force {
			os.Remove(dest)
		} else {
			return nil
		}
	}

	err := os.MkdirAll(filepath.Dir(dest), 0755)
	if err != nil {
		return err
	}
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cDest := C.CString(dest)
	defer C.free(unsafe.Pointer(cDest))
	C.gtk_thumbnail(cName, cDest, C.int(width), C.int(height))

	// check thumbnail result
	_, err = os.Stat(dest)
	return err
}
