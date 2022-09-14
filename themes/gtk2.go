// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package themes

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
)

const (
	gtk2ConfDelim = "="
)

type gtk2ConfInfo struct {
	key   string
	value string
}
type gtk2ConfInfos []*gtk2ConfInfo

var (
	gtk2Locker   sync.Mutex
	gtk2ConfFile = path.Join(os.Getenv("HOME"), ".gtkrc-2.0")
)

func setGtk2Theme(name string) error {
	return setGtk2Prop("gtk-theme-name",
		"\""+name+"\"", gtk2ConfFile)
}

func setGtk2Icon(name string) error {
	return setGtk2Prop("gtk-icon-theme-name",
		"\""+name+"\"", gtk2ConfFile)
}

func setGtk2Cursor(name string) error {
	return setGtk2Prop("gtk-cursor-theme-name",
		"\""+name+"\"", gtk2ConfFile)
}

func setGtk2Prop(key, value, file string) error {
	gtk2Locker.Lock()
	defer gtk2Locker.Unlock()

	infos := gtk2FileReader(file)
	info := infos.Get(key)
	if info == nil {
		infos = infos.Add(key, value)
	} else {
		if info.value == value {
			return nil
		}
		info.value = value
	}
	return gtk2FileWriter(infos, file)
}

func (infos gtk2ConfInfos) Get(key string) *gtk2ConfInfo {
	for _, info := range infos {
		if info.key == key {
			return info
		}
	}
	return nil
}

func (infos gtk2ConfInfos) Add(key, value string) gtk2ConfInfos {
	for _, info := range infos {
		if info.key == key {
			info.value = value
			return infos
		}
	}

	infos = append(infos, &gtk2ConfInfo{
		key:   key,
		value: value,
	})
	return infos
}

func gtk2FileWriter(infos gtk2ConfInfos, file string) error {
	var content string
	length := len(infos)
	for i, info := range infos {
		content += info.key + gtk2ConfDelim + info.value
		if i != length-1 {
			content += "\n"
		}
	}

	err := os.MkdirAll(path.Dir(file), 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, []byte(content), 0644)
}

func gtk2FileReader(file string) gtk2ConfInfos {
	var infos gtk2ConfInfos
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return infos
	}

	var lines = strings.Split(string(content), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		array := strings.Split(line, gtk2ConfDelim)
		if len(array) != 2 {
			continue
		}

		infos = append(infos, &gtk2ConfInfo{
			key:   array[0],
			value: array[1],
		})
	}
	return infos
}
