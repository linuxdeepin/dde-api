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

func (s *UnitTestSuite) Test_GetInterfaceName() {
	s.m.GetInterfaceName()
}

func (s *UnitTestSuite) Test_GetLunarInfoBySolar() {
	_, _, err := s.m.GetLunarInfoBySolar(2021, 10, 1)
	s.Require().Nil(err)
}

func (s *UnitTestSuite) Test_GetFestivalsInRange() {
	_, err := s.m.GetFestivalsInRange("2021-01-02", "2021-10-01")
	s.Require().Nil(err)
}

func (s *UnitTestSuite) Test_GetLunarMonthCalendar() {
	_, _, err := s.m.GetLunarMonthCalendar(2021, 10, true)
	s.Require().Nil(err)
}

func (s *UnitTestSuite) Test_GetHuangLiDay() {
	_, err := s.m.GetHuangLiDay(2021, 10, 1)
	s.Require().Nil(err)
}

func (s *UnitTestSuite) Test_GetHuangLiMonth() {
	_, err := s.m.GetHuangLiMonth(2021, 10, true)
	s.Require().Nil(err)
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
