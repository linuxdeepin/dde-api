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
	"os"
	"os/exec"
	"sync"
	"time"

	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/log"
)

//go:generate dbusutil-gen em -type Helper

const (
	dbusServiceName = "com.deepin.api.LocaleHelper"
	dbusPath        = "/com/deepin/api/LocaleHelper"
	dbusInterface   = dbusServiceName
	localeGenBin    = "locale-gen"
)

type Helper struct {
	service *dbusutil.Service
	mu      sync.Mutex
	running bool

	//nolint
	signals *struct {
		/**
		 * if failed, Success(false, reason), else Success(true, "")
		 **/
		Success struct {
			ok     bool
			reason string
		}
	}
}

var (
	logger = log.NewLogger(dbusServiceName)
)

func main() {
	os.Setenv("PATH", "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")
	logger.BeginTracing()
	defer logger.EndTracing()

	service, err := dbusutil.NewSystemService()
	if err != nil {
		logger.Fatal("failed to new system service:", err)
	}

	hasOwner, err := service.NameHasOwner(dbusServiceName)
	if err != nil {
		logger.Fatal(err)
	}
	if hasOwner {
		logger.Fatalf("name %q already has the owner", dbusServiceName)
	}

	var h = &Helper{
		running: false,
		service: service,
	}
	err = service.Export(dbusPath, h)
	if err != nil {
		logger.Fatal("failed to export:", err)
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		logger.Fatal("failed to request name:", err)
	}

	service.SetAutoQuitHandler(30*time.Second, h.canQuit)
	service.Wait()
}

func (*Helper) GetInterfaceName() string {
	return dbusInterface
}

func (h *Helper) canQuit() bool {
	h.mu.Lock()
	running := h.running
	h.mu.Unlock()
	return !running
}

func (h *Helper) doGenLocale() error {
	return exec.Command(localeGenBin).Run()
}

// locales version <= 2.13
func (h *Helper) doGenLocaleWithParam(locale string) error {
	return exec.Command(localeGenBin, locale).Run()
}
