// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package themes

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	dutils "github.com/linuxdeepin/go-lib/utils"
)

func TestGtk3Prop(t *testing.T) {
	kfile, err := dutils.NewKeyFileFromFile("testdata/settings.ini")
	assert.Nil(t, err)
	defer kfile.Free()

	assert.Equal(t, isGtk3PropEqual(gtk3KeyTheme, "Paper",
		kfile), true)
	assert.Equal(t, isGtk3PropEqual("gtk-menu-images", "1",
		kfile), true)
	assert.Equal(t, isGtk3PropEqual("gtk-modules", "gail:atk-bridge",
		kfile), true)
	assert.Equal(t, isGtk3PropEqual("test-list", "1;2;3;",
		kfile), true)

	err = setGtk3Prop("test-gtk3", "test", "testdata/tmp-gtk3")
	defer os.Remove("testdata/tmp-gtk3")
	assert.Nil(t, err)
}
