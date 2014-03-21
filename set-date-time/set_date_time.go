package main

import (
        "dlib/dbus"
        dlogger "dlib/logger"
        "net"
        "os"
        "os/exec"
        "strconv"
        "time"
)

const (
        _NTP_HOST           = "0.pool.ntp.org"
        _SET_DATE_TIME_DEST = "com.deepin.api.SetDateTime"
        _SET_DATE_TIME_PATH = "/com/deepin/api/SetDateTime"
        _SET_DATA_TIME_IFC  = "com.deepin.api.SetDateTime"
)

var (
        logger   = dlogger.NewLogger("dde-api/set-date-time")
        quitChan chan bool
)

type SetDateTime struct {
        ntpRunFlag bool
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

func (sdt *SetDateTime) SyncNtpTime() bool {
        for i := 0; i < 10; i++ {
                t, err := GetNtpNow()
                if err == nil && t != nil {
                        dStr, tStr := GetDateTimeAny(t)
                        logger.Infof("Data: %s, Time: %s", dStr, tStr)
                        sdt.SetCurrentDate(dStr)
                        sdt.SetCurrentTime(tStr)
                        return true
                } else {
                        logger.Info(err.Error())
                        //return false
                }
        }

        return false
}

func (sdt *SetDateTime) SetNtpUsing(using bool) bool {
        if using {
                if sdt.ntpRunFlag {
                        sdt.SyncNtpTime()
                        logger.Info("Ntp is running....")
                        return true
                }

                sdt.ntpRunFlag = true
                go SetNtpThread(sdt)
        } else {
                if sdt.ntpRunFlag {
                        logger.Info("Ntp will quit....")
                        quitChan <- true
                }

                sdt.ntpRunFlag = false
        }
        return true
}

func SetNtpThread(sdt *SetDateTime) {
        for {
                sdt.SyncNtpTime()
                timer := time.NewTimer(time.Minute * 1)
                select {
                case <-timer.C:
                case <-quitChan:
                        return
                }
        }
}

func NewSetDateTime() *SetDateTime {
        sdt := SetDateTime{}
        sdt.ntpRunFlag = false

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

func GetNtpNow() (*time.Time, error) {
        raddr, err := net.ResolveUDPAddr("udp", _NTP_HOST+":123")
        if err != nil {
                return nil, err
        }

        data := make([]byte, 48)
        data[0] = 3<<3 | 3

        con, err := net.DialUDP("udp", nil, raddr)
        if err != nil {
                return nil, err
        }

        defer con.Close()

        _, err = con.Write(data)
        if err != nil {
                return nil, err
        }

        con.SetDeadline(time.Now().Add(5 * time.Second))

        _, err = con.Read(data)
        if err != nil {
                return nil, err
        }

        var sec, frac uint64
        sec = uint64(data[43]) | uint64(data[42])<<8 | uint64(data[41])<<16 |
                uint64(data[40])<<24
        frac = uint64(data[47]) | uint64(data[46])<<8 | uint64(data[45])<<16 |
                uint64(data[44])<<24

        nsec := sec * 1e9
        nsec += (frac * 1e9) >> 32

        t := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC).
                Add(time.Duration(nsec)).Local()

        return &t, nil
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
        defer func() {
                if err := recover(); err != nil {
                        logger.Error("recover err:", err)
                }
        }()

        // configure logger
        logger.SetRestartCommand("/usr/lib/deepin-api/set-date-time", "--debug")
        if stringInSlice("-d", os.Args) || stringInSlice("--debug", os.Args) {
                logger.SetLogLevel(dlogger.LEVEL_DEBUG)
        }

        quitChan = make(chan bool)
        sdt := NewSetDateTime()
        err := dbus.InstallOnSystem(sdt)
        if err != nil {
                panic(err)
        }
        dbus.DealWithUnhandledMessage()
        //select {}
        if err = dbus.Wait(); err != nil {
                logger.Error("lost dbus session:", err)
                os.Exit(1)
        } else {
                os.Exit(0)
        }
}
