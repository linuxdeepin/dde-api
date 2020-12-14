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
	"time"

	"log"

	"github.com/godbus/dbus"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/pinyin"
)

//go:generate dbusutil-gen em -type Manager

const (
	dbusServiceName = "com.deepin.api.Pinyin"
	dbusPath        = "/com/deepin/api/Pinyin"
	dbusInterface   = dbusServiceName
)

type Manager struct {
	service *dbusutil.Service
}

func (m *Manager) Query(hans string) (pinyin []string, busErr *dbus.Error) {
	m.service.DelayAutoQuit()
	return queryPinyin(hans), nil
}

// QueryList query pinyin for hans list, return a json data.
func (m *Manager) QueryList(hansList []string) (jsonStr string, err *dbus.Error) {
	m.service.DelayAutoQuit()
	var data = make(map[string][]string)
	for _, hans := range hansList {
		data[hans] = queryPinyin(hans)
	}

	content, _ := json.Marshal(data)
	return string(content), nil
}

func (*Manager) GetInterfaceName() string {
	return dbusInterface
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

	service, err := dbusutil.NewSessionService()
	if err != nil {
		log.Fatal("failed to new session service", err)
	}

	hasOwner, err := service.NameHasOwner(dbusServiceName)
	if err != nil {
		log.Fatal("failed to call NameHasOwner:", err)
	}
	if hasOwner {
		log.Fatalf("name %q already has the owner", dbusServiceName)
	}

	m := &Manager{
		service: service,
	}
	err = service.Export(dbusPath, m)
	if err != nil {
		log.Fatal("failed to export:", err)
	}
	err = service.RequestName(dbusServiceName)
	if err != nil {
		log.Fatal("failed to request name:", err)
	}
	service.SetAutoQuitHandler(time.Second*5, nil)
	service.Wait()
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
