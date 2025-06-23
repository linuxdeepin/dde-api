// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"os/exec"
	"sync"
	"time"

	"github.com/linuxdeepin/go-lib/dbusutil"
	"github.com/linuxdeepin/go-lib/log"
	dutils "github.com/linuxdeepin/go-lib/utils"
)

//go:generate dbusutil-gen em -type Helper

const (
	dbusServiceName       = "org.deepin.dde.LocaleHelper1"
	dbusPath              = "/org/deepin/dde/LocaleHelper1"
	dbusInterface         = dbusServiceName
	localeGenBin          = "/usr/sbin/locale-gen"
	deepinImmutableCtlBin = "/usr/sbin/deepin-immutable-ctl"
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
	if !dutils.IsFileExist(deepinImmutableCtlBin) {
		logger.Warning("deepin-immutable-ctl not found, use locale-gen directly")
		return exec.Command(localeGenBin).Run()
	} else {
		// TODO 在磐石适配 locale-gen 前使用 deepin-immutable-ctl 执行 locale-gen，否则有权限问题
		output, err := exec.Command(deepinImmutableCtlBin, "admin", "exec", localeGenBin).CombinedOutput()
		if err != nil {
			logger.Warning("deepin-immutable-ctl exec locale-gen failed, err:", err, "output:", string(output))
			return err
		}
		return nil
	}
}

// locales version <= 2.13
func (h *Helper) doGenLocaleWithParam(locale string) error {
	if !dutils.IsFileExist(deepinImmutableCtlBin) {
		logger.Warning("deepin-immutable-ctl not found, use locale-gen directly")
		return exec.Command(localeGenBin, locale).Run()
	} else {
		// TODO 在磐石适配 locale-gen 前使用 deepin-immutable-ctl 执行 locale-gen，否则有权限问题
		output, err := exec.Command(deepinImmutableCtlBin, "admin", "exec", "--", localeGenBin, locale).CombinedOutput()
		if err != nil {
			logger.Warning("deepin-immutable-ctl exec locale-gen failed, err:", err, "output:", string(output))
			return err
		}
		return nil
	}
}
