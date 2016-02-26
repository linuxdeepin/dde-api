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
	"path"

	"gir/gio-2.0"
	"pkg.deepin.io/lib/mime"
	dutils "pkg.deepin.io/lib/utils"
)

type AppInfo struct {
	// Desktop id
	Id string
	// App name
	Name string
	// Display name
	DisplayName string
	// Comment
	Description string
	// Icon
	Icon string
	// Commandline
	Exec string
}
type AppInfos []*AppInfo

func GetAppInfo(ty string) (*AppInfo, error) {
	id, err := mime.GetDefaultApp(ty, false)
	if err != nil {
		return nil, err
	}

	return newAppInfoById(id), nil
}

func (infos AppInfos) Delete(id string) AppInfos {
	var ret AppInfos
	for _, info := range infos {
		if info.Id == id {
			continue
		}
		ret = append(ret, info)
	}
	return ret
}

func SetAppInfo(ty, id string) error {
	return mime.SetDefaultApp(ty, id)
}

func GetAppInfos(ty string) AppInfos {
	var infos AppInfos
	for _, id := range mime.GetAppList(ty) {
		infos = append(infos, newAppInfoById(id))
	}
	return infos
}

func newAppInfoById(id string) *AppInfo {
	ginfo := gio.NewDesktopAppInfo(id)
	defer ginfo.Unref()
	var info = &AppInfo{
		Id:          id,
		Name:        ginfo.GetName(),
		DisplayName: ginfo.GetGenericName(),
		Description: ginfo.GetDescription(),
		Exec:        ginfo.GetCommandline(),
	}
	iconObj := ginfo.GetIcon()
	if iconObj != nil {
		info.Icon = iconObj.ToString()
		iconObj.Unref()
	}

	return info
}

func findFilePath(file string) string {
	data := path.Join(os.Getenv("HOME"), ".local/share", file)
	if dutils.IsFileExist(data) {
		return data
	}

	data = path.Join("/usr/local/share", file)
	if dutils.IsFileExist(data) {
		return data
	}

	return path.Join("/usr/share", file)
}
