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
	"fmt"
	"testing"

	"pkg.deepin.io/lib/dbusutil"

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
