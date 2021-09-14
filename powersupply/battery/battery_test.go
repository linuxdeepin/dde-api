/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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
