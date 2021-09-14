/*
 * Copyright (C) 2019 ~ 2021 Uniontech Software Technology Co.,Ltd
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
