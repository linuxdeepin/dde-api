/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package cursor

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

const (
	xcur2pngTool  = "xcur2png"
	xcur2pngCache = "/tmp/xcur2png-cache"
)

func XCursorToPng(file string) (string, error) {
	os.MkdirAll(xcur2pngCache, 0755)
	out, err := exec.Command("/bin/sh", "-c",
		fmt.Sprintf("%s -c %s -d %s -q %s", xcur2pngTool,
			xcur2pngCache, xcur2pngCache, file)).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf(string(out))
	}

	// 000: 24x24
	// some images are only size of 24x24
	return path.Join(xcur2pngCache, path.Base(file)+"_000.png"), nil
}
