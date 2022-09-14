// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"log"

	"github.com/godbus/dbus"
	"github.com/linuxdeepin/go-lib/dbusutil"
	"github.com/linuxdeepin/go-lib/pinyin"
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
