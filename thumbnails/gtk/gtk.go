/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

// Gtk theme thumbnail generator
package gtk

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	. "pkg.deepin.io/dde/api/thumbnails/loader"
	"pkg.deepin.io/lib/mime"
	dutils "pkg.deepin.io/lib/utils"
	"strconv"
)

const (
	sysThemeThumbDir = "/var/cache/thumbnails/appearance"

	cmdGtkThumbnailer = "/usr/lib/deepin-api/gtk-thumbnailer"
)

var themeThumbDir = path.Join(os.Getenv("HOME"),
	".cache", "thumbnails", "appearance")

func init() {
	for _, ty := range SupportedTypes() {
		Register(ty, genGtkThumbnail)
	}
}

func SupportedTypes() []string {
	return []string{
		mime.MimeTypeGtk,
	}
}

func GenThumbnail(src, bg string, width, height int, force bool) (string, error) {
	if width <= 0 || height <= 0 {
		return "", fmt.Errorf("Invalid width or height")
	}

	ty, err := mime.Query(src)
	if err != nil {
		return "", err
	}
	if ty != mime.MimeTypeGtk {
		return "", fmt.Errorf("Unspported mime: %s", ty)
	}

	return genGtkThumbnail(src, bg, width, height, force)
}

func ThumbnailForTheme(src, bg string, width, height int, force bool) (string, error) {
	if width <= 0 || height <= 0 {
		return "", fmt.Errorf("Invalid width or height")
	}

	dest, err := getThumbDest(src, width, height, true)
	if err != nil {
		return "", err
	}

	thumb := path.Join(sysThemeThumbDir, path.Base(dest))
	if !force && dutils.IsFileExist(thumb) {
		return thumb, nil
	}

	return doGenThumbnail(path.Base(path.Dir(src)), bg, dest,
		width, height, force)
}

func genGtkThumbnail(src, bg string, width, height int, force bool) (string, error) {
	dest, err := getThumbDest(src, width, height, false)
	if err != nil {
		return "", err
	}

	return doGenThumbnail(path.Base(path.Dir(src)), bg, dest,
		width, height, force)
}

func getThumbDest(src string, width, height int, theme bool) (string, error) {
	dest, err := GetThumbnailDest(src, width, height)
	if err != nil {
		return "", err
	}

	if theme {
		dest = path.Join(themeThumbDir, "gtk-"+path.Base(dest))
	}
	return dest, nil
}

func doGenThumbnail(name, bg, dest string, width, height int, force bool) (string, error) {
	var args = []string{
		"-theme", name,
		"-dest", dest,
		"-width", strconv.Itoa(width),
		"-height", strconv.Itoa(height),
	}
	if force {
		args = append(args, "-force")
	}
	_, err := exec.Command(cmdGtkThumbnailer, args...).CombinedOutput()
	if err != nil {
		return "", err
	}
	return dest, nil
}
