// Copyright (c) 2015 Deepin Ltd. All rights reserved.
// Use of this source is govered by General Public License that can be found
// in the LICENSE file.
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
// * Username cannot already be userd others or by system.
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
