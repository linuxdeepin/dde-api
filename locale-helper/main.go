// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/linuxdeepin/go-lib/dbusutil"
	"github.com/linuxdeepin/go-lib/log"
)

//go:generate dbusutil-gen em -type Helper

const (
	dbusServiceName = "org.deepin.dde.LocaleHelper1"
	dbusPath        = "/org/deepin/dde/LocaleHelper1"
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
