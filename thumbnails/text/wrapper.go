// TODO: use 'code.google.com/p/freetype-go/freetype' to draw text on image
package text

// #cgo pkg-config: gdk-3.0
// #include <stdlib.h>
// #include "text.h"
import "C"

import (
	"fmt"
	"io/ioutil"
	"path"
	"pkg.deepin.io/dde/api/thumbnails/images"
	"strings"
	"unsafe"
)

func doGenThumbnail(src, dest string, width, height int) (string, error) {
	strv, err := readFileContent(src)
	if err != nil {
		return "", err
	}
	defer freeCStrv(strv)

	tmp := path.Join("/tmp", path.Base(src)+".png")
	cTmp := C.CString(tmp)
	defer C.free(unsafe.Pointer(cTmp))
	ret := C.do_gen_thumbnail(&strv[0], cTmp)
	if ret != 0 {
		return "", fmt.Errorf("Draw text on image failed")
	}

	err = images.Scale(tmp, dest, width, height)
	if err != nil {
		return "", err
	}
	return dest, nil
}

func readFileContent(file string) ([]*C.char, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	length := len(lines)
	var strv []*C.char
	for i, line := range lines {
		if i == length-1 && len(line) == 0 {
			continue
		}

		var tmp string
		for _, ch := range line {
			// tab: 4 space
			if ch == '\t' {
				tmp += "    "
				continue
			}
			tmp += string(ch)
		}
		strv = append(strv, C.CString(tmp))
	}
	strv = append(strv, nil)

	return strv, nil
}

func freeCStrv(strv []*C.char) {
	for _, s := range strv {
		C.free(unsafe.Pointer(s))
	}
}
