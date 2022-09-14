// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"regexp"

	"github.com/godbus/dbus"
	"github.com/linuxdeepin/go-lib/users/passwd"
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
