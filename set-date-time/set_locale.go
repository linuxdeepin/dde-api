/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

import (
	"os/exec"
)

const (
	LOCALE_GEN_CMD = "/usr/sbin/locale-gen"
)

func (op *SetDateTime) GenLocale(locale string) {
	if len(locale) < 1 {
		logger.Warning("GenLocale Arg Error")
		op.GenLocaleStatus(false, "Arg Error")
		return
	}

	go func() {
		err := exec.Command(LOCALE_GEN_CMD, locale).Run()
		if err != nil {
			logger.Warningf("locale-gen '%s' failed: %v\n",
				locale, err)
			op.GenLocaleStatus(false, err.Error())
			return
		}
		op.GenLocaleStatus(true, locale)
	}()
}
