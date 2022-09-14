// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package themes

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeThemeList(t *testing.T) {
	src := []string{"Deepin", "Adwaita", "Zukitwo"}
	target := []string{"Deepin", "Evolve"}
	ret := []string{"Deepin", "Adwaita", "Zukitwo", "Evolve"}

	assert.ElementsMatch(t, mergeThemeList(src, target), ret)
}

func TestSetQt4Theme(t *testing.T) {
	config := "/tmp/Trolltech.conf"
	assert.Equal(t, setQt4Theme(config), true)
	os.Remove(config)
}
