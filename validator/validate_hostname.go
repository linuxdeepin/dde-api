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
	"regexp"
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
func (validator *Validator) ValidateHostname(hostname string) bool {
	return hostnameRegex.MatchString(hostname)
}

// ValidateHostnameTemp validates part of hostname.
// This function is used to check hostname when it is being input.
// Unlike @ValidateHostname, dot(.), underscore(_) and dash(-) are allowed to
// be at the end of hostname.
func (validator *Validator) ValidateHostnameTemp(hostname string) bool {
	return hostnameTempRegex.MatchString(hostname)
}
