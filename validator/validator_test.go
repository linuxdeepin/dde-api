// Copyright (c) 2015 Deepin Ltd. All rights reserved.
// Use of this source is govered by General Public License that can be found
// in the LICENSE file.
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
}

func (validator *Validator) TestValidateHostnameTemp(c *C.C) {
	c.Check(validator.ValidateHostnameTemp("sub.domain."), C.Equals, true)
	c.Check(validator.ValidateHostnameTemp("sub-domain."), C.Equals, true)
	c.Check(validator.ValidateHostnameTemp("sub-domain$"), C.Equals, false)
}

func (validator *Validator) TestValiateUsername(c *C.C) {
	state, _ := validator.ValidateUsername("root")
	c.Check(state, C.Equals, UsernameSystemUsed)

	state, _ = validator.ValidateUsername("nonexst")
	c.Check(state, C.Equals, UsernameOk)

	state, _ = validator.ValidateUsername("-first-char")
	c.Check(state, C.Equals, UsernameFirstCharInvalid)

	state, _ = validator.ValidateUsername("upperCase")
	c.Check(state, C.Equals, UsernameInvalidChars)
}
