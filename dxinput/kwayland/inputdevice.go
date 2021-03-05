package kwayland

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/godbus/dbus"
	kwin "github.com/linuxdeepin/go-dbus-factory/org.kde.kwin"
	. "pkg.deepin.io/dde/api/dxinput/common"
)

const (
	SysNamePrefix = "event"

	eventPathPrefix = "/org/kde/KWin/InputDevice/"
)

var (
	_conn *dbus.Conn

	errUnsupported = fmt.Errorf("unsupported this operation")
)

func getSessionBus() *dbus.Conn {
	if _conn == nil {
		conn, err := dbus.SessionBus()
		if err != nil {
			fmt.Println("Failed to connect session bus:", err)
			panic(err)
		}
		_conn = conn
	}

	return _conn
}

func ListDevice() (DeviceInfos, error) {
	var manager = kwin.NewInputDeviceManager(getSessionBus())
	var infos DeviceInfos

	sysNames, err := manager.DevicesSysNames().Get(0)
	if err != nil {
		return nil, err
	}
	for _, sysName := range sysNames {
		info, err := NewDeviceInfo(sysName)
		if err != nil {
			fmt.Println("Failed to new input device:", err)
			continue
		}
		if info != nil {
			infos = append(infos, info)
		}
	}

	return infos, nil
}

func Enable(sysName string, enabled bool) error {
	dev, err := newInputDeviceObj(sysName)
	if err != nil {
		return err
	}
	return dev.Enabled().Set(0, enabled)
}

func CanEnabled(sysName string) bool {
	dev, err := newInputDeviceObj(sysName)
	if err != nil {
		return false
	}
	enabled, _ := dev.Enabled().Get(0)
	return enabled
}

func EnableLeftHanded(sysName string, enabled bool) error {
	dev, err := newInputDeviceObj(sysName)
	if err != nil {
		return err
	}
	if supp, _ := dev.SupportsLeftHanded().Get(0); !supp {
		return errUnsupported
	}

	return dev.LeftHanded().Set(0, enabled)
}

func CanLeftHanded(sysName string) bool {
	dev, err := newInputDeviceObj(sysName)
	if err != nil {
		return false
	}
	enabled, _ := dev.LeftHanded().Get(0)
	return enabled
}

func EnableMiddleEmulation(sysName string, enabled bool) error {
	dev, err := newInputDeviceObj(sysName)
	if err != nil {
		return err
	}
	if supp, _ := dev.SupportsMiddleEmulation().Get(0); !supp {
		return errUnsupported
	}

	return dev.MiddleEmulation().Set(0, enabled)
}

func CanMiddleButtonEmulation(sysName string) bool {
	dev, err := newInputDeviceObj(sysName)
	if err != nil {
		return false
	}
	enabled, _ := dev.MiddleEmulation().Get(0)
	return enabled
}

func SetPointerAccel(sysName string, v float64) error {
	dev, err := newInputDeviceObj(sysName)
	if err != nil {
		return err
	}

	if supp, _ := dev.SupportsPointerAcceleration().Get(0); !supp {
		return errUnsupported
	}

	return dev.PointerAcceleration().Set(0, v)
}

func GetPointerAccel(sysName string) (float64, error) {
	dev, err := newInputDeviceObj(sysName)
	if err != nil {
		return 0, err
	}

	return dev.PointerAcceleration().Get(0)
}

func EnableTapToClick(sysName string, enabled bool) error {
	dev, err := newInputDeviceObj(sysName)
	if err != nil {
		return err
	}

	return dev.TapToClick().Set(0, enabled)
}

func CanTapToClick(sysName string) bool {
	dev, err := newInputDeviceObj(sysName)
	if err != nil {
		return false
	}
	enabled, _ := dev.TapToClick().Get(0)
	return enabled
}

func EnableDisableWhileTyping(sysName string, enabled bool) error {
	dev, err := newInputDeviceObj(sysName)
	if err != nil {
		return err
	}
	if supp, _ := dev.SupportsDisableWhileTyping().Get(0); !supp {
		return errUnsupported
	}

	return dev.DisableWhileTyping().Set(0, enabled)
}

func CanDisableWhileTyping(sysName string) bool {
	dev, err := newInputDeviceObj(sysName)
	if err != nil {
		return false
	}
	enabled, _ := dev.DisableWhileTyping().Get(0)
	return enabled
}

func EnableAdaptiveAccelProfile(sysName string, enabled bool) error {
	dev, err := newInputDeviceObj(sysName)
	if err != nil {
		return err
	}
	if supp, _ := dev.SupportsPointerAccelerationProfileAdaptive().Get(0); !supp {
		return errUnsupported
	}

	return dev.PointerAccelerationProfileAdaptive().Set(0, enabled)
}

func CanAdaptiveAccelProfile(sysName string) bool {
	dev, err := newInputDeviceObj(sysName)
	if err != nil {
		return false
	}
	enabled, _ := dev.PointerAccelerationProfileAdaptive().Get(0)
	return enabled
}

func SetScrollButton(sysName string, v uint32) error {
	dev, err := newInputDeviceObj(sysName)
	if err != nil {
		return err
	}

	return dev.ScrollButton().Set(0, v)
}

func GetScrollButton(sysName string) (uint32, error) {
	dev, err := newInputDeviceObj(sysName)
	if err != nil {
		return 0, err
	}

	return dev.ScrollButton().Get(0)
}

func EnableScrollTwoFinger(sysName string, enabled bool) error {
	dev, err := newInputDeviceObj(sysName)
	if err != nil {
		return err
	}
	if supp, _ := dev.SupportsScrollTwoFinger().Get(0); !supp {
		return errUnsupported
	}

	return dev.ScrollTwoFinger().Set(0, enabled)
}

func CanScrollTwoFinger(sysName string) bool {
	dev, err := newInputDeviceObj(sysName)
	if err != nil {
		return false
	}
	enabled, _ := dev.ScrollTwoFinger().Get(0)
	return enabled
}

func EnableScrollEdge(sysName string, enabled bool) error {
	dev, err := newInputDeviceObj(sysName)
	if err != nil {
		return err
	}
	if supp, _ := dev.SupportsScrollEdge().Get(0); !supp {
		return errUnsupported
	}

	return dev.ScrollEdge().Set(0, enabled)
}

func CanScrollEdge(sysName string) bool {
	dev, err := newInputDeviceObj(sysName)
	if err != nil {
		return false
	}
	enabled, _ := dev.ScrollEdge().Get(0)
	return enabled
}

func EnableNaturalScroll(sysName string, enabled bool) error {
	dev, err := newInputDeviceObj(sysName)
	if err != nil {
		return err
	}
	if supp, _ := dev.SupportsNaturalScroll().Get(0); !supp {
		return errUnsupported
	}

	return dev.NaturalScroll().Set(0, enabled)
}

func CanNaturalScroll(sysName string) bool {
	dev, err := newInputDeviceObj(sysName)
	if err != nil {
		return false
	}
	enabled, _ := dev.NaturalScroll().Get(0)
	return enabled
}

func EnableLmrTapButtonMap(sysName string, enabled bool) error {
	dev, err := newInputDeviceObj(sysName)
	if err != nil {
		return err
	}
	if supp, _ := dev.SupportsLmrTapButtonMap().Get(0); !supp {
		return errUnsupported
	}

	return dev.LmrTapButtonMap().Set(0, enabled)
}

func CanLmrTapButtonMap(sysName string) bool {
	dev, err := newInputDeviceObj(sysName)
	if err != nil {
		return false
	}
	enabled, _ := dev.LmrTapButtonMap().Get(0)
	return enabled
}

func NewDeviceInfo(sysName string) (*DeviceInfo, error) {
	dev, err := kwin.NewInputDevice(getSessionBus(), dbus.ObjectPath(eventPathPrefix+sysName))
	if err != nil {
		fmt.Println("Failed to new input device:", err)
		return nil, err
	}
	//dumpInputDevice(dev)

	var info DeviceInfo
	kbd, _ := dev.Keyboard().Get(0)
	numKbd, _ := dev.AlphaNumericKeyboard().Get(0)
	if kbd && numKbd {
		if isMouseDevice(dev) {
			info.Type = DevTypeMouse
		} else {
			info.Type = DevTypeKeyboard
		}
		goto fill
	}

	if tpad, _ := dev.Touchpad().Get(0); tpad {
		info.Type = DevTypeTouchpad
		goto fill
	}

	if isMouseDevice(dev) {
		info.Type = DevTypeMouse
		goto fill
	}

	return nil, nil

fill:
	id, _ := strconv.Atoi(strings.Split(sysName, SysNamePrefix)[1])
	info.Id = int32(id)
	info.Name, _ = dev.Name().Get(0)
	info.Enabled, _ = dev.Enabled().Get(0)

	return &info, nil
}

func isMouseDevice(dev kwin.InputDevice) bool {
	fcount, _ := dev.TapFingerCount().Get(0)
	suppLeftHanded, _ := dev.SupportsLeftHanded().Get(0)
	suppBtns, _ := dev.SupportedButtons().Get(0)
	sbtn, _ := dev.ScrollButton().Get(0)
	return (fcount == 0) && suppLeftHanded && (suppBtns > 0) && (sbtn != 0)
}

func newInputDeviceObj(sysName string) (kwin.InputDevice, error) {
	return kwin.NewInputDevice(getSessionBus(), dbus.ObjectPath(eventPathPrefix+sysName))
}
