package dxinput

import (
	"pkg.deepin.io/dde/api/dxinput/utils"
)

func SetKeyboardRepeat(enabled bool, delay, interval uint32) error {
	return utils.SetKeyboardRepeat(enabled, delay, interval)
}
