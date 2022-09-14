// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/linuxdeepin/dde-api/themes"
)

func Test_setTheme(t *testing.T) {
	cursorTheme := themes.GetCursorTheme()
	tests := []struct {
		Input    string
		Expected error
	}{
		{
			cursorTheme,
			nil,
		},
		{
			"fake1Theme",
			fmt.Errorf("invalid theme '%s'", "fake1Theme"),
		},
		{
			"fake2Theme",
			fmt.Errorf("invalid theme '%s'", "fake2Theme"),
		},
	}
	for i, test := range tests {
		t.Run("Test_setTheme"+strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, test.Expected, setTheme(test.Input))
		})
	}
}
