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
	"testing"

	C "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	C.TestingT(t)
}

var _ = C.Suite(&Validator{})

func (validator *Validator) TestValidateHostname(c *C.C) {
	var res bool

	res, _ = validator.ValidateHostname("hostname")
	c.Check(res, C.Equals, true)

	res, _ = validator.ValidateHostname("#1")
	c.Check(res, C.Equals, false)

	res, _ = validator.ValidateHostname("sub.domain.")
	c.Check(res, C.Equals, false)
}

func (validator *Validator) TestValidateHostnameTemp(c *C.C) {
	var res bool

	res, _ = validator.ValidateHostnameTemp("sub.domain.")
	c.Check(res, C.Equals, true)

	res, _ = validator.ValidateHostnameTemp("sub.domain$")
	c.Check(res, C.Equals, false)
}

func (validator *Validator) TestValiateUsername(c *C.C) {
	state, _, _ := validator.ValidateUsername("root")
	c.Check(state, C.Equals, UsernameSystemUsed)

	state, _, _ = validator.ValidateUsername("nonexst")
	c.Check(state, C.Equals, UsernameOk)

	state, _, _ = validator.ValidateUsername("-first-char")
	c.Check(state, C.Equals, UsernameFirstCharInvalid)

	state, _, _ = validator.ValidateUsername("upperCase")
	c.Check(state, C.Equals, UsernameInvalidChars)
}
