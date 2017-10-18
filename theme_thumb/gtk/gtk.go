package gtk

import (
	"os/exec"
	"strconv"
)

const Version = 0
const cmd = "/usr/lib/deepin-api/gtk-thumbnailer"

func Gen(name string, width, height int, scaleFactor float64, dest string) error {
	var gdkWinScalingFactor float64 = 1.0
	if scaleFactor > 1.7 {
		// 根据 startdde 的逻辑，此种条件下 gtk 窗口放大为 2 倍
		gdkWinScalingFactor = 2.0
	}

	width = int(float64(width) * scaleFactor / gdkWinScalingFactor)
	height = int(float64(height) * scaleFactor / gdkWinScalingFactor)

	var args = []string{
		"-theme", name,
		"-dest", dest,
		"-width", strconv.Itoa(width),
		"-height", strconv.Itoa(height),
		"-force",
	}
	_, err := exec.Command(cmd, args...).CombinedOutput()
	return err
}
