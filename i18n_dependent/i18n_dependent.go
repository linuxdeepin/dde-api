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

package i18n_dependent

type DependentInfo struct {
	Dependent string
	Packages  []string
}
type DependentInfos []DependentInfo

const (
	pkgDependsFile = "/usr/share/i18n/i18n_dependent.json"
)

func GetByPackage(locale, pkg string) ([]string, []string, error) {
	categories, err := getDependentCategories(pkgDependsFile)
	if err != nil {
		return nil, nil, err
	}

	infos := categories.GetAllDependentInfos(locale)
	pkgs := infos.GetPackagesByDependent(pkg)
	return pkgs, getConflictPackages(pkgs), nil
}

func GetByLocale(locale string) (DependentInfos, DependentInfos, error) {
	categories, err := getDependentCategories(pkgDependsFile)
	if err != nil {
		return nil, nil, err
	}

	infos := categories.GetAllDependentInfos(locale)
	return infos, infos.GetConflictPackages(), nil
}

func (infos DependentInfos) GetPackagesByDependent(dependent string) []string {
	var list []string
	for _, info := range infos {
		if info.Dependent != dependent {
			continue
		}
		list = append(list, info.Packages...)
	}
	return list
}

func (infos DependentInfos) GetConflictPackages() DependentInfos {
	var conflicts DependentInfos
	for _, info := range infos {
		list := getConflictPackages(info.Packages)
		if len(list) == 0 {
			continue
		}
		conflicts = append(conflicts, DependentInfo{
			Dependent: info.Dependent,
			Packages:  list,
		})
	}
	return conflicts
}

func getConflictPackages(pkgs []string) []string {
	var list []string
	for _, pkg := range pkgs {
		v, ok := conflictPkgMap[pkg]
		if !ok {
			continue
		}
		list = append(list, v...)
	}
	return list
}
