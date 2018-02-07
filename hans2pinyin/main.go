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

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/pinyin"
	"time"
)

const (
	dbusDest = "com.deepin.api.Pinyin"
	dbusPath = "/com/deepin/api/Pinyin"
	dbusIFC  = dbusDest
)

type Manager struct{}

func (*Manager) Query(hans string) []string {
	return queryPinyin(hans)
}

// Querylist query pinyin for hans list, return a json data.
func (*Manager) QueryList(hansList []string) string {
	var data = make(map[string][]string)
	for _, hans := range hansList {
		data[hans] = queryPinyin(hans)
	}

	content, _ := json.Marshal(data)
	return string(content)
}

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func main() {
	if len(os.Args) == 2 {
		if os.Args[1] == "--help" || os.Args[1] == "-h" {
			usage()
			return
		}

		fmt.Println(queryPinyin(os.Args[1]))
		return
	}

	err := dbus.InstallOnSession(new(Manager))
	if err != nil {
		fmt.Println("Install dbus failed:", err)
		return
	}

	dbus.DealWithUnhandledMessage()
	dbus.SetAutoDestroyHandler(time.Second*5, nil)
	err = dbus.Wait()
	if err != nil {
		fmt.Println("Lost dbus:", err)
		os.Exit(-1)
	}

	os.Exit(0)
}

func usage() {
	fmt.Println("Usage: hans2pinyin <hans>")
	fmt.Println("Example:")
	fmt.Println("\thans2pinyin 重")
	fmt.Println("\t[zhong chong]")
}

func queryPinyin(hans string) []string {
	return pinyin.HansToPinyin(hans)
}
