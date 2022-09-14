// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package powersupply

import (
	"testing"
)

func TestACOnline(t *testing.T) {
	acExist, acOnline, err := ACOnline()
	t.Logf("acExist %v, acOnline %v, err %v", acExist, acOnline, err)
}

func TestGetSystemBatteryInfos(t *testing.T) {
	batInfos, err := GetSystemBatteryInfos()
	if err != nil {
		t.Log("err", err)
		return
	}
	for _, batInfo := range batInfos {
		t.Logf("%+v", batInfo)
	}
}
