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
	"fmt"
	"regexp"

	"github.com/godbus/dbus"
	"pkg.deepin.io/lib/users/passwd"
)

const (
	UsernameMaxLength = 32
	UidMin            = 1000
)

const (
	UsernameOk = 0
	UsernameEmpty
	UsernameInvalidChars
	UsernameFirstCharInvalid
	UsernameExists
	UsernameSystemUsed
	UsernameExceed
)

var usernameRegex *regexp.Regexp

func init() {
	usernameRegex, _ = regexp.Compile("[a-z][a-z0-9_-]+")
}

// ValidateUsername valiate username based on the following rules:
// * Length of username is (0, 32].
// * First char of username bust be in lower case.
// * Username can only contain lower letters(a~z), numbers(0~9), dash(-) and
//   userscore(_).
// * Username cannot be used by others or by system.
func (validator *Validator) ValidateUsername(username string) (code int, result string, busErr *dbus.Error) {
	if len(username) == 0 {
		return UsernameEmpty, "Username can not be empty.", nil
	}

	if len(username) > UsernameMaxLength {
		return UsernameExceed, fmt.Sprintf("The length of username cannot exceed %d characters", UsernameMaxLength), nil
	}

	if pw, _ := passwd.GetPasswdByName(username); pw != nil {
		if pw.Uid >= UidMin {
			return UsernameExists, "The username exists.", nil
		} else {
			return UsernameSystemUsed, "The username has been used by system.", nil
		}
	}

	if username[0] < 'a' || username[0] > 'z' {
		return UsernameFirstCharInvalid, "The first character msut be in lower case.", nil
	}

	if !usernameRegex.MatchString(username) {
		return UsernameInvalidChars, "Username must comprise a~z, 0-9, - or _", nil
	}

	return UsernameOk, "", nil
}
