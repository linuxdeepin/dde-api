// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package battery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseStatus(t *testing.T) {
	assert.Equal(t, parseStatus("Unknown"), StatusUnknown)
	assert.Equal(t, parseStatus("Charging"), StatusCharging)
	assert.Equal(t, parseStatus("Discharging"), StatusDischarging)
	assert.Equal(t, parseStatus("Not charging"), StatusNotCharging)
	assert.Equal(t, parseStatus("Full"), StatusFull)
	assert.Equal(t, parseStatus("Other"), StatusUnknown)
}

func Test_GetDisplayStatus(t *testing.T) {
	// one
	one := []Status{StatusDischarging}
	assert.Equal(t, GetDisplayStatus(one), StatusDischarging)
	one[0] = StatusNotCharging
	assert.Equal(t, GetDisplayStatus(one), StatusNotCharging)

	// two
	two := []Status{StatusFull, StatusFull}
	assert.Equal(t, GetDisplayStatus(two), StatusFull)
	two[0] = StatusDischarging
	two[1] = StatusFull
	assert.Equal(t, GetDisplayStatus(two), StatusDischarging)

	two[0] = StatusCharging
	two[1] = StatusFull
	assert.Equal(t, GetDisplayStatus(two), StatusCharging)

	two[0] = StatusCharging
	two[1] = StatusDischarging
	assert.Equal(t, GetDisplayStatus(two), StatusDischarging)
}
