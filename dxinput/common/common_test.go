// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Get(t *testing.T) {
	info := DeviceInfo{
		Id:      111,
		Type:    1,
		Name:    "test",
		Enabled: true,
	}

	infos := DeviceInfos{
		&info,
	}

	assert.Equal(t, infos.Get(111).Id, info.Id)
}
