/**
 * Copyright (c) 2013 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

import (
	"dlib/dbus"
	libgraphic "dlib/graphic"
	liblogger "dlib/logger"
	"os"
)

var logger = liblogger.NewLogger("dde-api/graphic")

// Graphic is a dbus interface wrapper for dlib/graphic.
type Graphic struct {
	BlurPictChanged func(string, string)
}

// GetDBusInfo implement interface of dbus.DBusObject
func (graphic *Graphic) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.api.Graphic",
		"/com/deepin/api/Graphic",
		"com.deepin.api.Graphic",
	}
}

// RGB2HSV convert color format from RGB(r, g, b=[0..255]) to HSV(h=[0..360), s,v=[0..1]).
func (graphic *Graphic) RGB2HSV(r, g, b uint8) (h, s, v float64) {
	return libgraphic.RGB2HSV(r, g, b)
}

// HSV2RGB convert color format from HSV(h=[0..360), s,v=[0..1]) to RGB(r, g, b=[0..255]).
func (graphic *Graphic) HSV2RGB(h, s, v float64) (r, g, b uint8) {
	return libgraphic.HSV2RGB(h, s, v)
}

// GetImageSize return a image's width and height.
func (graphic *Graphic) GetImageSize(imgfile string) (w, h int32, err error) {
	w, h, err = libgraphic.GetImageSize(imgfile)
	if err != nil {
		logger.Error("%v", err)
	}
	return
}

// GetDominantColorOfImage return the dominant hsv color of a image.
func (graphic *Graphic) GetDominantColorOfImage(imgfile string) (h, s, v float64, err error) {
	h, s, v, err = libgraphic.GetDominantColorOfImage(imgfile)
	if err != nil {
		logger.Error("%v", err)
	}
	return
}

// ConvertImage converts from any recognized format imaget to target
// format image which could be "png" or "jpeg".
func (graphic *Graphic) ConvertImage(srcfile, dstfile, f string) (err error) {
	err = libgraphic.ConvertImage(srcfile, dstfile, libgraphic.Format(f))
	if err != nil {
		logger.Error("%v", err)
	}
	return
}

// ClipImage clip any recognized format image to target format image
// which could be "png" or "jpeg".
func (graphic *Graphic) ClipImage(srcfile, dstfile string, x0, y0, x1, y1 int32, f string) (err error) {
	err = libgraphic.ClipImage(srcfile, dstfile, x0, y0, x1, y1, libgraphic.Format(f))
	if err != nil {
		logger.Error("%v", err)
	}
	return
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			logger.Fatal("%v", err)
		}
	}()

	// configure logger
	logger.SetRestartCommand("/usr/lib/deepin-api/graphic", "--debug")
	if stringInSlice("-d", os.Args) || stringInSlice("--debug", os.Args) {
		logger.SetLogLevel(liblogger.LEVEL_DEBUG)
	}

	jobInHand = make(map[string]bool) // used by blur pict

	graphic := &Graphic{}
	err := dbus.InstallOnSession(graphic)
	if err != nil {
		logger.Error("register dbus interface failed: %v", err)
		os.Exit(1)
	}
	dbus.DealWithUnhandledMessage()

	if err := dbus.Wait(); err != nil {
		logger.Error("lost dbus session: %v", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
