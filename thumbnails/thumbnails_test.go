// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package thumbnails

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/linuxdeepin/dde-api/thumbnails/loader"
)

func TestCorrectSize(t *testing.T) {
	assert.Equal(t, correctSize(64), loader.SizeFlagSmall)
	assert.Equal(t, correctSize(128), loader.SizeFlagNormal)
	assert.Equal(t, correctSize(176), loader.SizeFlagNormal)
	assert.Equal(t, correctSize(256), loader.SizeFlagLarge)
}
