/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dxinput

import (
	"pkg.deepin.io/dde/api/dxinput/utils"
)

func SetKeyboardRepeat(enabled bool, delay, interval uint32) error {
	return utils.SetKeyboardRepeat(enabled, delay, interval)
}
