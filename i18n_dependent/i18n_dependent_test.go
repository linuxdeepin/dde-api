// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package i18n_dependent

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UnitTestSuite struct {
	suite.Suite
	testCategories     jsonDependentCategories
	testDependentInfos jsonDependentInfos
}

func (s *UnitTestSuite) SetupSuite() {
	var err error
	s.testCategories, err = getDependentCategories("testdata/i18n_dependent.json")
	s.Require().Nil(err)
	s.testDependentInfos = s.testCategories[0].Infos
}

func (s *UnitTestSuite) Test_GetAllDependentInfos() {
	testDependents := s.testCategories.GetAllDependentInfos("zh_HK.UTF-8")
	s.Assert().NotNil(testDependents)
}

func (s *UnitTestSuite) Test_GetDependentInfos() {
	testDependents := s.testCategories.GetDependentInfos("tr", "zh_HK.UTF-8")
	s.Assert().NotNil(testDependents)
}

func (s *UnitTestSuite) Test_GetInfos() {
	testDependents := s.testCategories.GetInfos("tr")
	s.Assert().NotNil(testDependents)
}

func (s *UnitTestSuite) Test_GetDependentInfosByLocale() {
	s.testDependentInfos.GetDependentInfos("zh_HK.UTF-8")
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
