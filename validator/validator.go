// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"github.com/linuxdeepin/go-lib/dbusutil"
)

//go:generate dbusutil-gen em -type Validator

type Validator struct {
	service *dbusutil.Service
}

func (v *Validator) GetInterfaceName() string {
	return DBusInterface
}
