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
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"

	"pkg.deepin.io/lib/dbus"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	dbusDest = "com.deepin.api.GreeterHelper"
	dbusPath = "/com/deepin/api/GreeterHelper"
	dbusIFC  = "com.deepin.api.GreeterHelper"

	defaultConfig = "/var/lib/greeter/users.ini"

	kfKeyTheme      = "GreeterTheme"
	kfKeyLayout     = "KeyboardLayout"
	kfKeyLayoutList = "KeyboardLayoutList"

	layoutDelim = "|"
	listDelim   = " "
)

type Manager struct {
	locker sync.Mutex
}

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func (m *Manager) set(group, key, value string) error {
	m.locker.Lock()
	defer m.locker.Unlock()
	return doSet(defaultConfig, group, key, value)
}

func formatLayoutList(list []string) string {
	var ret string
	for _, v := range list {
		tmp := formatLayout(v)
		if len(tmp) == 0 {
			continue
		}
		ret += tmp + listDelim
	}
	return strings.TrimSpace(ret)
}

func formatLayout(layout string) string {
	array := strings.Split(layout, ";")
	if len(array) != 2 {
		return ""
	}

	layout = array[0] + layoutDelim
	if len(array[1]) > 0 {
		layout += array[1]
	}
	return layout
}

func doSet(file, group, key, value string) error {
	err := ensureConfigExist(file)
	if err != nil {
		return err
	}

	kf, err := dutils.NewKeyFileFromFile(file)
	if err != nil {
		return err
	}
	defer kf.Free()

	v, _ := kf.GetString(group, key)
	if v == value {
		return nil
	}

	kf.SetString(group, key, value)
	_, content, err := kf.ToData()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, []byte(content), 0644)
}

func ensureConfigExist(file string) error {
	if dutils.IsFileExist(file) {
		return nil
	}

	err := os.MkdirAll(path.Dir(file), 0755)
	if err != nil {
		return err
	}
	return dutils.CreateFile(file)
}
