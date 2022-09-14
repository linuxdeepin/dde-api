// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/linuxdeepin/go-dbus-factory/com.deepin.sessionmanager"
	"github.com/linuxdeepin/go-gir/gio-2.0"
	"github.com/godbus/dbus"
	"github.com/linuxdeepin/go-lib/log"
)

var logger = log.NewLogger("dde-open")

var optVersion bool

func init() {
	flag.BoolVar(&optVersion, "version", false, "show version")
}

func main() {
	flag.Parse()
	if optVersion {
		fmt.Println("1.0")
		os.Exit(0)
	}

	if len(flag.Args()) != 1 {
		fmt.Println("usage: dde-open { file | URL }")
		os.Exit(1)
	}
	arg := flag.Arg(0)

	u, err := url.Parse(arg)
	if err != nil {
		logger.Warningf("failed to parse url %q: %v", arg, err)
		err = openFile(arg)

	} else {
		switch u.Scheme {
		case "file":
			err = openFile(u.Path)

		case "":
			err = openFile(arg)

		default:
			err = openScheme(u.Scheme, arg)
		}
	}
	if err != nil {
		logger.Warning("open failed:", err)
		os.Exit(2)
	}
}

func launchApp(desktopFile, filename string) error {
	sessionBus, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	startManager := sessionmanager.NewStartManager(sessionBus)
	err = startManager.LaunchApp(dbus.FlagNoAutoStart, desktopFile, 0,
		[]string{filename})
	return err
}

func openFile(filename string) error {
	logger.Debugf("openFile: %q", filename)
	filename, err := filepath.Abs(filename)
	if err != nil {
		return err
	}
	_, err = os.Stat(filename)
	if err != nil {
		return err
	}
	file := gio.FileNewForPath(filename)
	defer file.Unref()

	fileInfo, err := file.QueryInfo(gio.FileAttributeStandardContentType, gio.FileQueryInfoFlagsNone, nil)
	if err != nil {
		return err
	}
	defer fileInfo.Unref()
	contentType := fileInfo.GetAttributeString(gio.FileAttributeStandardContentType)
	if contentType == "" {
		return errors.New("failed to get file content type")
	}

	appInfo := gio.AppInfoGetDefaultForType(contentType, false)
	if appInfo == nil {
		return errors.New("failed to get appInfo")
	}
	defer appInfo.Unref()

	dAppInfo := gio.ToDesktopAppInfo(appInfo)
	desktopFile := dAppInfo.GetFilename()
	logger.Debug("desktop file:", desktopFile)
	err = launchApp(desktopFile, filename)
	if err != nil {
		return err
	}
	return nil
}

func openScheme(scheme, url string) error {
	logger.Debugf("openScheme: %q, %q", scheme, url)
	appInfo := gio.AppInfoGetDefaultForUriScheme(scheme)
	if appInfo == nil {
		return errors.New("failed to get appInfo")
	}
	defer appInfo.Unref()

	dAppInfo := gio.ToDesktopAppInfo(appInfo)
	desktopFile := dAppInfo.GetFilename()
	logger.Debug("desktop file:", desktopFile)
	err := launchApp(desktopFile, url)
	if err != nil {
		return err
	}
	return nil
}
