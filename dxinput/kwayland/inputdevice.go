package kwayland

import (
	"fmt"
	"strconv"
	"strings"

	kwin "github.com/linuxdeepin/go-dbus-factory/org.kde.kwin"
	. "pkg.deepin.io/dde/api/dxinput/common"
	dbus "pkg.deepin.io/lib/dbus1"
)

const (
	SysNamePrefix = "event"

	eventPathPrefix = "/org/kde/KWin/InputDevice/"
)

var (
	_conn *dbus.Conn

	errUnsupported = fmt.Errorf("unsupported this operation")
)

func init() {
	conn, err := dbus.SessionBus()
	if err != nil {
		fmt.Println("Failed to connect session bus:", err)
		panic(err)
	}
	_conn = conn
}

func ListDevice() (DeviceInfos, error) {
	var manager = kwin.NewInputDeviceManager(_conn)
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
	dev, err := kwin.NewInputDevice(_conn, dbus.ObjectPath(eventPathPrefix+sysName))
	if err != nil {
		fmt.Println("Failed to new input device:", err)
		return nil, err
	}
	dumpInputDevice(dev)

	var info DeviceInfo
	kbd, _ := dev.Keyboard().Get(0)
	numKbd, _ := dev.AlphaNumericKeyboard().Get(0)
	if kbd && numKbd {
		info.Type = DevTypeKeyboard
		goto fill
	}

	if tpad, _ := dev.Touchpad().Get(0); tpad {
		info.Type = DevTypeTouchpad
		goto fill
	}

	{
		fcount, _ := dev.TapFingerCount().Get(0)
		suppLeftHanded, _ := dev.SupportsLeftHanded().Get(0)
		suppBtns, _ := dev.SupportedButtons().Get(0)
		sbtn, _ := dev.ScrollButton().Get(0)
		if (fcount == 0) && suppLeftHanded && (suppBtns > 0) && (sbtn != 0) {
			info.Type = DevTypeMouse
			goto fill
		}
	}

	return nil, nil

fill:
	id, _ := strconv.Atoi(strings.Split(sysName, SysNamePrefix)[1])
	info.Id = int32(id)
	info.Name, _ = dev.Name().Get(0)
	info.Enabled, _ = dev.Enabled().Get(0)

	return &info, nil
}

func newInputDeviceObj(sysName string) (*kwin.InputDevice, error) {
	return kwin.NewInputDevice(_conn, dbus.ObjectPath(eventPathPrefix+sysName))
}

func dumpInputDevice(dev *kwin.InputDevice) {
	name, _ := dev.Name().Get(0)
	fmt.Printf("%s\n", name)
	enabled, _ := dev.Enabled().Get(0)
	fmt.Println("\tEnabled:", enabled)
	sysName, _ := dev.SysName().Get(0)
	fmt.Println("\tSysName:", sysName)
	vendor, _ := dev.Vendor().Get(0)
	fmt.Println("\tVendor:", vendor)
	product, _ := dev.Product().Get(0)
	fmt.Println("\tProduct:", product)
	kbd, _ := dev.Keyboard().Get(0)
	fmt.Println("\tIs Keybaord:", kbd)
	numKbd, _ := dev.AlphaNumericKeyboard().Get(0)
	fmt.Println("\tIs Number Keyboard:", numKbd)
	tpad, _ := dev.Touchpad().Get(0)
	fmt.Println("\tIs Touchpad:", tpad)
	touch, _ := dev.Touch().Get(0)
	fmt.Println("\tCan Touch:", touch)
	tabTool, _ := dev.TabletTool().Get(0)
	fmt.Println("\tIs Tablet Tool:", tabTool)
	tms, _ := dev.TabletModeSwitch().Get(0)
	fmt.Println("\tTablet Mode Switch:", tms)
	cma, _ := dev.ClickMethodAreas().Get(0)
	fmt.Println("\tClick Method Areas:", cma)
	dcma, _ := dev.DefaultClickMethodAreas().Get(0)
	fmt.Println("\tDefault Click Method Areas:", dcma)
	cmcf, _ := dev.ClickMethodClickfinger().Get(0)
	fmt.Println("\tClick Method Click Finger:", cmcf)
	dcmcf, _ := dev.DefaultClickMethodClickfinger().Get(0)
	fmt.Println("\tDefault Click Method Click Finger:", dcmcf)
	paccel, _ := dev.PointerAcceleration().Get(0)
	fmt.Println("\tPointer Acceleration:", paccel)
	dpaccel, _ := dev.DefaultPointerAcceleration().Get(0)
	fmt.Println("\tDefault Pointer Acceleration:", dpaccel)
	paccelpa, _ := dev.PointerAccelerationProfileAdaptive().Get(0)
	fmt.Println("\tPointer Acceleration Profile Adaptive:", paccelpa)
	dpaccelpa, _ := dev.DefaultPointerAccelerationProfileAdaptive().Get(0)
	fmt.Println("\tDefault Pointer Acceleration Profile Adaptive:", dpaccelpa)
	paccelpf, _ := dev.PointerAccelerationProfileFlat().Get(0)
	fmt.Println("\tPointer Acceleration Profile Flat:", paccelpf)
	dpaccelpf, _ := dev.DefaultPointerAccelerationProfileFlat().Get(0)
	fmt.Println("\tDefault Pointer Acceleration Profile Flat:", dpaccelpf)
	sbtn, _ := dev.ScrollButton().Get(0)
	fmt.Println("\tScroll Button:", sbtn)
	dsbtn, _ := dev.DefaultScrollButton().Get(0)
	fmt.Println("\tDefault Scroll Button:", dsbtn)
	dwtyping, _ := dev.DisableWhileTyping().Get(0)
	fmt.Println("\tDisable While Typing:", dwtyping)
	dwtyinged, _ := dev.DisableWhileTypingEnabledByDefault().Get(0)
	fmt.Println("\tDisable While Typing Enabled Default:", dwtyinged)
	gestureSupp, _ := dev.GestureSupport().Get(0)
	fmt.Println("\tGesture Support:", gestureSupp)
	leftHanded, _ := dev.LeftHanded().Get(0)
	fmt.Println("\tLeft Handed:", leftHanded)
	leftHandedD, _ := dev.LeftHandedEnabledByDefault().Get(0)
	fmt.Println("\tLeft Handed Default:", leftHandedD)
	lidS, _ := dev.LidSwitch().Get(0)
	fmt.Println("\tLid Switch:", lidS)
	lmrTapBtnMap, _ := dev.LmrTapButtonMap().Get(0)
	fmt.Println("\tLmr Tap Button Map:", lmrTapBtnMap)
	lmrTapBtnMapD, _ := dev.LmrTapButtonMapEnabledByDefault().Get(0)
	fmt.Println("\tLmr Tap Button Map Enable Default:", lmrTapBtnMapD)
	midEmu, _ := dev.MiddleEmulation().Get(0)
	fmt.Println("\tMiddle Emulation:", midEmu)
	midEmuD, _ := dev.MiddleEmulationEnabledByDefault().Get(0)
	fmt.Println("\tMiddle Emulation Default:", midEmuD)
	natureS, _ := dev.NaturalScroll().Get(0)
	fmt.Println("\tNature Scroll:", natureS)
	natureSD, _ := dev.NaturalScrollEnabledByDefault().Get(0)
	fmt.Println("\tNature Scroll Default:", natureSD)
	outName, _ := dev.OutputName().Get(0)
	fmt.Println("\tOutput Name:", outName)
	sedge, _ := dev.ScrollEdge().Get(0)
	fmt.Println("\tScroll Edge:", sedge)
	sedgeD, _ := dev.ScrollEdgeEnabledByDefault().Get(0)
	fmt.Println("\tScroll Edge Default:", sedgeD)
	sobtnd, _ := dev.ScrollOnButtonDown().Get(0)
	fmt.Println("\tScroll On Button Down:", sobtnd)
	sobtndD, _ := dev.ScrollOnButtonDownEnabledByDefault().Get(0)
	fmt.Println("\tScroll On Button Down Default:", sobtndD)
	stfinger, _ := dev.ScrollTwoFinger().Get(0)
	fmt.Println("\tScroll Two Finger:", stfinger)
	stfingerD, _ := dev.ScrollTwoFingerEnabledByDefault().Get(0)
	fmt.Println("\tScroll Two Finger Default:", stfingerD)
	suppBtns, _ := dev.SupportedButtons().Get(0)
	fmt.Println("\tSupported Buttons:", suppBtns)
	suppCaliMatrix, _ := dev.SupportsCalibrationMatrix().Get(0)
	fmt.Println("\tSupported Calibration Matrix:", suppCaliMatrix)
	suppClickMA, _ := dev.SupportsClickMethodAreas().Get(0)
	fmt.Println("\tSupported Click Method Areas:", suppClickMA)
	suppClickMCf, _ := dev.SupportsClickMethodClickfinger().Get(0)
	fmt.Println("\tSupported Click Method Clickfinger:", suppClickMCf)
	suppDisWTyping, _ := dev.SupportsDisableWhileTyping().Get(0)
	fmt.Println("\tSupport Disable While Typing:", suppDisWTyping)
	suppLeftHanded, _ := dev.SupportsLeftHanded().Get(0)
	fmt.Println("\tSupport Left Handed:", suppLeftHanded)
	suppLmrTapBtn, _ := dev.SupportsLmrTapButtonMap().Get(0)
	fmt.Println("\tSupport Lmr Tap Button Map:", suppLmrTapBtn)
	suppMidEmu, _ := dev.SupportsMiddleEmulation().Get(0)
	fmt.Println("\tSupport Middle Emulation:", suppMidEmu)
	suppNaturalS, _ := dev.SupportsNaturalScroll().Get(0)
	fmt.Println("\tSupport Natural Scroll:", suppNaturalS)
	suppPAccel, _ := dev.SupportsPointerAcceleration().Get(0)
	fmt.Println("\tSupport Pointer Acceleration:", suppPAccel)
	suppPAccelP, _ := dev.SupportsPointerAccelerationProfileAdaptive().Get(0)
	fmt.Println("\tSupport Pointer Acceleration Profile Adaptive:", suppPAccelP)
	suppPAccelPF, _ := dev.SupportsPointerAccelerationProfileFlat().Get(0)
	fmt.Println("\tSupport Pointer Acceleration Profile Flat:", suppPAccelPF)
	suppSEdge, _ := dev.SupportsScrollEdge().Get(0)
	fmt.Println("\tSupport Scroll Edge:", suppSEdge)
	suppSOBtn, _ := dev.SupportsScrollOnButtonDown().Get(0)
	fmt.Println("\tSupport Scroll On Button Down:", suppSOBtn)
	suppSTFinger, _ := dev.SupportsScrollTwoFinger().Get(0)
	fmt.Println("\tSupport Scroll Two Finger:", suppSTFinger)
	tad, _ := dev.TapAndDrag().Get(0)
	fmt.Println("\tTap And Drag:", tad)
	tadD, _ := dev.TapAndDragEnabledByDefault().Get(0)
	fmt.Println("\tTap And Drag Default:", tadD)
	tdl, _ := dev.TapDragLock().Get(0)
	fmt.Println("\tTap Drag Lock:", tdl)
	tdlD, _ := dev.TapDragLockEnabledByDefault().Get(0)
	fmt.Println("\tTap Drag Lock Default:", tdlD)
	ttc, _ := dev.TapToClick().Get(0)
	fmt.Println("\tTap To Click:", ttc)
	ttcD, _ := dev.TapToClickEnabledByDefault().Get(0)
	fmt.Println("\tTap To Click Default:", ttcD)
	tapFinger, _ := dev.TapFingerCount().Get(0)
	fmt.Println("\tTap Finger Count:", tapFinger)
}
