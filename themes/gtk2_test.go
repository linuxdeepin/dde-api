// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package themes

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGtk2Infos(t *testing.T) {
	infos := gtk2FileReader("testdata/gtkrc-2.0")
	assert.Equal(t, len(infos), 16)

	info := infos.Get("gtk-theme-name")
	assert.Equal(t, info.value, "\"Paper\"")

	info.value = "\"Deepin\""
	assert.Equal(t, info.value, "\"Deepin\"")

	infos = infos.Add("gtk2-test", "test")
	assert.Equal(t, len(infos), 17)

	infos = gtk2FileReader("testdata/xxx")
	infos = infos.Add("gtk2-test", "test")
	assert.Equal(t, len(infos), 1)
	info = infos.Get("gtk2-test")
	assert.Equal(t, info.value, "test")

	err := gtk2FileWriter(infos, "testdata/tmp-gtk2rc")
	defer os.Remove("testdata/tmp-gtk2rc")
	assert.Nil(t, err)
}
