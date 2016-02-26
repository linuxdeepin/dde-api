/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

import (
	"os"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

const (
	dbusDest = "com.deepin.api.GtkThumbnailer"
	dbusPath = "/com/deepin/api/GtkThumbnailer"
	dbusIFC  = "com.deepin.api.GtkThumbnailer"
)

var (
	logger = log.NewLogger("api/GtkThumbnailer")

	_force  = kingpin.Flag("force", "Force to generate thumbnail").Short('f').Bool()
	_src    = kingpin.Arg("src", "The source").String()
	_bg     = kingpin.Arg("bg", "The background").String()
	_dest   = kingpin.Arg("dest", "The dest").String()
	_width  = kingpin.Arg("width", "The thumbnail width").Int()
	_height = kingpin.Arg("height", "The thumbnail height").Int()
)

type Manager struct {
	running bool
}

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func (m *Manager) Thumbnail(name, bg, dest string, width, height int32, force bool) error {
	m.running = true
	defer func() {
		m.running = false
	}()
	return doGenThumbnail(name, bg, dest, int(width), int(height), force)
}

func main() {
	err := initGtkEnv()
	if err != nil {
		logger.Error(err)
		return
	}

	kingpin.Parse()
	if len(os.Args) > 5 {
		err := doGenThumbnail(*_src, *_bg, *_dest, *_width, *_height, *_force)
		if err != nil {
			logger.Error("Generate gtk thumbnail failed:", *_src, *_dest, err)
		}
		return
	}

	if !lib.UniqueOnSession(dbusDest) {
		logger.Warning("There already has a gtk thumbnailer running...")
		return
	}

	var m = new(Manager)
	m.running = false
	err = dbus.InstallOnSession(m)
	if err != nil {
		logger.Error("Install dbus session failed:", err)
		return
	}
	dbus.DealWithUnhandledMessage()

	dbus.SetAutoDestroyHandler(time.Second*5, func() bool {
		if m.running {
			return false
		}
		return true
	})

	err = dbus.Wait()
	if err != nil {
		logger.Error("Lost dbus connect:", err)
		os.Exit(-1)
	}
	os.Exit(0)
}
