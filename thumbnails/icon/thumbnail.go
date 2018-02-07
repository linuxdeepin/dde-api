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

import (
	"os"
	"path/filepath"
	"pkg.deepin.io/dde/api/thumbnails/images"
	"pkg.deepin.io/dde/api/thumbnails/loader"
	"pkg.deepin.io/lib/graphic"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	defaultWidth    = 320
	defaultHeight   = 70
	defaultIconSize = 48
	defaultPadding  = 4
)

func doGenThumbnail(src, bg, dest string, width, height int, force, theme bool) (string, error) {
	if !force && dutils.IsFileExist(dest) {
		return dest, nil
	}

	src = dutils.DecodeURI(src)
	bg = dutils.DecodeURI(bg)
	dir := filepath.Dir(src)
	tmp := loader.GetTmpImage()
	themeName := filepath.Base(dir)
	iconFiles := getIconFiles(themeName)
	err := loader.CompositeIcons(iconFiles, bg, tmp,
		defaultIconSize, defaultWidth, defaultHeight, defaultPadding)
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp)

	if !theme {
		err = loader.ThumbnailImage(tmp, dest, width, height)
	} else {
		err = loader.ScaleImage(tmp, dest, width, height)
	}
	if err != nil {
		return "", err
	}

	return dest, nil
}

// default present icons
var presentIcons = [][]string{
	// file manager:
	{"dde-file-manager", "system-file-manager"},
	// music player:
	{"deepin-music", "banshee", "amarok", "deadbeef", "clementine", "rhythmbox"},
	// image viewer:
	{"deepin-image-viewer", "eog", "gthumb", "gwenview", "gpicview", "showfoto", "phototonic"},
	// media/video player:
	{"deepin-movie", "media-player", "totem", "smplayer", "vlc", "dragonplayer", "kmplayer"},
	// web browser:
	{"google-chrome", "firefox", "chromium", "opear", "internet-web-browser", "web-browser", "browser"},
	// system settings:
	{"preferences-system"},
	// text editor:
	{"accessories-text-editor", "text-editor", "gedit", "kedit", "xfce-edit"},
	// terminal:
	{"deepin-terminal", "utilities-terminal", "terminal", "gnome-terminal", "xfce-terminal", "terminator", "openterm"},
}

func getIconFiles(theme string) []string {
	var files []string
	for _, iconNames := range presentIcons {
		file := ChooseIcon(theme, iconNames)
		if file != "" {
			files = append(files, file)
			if len(files) == 6 {
				break
			}
		}
	}

	return fixIconFiles(files)
}

func fixIconFiles(files []string) []string {
	var ret []string
	for _, file := range files {
		ext := filepath.Ext(file)
		genThumbnail := false
		if ext == ".svg" {
			genThumbnail = true
		} else {
			// check size
			w, h, err := graphic.GetImageSize(file)
			if err != nil {
				continue
			}
			if !(w == defaultIconSize && w == h) {
				genThumbnail = true
			}
		}

		if genThumbnail {
			var err error
			file, err = images.GenThumbnail(file, defaultIconSize, defaultIconSize, true)
			if err != nil {
				continue
			}
		}
		ret = append(ret, file)
	}

	return ret
}
