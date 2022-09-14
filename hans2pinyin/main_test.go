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

type UnitTestSuite struct {
	suite.Suite
	m *Manager
}

func (s *UnitTestSuite) SetupSuite() {
	var err error
	s.m = &Manager{}
	s.m.service, err = dbusutil.NewSessionService()
	if err != nil {
		s.T().Skip(fmt.Sprintf("failed to get service: %v", err))
	}
}

func (s *UnitTestSuite) Test_Query() {
	pinyin, err := s.m.Query("Hanz")
	s.Require().Nil(err)
	s.Assert().ElementsMatch(pinyin, []string{"Hanz"})
}

func (s *UnitTestSuite) Test_QueryQueryList() {
	hansList := []string{"Hanz"}
	jsonStr, err := s.m.QueryList(hansList)
	testStr := `{"Hanz":["Hanz"]}`
	s.Require().Nil(err)
	s.Equal(jsonStr, testStr)
}

func (s *UnitTestSuite) Test_GetInterfaceName() {
	s.Equal(s.m.GetInterfaceName(), dbusServiceName)
}

func (s *UnitTestSuite) Test_usage() {
	usage()
}

func (s *UnitTestSuite) Test_queryPinyin() {
	pinyin := queryPinyin("Hanz")
	s.Assert().ElementsMatch(pinyin, []string{"Hanz"})
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
