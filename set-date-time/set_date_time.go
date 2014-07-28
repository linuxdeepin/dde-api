package main

import (
	"os"
	"os/exec"
	"pkg.linuxdeepin.com/lib"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/log"
	"strconv"
	"time"
)

const (
	_SET_DATE_TIME_DEST = "com.deepin.api.SetDateTime"
	_SET_DATE_TIME_PATH = "/com/deepin/api/SetDateTime"
	_SET_DATA_TIME_IFC  = "com.deepin.api.SetDateTime"
)

var (
	logger = log.NewLogger("dde-api/set-date-time")
)

type SetDateTime struct {
	GenLocaleStatus func(bool, string)
}

func (sdt *SetDateTime) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_SET_DATE_TIME_DEST,
		_SET_DATE_TIME_PATH,
		_SET_DATA_TIME_IFC,
	}
}

func (sdt *SetDateTime) SetCurrentDate(d string) bool {
	/* Date String Format: 2013-11-17 */
	if CountCharInString('-', d) != 2 {
		logger.Info("date string format error")
		return false
	}

	sysTime := time.Now()
	sysTmp := &sysTime
	_, tStr := GetDateTimeAny(sysTmp)
	cmd := exec.Command("/bin/date", "--set", d)
	_, err := cmd.Output()
	if err != nil {
		logger.Info("Set Date error:", err)
		return false
	}
	sdt.SetCurrentTime(tStr)
	return true
}

func (sdt *SetDateTime) SetCurrentTime(t string) bool {
	/* Time String Format: 12:23:09 */
	if CountCharInString(':', t) != 2 {
		logger.Info("time string format error")
		return false
	}

	cmd := exec.Command("/bin/date", "+%T", "-s", t)
	_, err := cmd.Output()
	if err != nil {
		logger.Info("Set time error:", err)
		return false
	}
	return true
}

func (sdt *SetDateTime) GetTimezone() (string, bool) {
	return getTimezone()
}

func (sdt *SetDateTime) SetTimezone(tz string) bool {
	return setTimezone(tz)
}

func NewSetDateTime() *SetDateTime {
	sdt := SetDateTime{}

	return &sdt
}

func CountCharInString(ch byte, str string) int {
	if l := len(str); l <= 0 {
		return 0
	}

	cnt := 0
	for i, _ := range str {
		if str[i] == ch {
			cnt++
		}
	}

	return cnt
}

func GetDateTimeAny(t *time.Time) (dStr, tStr string) {
	dStr += strconv.FormatInt(int64(t.Year()), 10) + "-" + strconv.FormatInt(int64(t.Month()), 10) + "-" + strconv.FormatInt(int64(t.Day()), 10)
	tStr += strconv.FormatInt(int64(t.Hour()), 10) + ":" + strconv.FormatInt(int64(t.Minute()), 10) + ":" + strconv.FormatInt(int64(t.Second()), 10)

	return dStr, tStr
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main() {
	defer logger.EndTracing()

	if !lib.UniqueOnSystem(_SET_DATE_TIME_DEST) {
		logger.Warning("There already has an SetDateTime daemon running.")
		return
	}

	// configure logger
	logger.SetRestartCommand("/usr/lib/deepin-api/set-date-time", "--debug")
	if stringInSlice("-d", os.Args) || stringInSlice("--debug", os.Args) {
		logger.SetLogLevel(log.LEVEL_DEBUG)
	}

	sdt := NewSetDateTime()
	err := dbus.InstallOnSystem(sdt)
	if err != nil {
		panic(err)
	}
	dbus.DealWithUnhandledMessage()

	dbus.SetAutoDestroyHandler(time.Second*1, func() bool {
		if genLcStart {
			if genLcEnd {
				genLcEnd = false
				genLcStart = false
				return true
			}
		} else {
			return true
		}

		return false
	})
	if err = dbus.Wait(); err != nil {
		logger.Error("lost dbus session:", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
