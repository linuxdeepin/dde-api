// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"testing"

	"github.com/linuxdeepin/go-lib/dbusutil"

	"github.com/stretchr/testify/suite"
)

type UnitTest struct {
	suite.Suite
	validator *Validator
}

func (ut *UnitTest) SetupSuite() {
	var err error
	ut.validator = &Validator{}
	ut.validator.service, err = dbusutil.NewSessionService()
	if err != nil {
		ut.T().Skip(fmt.Sprintf("failed to get service: %v", err))
	}
}

func (ut *UnitTest) TestValidateHostname() {
	var res bool

	res, _ = ut.validator.ValidateHostname("hostname")
	ut.Equal(res, true)

	res, _ = ut.validator.ValidateHostname("#1")
	ut.Equal(res, false)

	res, _ = ut.validator.ValidateHostname("sub.domain.")
	ut.Equal(res, false)
}

func (ut *UnitTest) TestValidateHostnameTemp() {
	var res bool

	res, _ = ut.validator.ValidateHostnameTemp("sub.domain.")
	ut.Equal(res, true)

	res, _ = ut.validator.ValidateHostnameTemp("sub.domain$")
	ut.Equal(res, false)
}

func (ut *UnitTest) TestValiateUsername() {
	state, _, _ := ut.validator.ValidateUsername("root")
	ut.Equal(state, UsernameSystemUsed)

	state, _, _ = ut.validator.ValidateUsername("nonexst")
	ut.Equal(state, UsernameOk)

	state, _, _ = ut.validator.ValidateUsername("-first-char")
	ut.Equal(state, UsernameFirstCharInvalid)

	state, _, _ = ut.validator.ValidateUsername("upperCase")
	ut.Equal(state, UsernameInvalidChars)
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTest))
}
