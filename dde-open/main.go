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
	"strings"

	"github.com/godbus/dbus/v5"
	appmanager "github.com/linuxdeepin/go-dbus-factory/session/org.deepin.dde.application1"
	gio "github.com/linuxdeepin/go-gir/gio-2.0"
	"github.com/linuxdeepin/go-lib/log"
)

const (
	dbusServiceName = "org.desktopspec.ApplicationManager1"
	dbusPath        = "/org/desktopspec/ApplicationManager1"
	dbusInterface   = dbusServiceName + ".Application"
)

type AppInfo struct {
	appId       string
	desktopFile string
}

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
	var scheme string
	u, err := url.Parse(arg)
	if err != nil || u == nil {
		gFile := gio.FileNewForCommandlineArg(arg)
		if gFile != nil {
			scheme = gFile.GetUriScheme()
		}
		if scheme == "" {
			logger.Warningf("failed to parse url %q: %v", arg, err)
		}
	} else {
		scheme = u.Scheme
	}
	logger.Debugf("scheme: %q", scheme)
	switch scheme {
	case "file":
		if u != nil {
			err = openFile(u.Path)
		} else {
			// 如果u为nil，说明url.Parse失败了，应该作为普通文件路径处理
			err = openFile(arg)
		}

	case "":
		err = openFile(arg)

	default:
		err = openScheme(scheme, arg)
	}
	if err != nil {
		logger.Warning("open failed:", err)
		os.Exit(2)
	}
}

// NOTE: these consts is copied from systemd-go
// https://github.com/coreos/go-systemd/blob/d843340ab4bd3815fda02e648f9b09ae2dc722a7/dbus/dbus.go#L30-L35
const (
	alpha    = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
	num      = `0123456789`
	alphanum = alpha + num
)

// PathBusEscape sanitizes a constituent string of a dbus ObjectPath using the
// rules that systemd uses for serializing special characters.
// NOTE: this function is copied from systemd-go
// https://github.com/coreos/go-systemd/blob/d843340ab4bd3815fda02e648f9b09ae2dc722a7/dbus/dbus.go#L47
func pathBusEscape(path string) string {
	// Special case the empty string
	if len(path) == 0 {
		return "_"
	}
	n := []byte{}
	for i := 0; i < len(path); i++ {
		c := path[i]
		if needsEscape(i, c) {
			e := fmt.Sprintf("_%x", c)
			n = append(n, []byte(e)...)
		} else {
			n = append(n, c)
		}
	}
	return string(n)
}

// needsEscape checks whether a byte in a potential dbus ObjectPath needs to be escaped
// NOTE: this function is copied from systemd-go
// https://github.com/coreos/go-systemd/blob/d843340ab4bd3815fda02e648f9b09ae2dc722a7/dbus/dbus.go#L38
func needsEscape(i int, b byte) bool {
	// Escape everything that is not a-z-A-Z-0-9
	// Also escape 0-9 if it's the first character
	return strings.IndexByte(alphanum, b) == -1 ||
		(i == 0 && strings.IndexByte(num, b) != -1)
}

func getDBusObjectFromAPPID(appId string) (dbus.BusObject, error) {
	sessionBus, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}

	escapeId := pathBusEscape(strings.TrimSuffix(appId, ".desktop"))
	return sessionBus.Object(
		"org.desktopspec.ApplicationManager1",
		dbus.ObjectPath(dbusPath+"/"+escapeId),
	), nil
}

func launchApp(appInfo AppInfo, filename string) error {
	obj, err := getDBusObjectFromAPPID(appInfo.appId)
	if err != nil {
		return err
	}
	err = obj.Call(
		"org.desktopspec.ApplicationManager1.Application.Launch", 0,
		"", []string{filename}, make(map[string]dbus.Variant)).Err
	if err == nil {
		return err
	}

	// NOTE: fallback to use old appmanager

	sessionBus, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	appManager := appmanager.NewManager(sessionBus)
	err = appManager.LaunchApp(dbus.FlagNoAutoStart, appInfo.desktopFile, 0, []string{filename})

	return err
}

func getAppInfo(appInfo *gio.AppInfo) AppInfo {
	dAppInfo := gio.ToDesktopAppInfo(appInfo)
	appId := dAppInfo.GetId()
	desktopFile := dAppInfo.GetFilename()
	logger.Debug("desktop file:", desktopFile)
	info := AppInfo{
		appId:       appId,
		desktopFile: desktopFile,
	}
	return info
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
	logger.Debugf("content type: %q", contentType)

	// 处理空文件的情况，根据文件扩展名猜测 content type
	if contentType == "application/x-zerosize" {
		uncerten, guessedType := gio.ContentTypeGuess(filename, nil)
		// 不确定，使用默认本文打开
		if uncerten || guessedType == "" {
			contentType = "text/plain"
		} else if guessedType == "application/x-zerosize" {
			contentType = "text/plain"
		} else {
			logger.Debugf("empty file detected, guessed content type: %q", guessedType)
			contentType = guessedType

		}
	}

	appInfo := gio.AppInfoGetDefaultForType(contentType, false)
	if appInfo == nil {
		return errors.New("failed to get appInfo")
	}
	defer appInfo.Unref()
	err = launchApp(getAppInfo(appInfo), filename)
	if err != nil {
		return err
	}
	return nil
}

func openScheme(scheme, url string) error {
	logger.Debugf("openScheme: %q, %q", scheme, url)
	appInfo := gio.AppInfoGetDefaultForUriScheme(scheme)
	if appInfo == nil && (scheme == "trash" || scheme == "computer") {
		// As file manager do not register scheme for trash and computer, which are the private protocols
		// supported by gio. To keep compatibility, we forward these protocols to inode/directory's handler.
		// By default, it will be the default file manager like dde-file-manager.
		appInfo = gio.AppInfoGetDefaultForType("inode/directory", true)
	}
	if appInfo == nil {
		return errors.New("failed to get appInfo")
	}
	defer appInfo.Unref()

	err := launchApp(getAppInfo(appInfo), url)
	if err != nil {
		return err
	}
	return nil
}
