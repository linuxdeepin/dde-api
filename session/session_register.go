/*
 * Copyright (C) 2015 ~ 2018 Deepin Technology Co., Ltd.
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

package session

import (
	"fmt"
	"os"

	"github.com/linuxdeepin/go-dbus-factory/com.deepin.sessionmanager"
	"github.com/godbus/dbus"
	"pkg.deepin.io/lib/utils"
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
