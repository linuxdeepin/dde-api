// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package session

import (
	"fmt"
	"os"

	"github.com/godbus/dbus/v5"
	sessionmanager "github.com/linuxdeepin/go-dbus-factory/session/org.deepin.dde.sessionmanager1"
	"github.com/linuxdeepin/go-lib/utils"
)

// Register will register to session manager if program is started from startdde.
func Register() {
	cookie := os.ExpandEnv("$DDE_SESSION_PROCESS_COOKIE_ID")
	err := utils.UnsetEnv("DDE_SESSION_PROCESS_COOKIE_ID")

	if cookie == "" {
		fmt.Println("get DDE_SESSION_PROCESS_COOKIE_ID failed")
		return
	}

	if err != nil {
		fmt.Println("unsetenv DDE_SESSION_PROCESS_COOKIE_ID failed")
	}

	go func() {
		sessionBus, err := dbus.SessionBus()
		if err != nil {
			fmt.Println("failed to get session bus:", err)
			return
		}
		manager := sessionmanager.NewSessionManager(sessionBus)
		_, err = manager.Register(dbus.FlagNoAutoStart, cookie)
		if err != nil {
			fmt.Println("failed to register:", err)
		}
	}()
}
