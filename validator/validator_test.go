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
	"testing"

	C "gopkg.in/check.v1"
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
