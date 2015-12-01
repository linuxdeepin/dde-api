package main

import (
	"encoding/json"
	"fmt"
	"os"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/pinyin"
	"time"
)

const (
	dbusDest = "com.deepin.api.Pinyin"
	dbusPath = "/com/deepin/api/Pinyin"
	dbusIFC  = dbusDest
)

type Manager struct{}

func (*Manager) Query(hans string) []string {
	return queryPinyin(hans)
}

// Querylist query pinyin for hans list, return a json data.
func (*Manager) QueryList(hansList []string) string {
	var data = make(map[string][]string)
	for _, hans := range hansList {
		data[hans] = queryPinyin(hans)
	}

	content, _ := json.Marshal(data)
	return string(content)
}

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func main() {
	if len(os.Args) == 2 {
		if os.Args[1] == "--help" || os.Args[1] == "-h" {
			usage()
			return
		}

		fmt.Println(queryPinyin(os.Args[1]))
		return
	}

	err := dbus.InstallOnSession(new(Manager))
	if err != nil {
		fmt.Println("Install dbus failed:", err)
		return
	}

	dbus.DealWithUnhandledMessage()
	dbus.SetAutoDestroyHandler(time.Second*5, nil)
	err = dbus.Wait()
	if err != nil {
		fmt.Println("Lost dbus:", err)
		os.Exit(-1)
	}

	os.Exit(0)
}

func usage() {
	fmt.Println("Usage: hans2pinyin <hans>")
	fmt.Println("Example:")
	fmt.Println("\thans2pinyin 重")
	fmt.Println("\t[zhong chong]")
}

func queryPinyin(hans string) []string {
	return pinyin.HansToPinyin(hans)
}
