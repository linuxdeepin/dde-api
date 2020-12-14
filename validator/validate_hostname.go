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
	"regexp"

	"github.com/godbus/dbus"
)

var hostnameRegex *regexp.Regexp
var hostnameTempRegex *regexp.Regexp

func init() {
	hostnameRegex, _ = regexp.Compile("^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]).)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9-]*[A-Za-z0-9])$")
	hostnameTempRegex, _ = regexp.Compile("^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9-.])*$")
}

// ValidateUsername valites user name based on following rules:
// * Only letters(a~z, A~z), numbers(0~9), dash(-) and underscore(_) are allowed.
// * Only lowercase of letters are allowed as the first character
// * Maximum size is 32
// * Username cannot be empty
//func (validator *Validator) ValidateUsername(username string) bool {
//}

// ValidateHostname validates hostname (machine name).
// Only letters(a~z, A~Z) and numbers(0~9) are allowed as prefix/suffix of
// hostname.
// Dot(.) is used to separator domain and subdomains
// Underscore(_) and dash(-) are used to concat letters and numbers.
func (validator *Validator) ValidateHostname(hostname string) (result bool, busErr *dbus.Error) {
	return hostnameRegex.MatchString(hostname), nil
}

// ValidateHostnameTemp validates part of hostname.
// This function is used to check hostname when it is being input.
// Unlike @ValidateHostname, dot(.), underscore(_) and dash(-) are allowed to
// be at the end of hostname.
func (validator *Validator) ValidateHostnameTemp(hostname string) (result bool, busErr *dbus.Error) {
	return hostnameTempRegex.MatchString(hostname), nil
}
