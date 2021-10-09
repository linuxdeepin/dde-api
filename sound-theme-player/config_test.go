/*
 * Copyright (C) 2020 ~ 2022 Uniontech Software Technology Co.,Ltd
 *
 * Author:     dengbo <dengbo@uniontech.com>
 *
 * Maintainer: dengbo <dengbo@uniontech.com>
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
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UnitTestSuite struct {
	suite.Suite
	filename string
}

func (s *UnitTestSuite) SetupSuite() {
	homeDir = "./"
	s.filename = "config-1000.json"
	data := []byte(`{"DesktopLoginEnabled":false,"Theme":"deepin","Card":"PCH","Device":"0","Mute":false}`)
	err := ioutil.WriteFile(s.filename, data, 0777)
	require.NoError(s.T(), err)
}

func (s *UnitTestSuite) TearDownSuite() {
	homeDir = ""
	err := os.Remove(s.filename)
	require.NoError(s.T(), err)
}

func (s *UnitTestSuite) Test_getConfigFile() {
	file := getConfigFile(int(1000))
	s.Assert().Equal(file, "config-1000.json")
}

func (s *UnitTestSuite) Test_loadUserConfig() {
	var testCfg config

	type test struct {
		name      string
		uid       int
		isSuccess bool
	}

	tests := []test{
		{name: "File Exists", uid: 1000, isSuccess: true},
		{name: "File Not Exists", uid: 2000, isSuccess: false},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			if tt.isSuccess {
				s.Nil(loadUserConfig(tt.uid, &testCfg))
			} else {
				s.Error(loadUserConfig(tt.uid, &testCfg))
			}
		})
	}
}

func (s *UnitTestSuite) Test_saveUserConfig() {
	var testCfg config

	s.Nil(saveUserConfig(int(1000), &testCfg))
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
