// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package dxinput

import (
	"github.com/linuxdeepin/dde-api/dxinput/utils"
)

func SetKeyboardRepeat(enabled bool, delay, interval uint32) error {
	return utils.SetKeyboardRepeat(enabled, delay, interval)
}
