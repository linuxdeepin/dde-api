/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

import (
	"fmt"
	"regexp"

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
func (validator *Validator) ValidateUsername(username string) (int, string) {
	if len(username) == 0 {
		return UsernameEmpty, "Username can not be empty."
	}

	if len(username) > UsernameMaxLength {
		return UsernameExceed, fmt.Sprintf("The length of username cannot exceed %s characters", UsernameMaxLength)
	}

	if pw, _ := passwd.GetPasswdByName(username); pw != nil {
		if pw.Uid >= UidMin {
			return UsernameExists, "The username exists."
		} else {
			return UsernameSystemUsed, "The username has been used by system."
		}
	}

	if username[0] < 'a' || username[0] > 'z' {
		return UsernameFirstCharInvalid, "The first character msut be in lower case."
	}

	if !usernameRegex.MatchString(username) {
		return UsernameInvalidChars, "Username must comprise a~z, 0-9, - or _"
	}

	return UsernameOk, ""
}
