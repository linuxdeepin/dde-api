/**
 * Copyright (c) 2015 Deepin, Inc.
 *               2015 Xu Shaohua
 *
 * Author:       Xu Shaohua<xushaohua@linuxdeepin.com>
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
	"testing"

	C "launchpad.net/gocheck"
)

func Test(t *testing.T) {
	C.TestingT(t)
}

var _ = C.Suite(&Validator{})

func (validator *Validator) TestValidateHostname(c *C.C) {
	c.Check(validator.ValidateHostname("hostname"), C.Equals, true)
	c.Check(validator.ValidateHostname("#1"), C.Equals, false)
	c.Check(validator.ValidateHostname("sub.domain."), C.Equals, false)
	c.Check(validator.ValidateHostnameTemp("sub.domain."), C.Equals, true)
	c.Check(validator.ValidateHostnameTemp("sub-domain."), C.Equals, true)
	c.Check(validator.ValidateHostnameTemp("sub-domain$"), C.Equals, false)
}
