package main

// #cgo pkg-config: x11 xcursor xfixes gtk+-3.0
// #include <stdlib.h>
// #include "cursor.h"
import "C"

import (
	"fmt"
	"os"
	"pkg.deepin.io/dde/api/themes"
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"strings"
	"sync"
	"time"
	"unsafe"
)

type Manager struct {
	locker  sync.Mutex
	running bool
}

const (
	dbusDest = "com.deepin.api.CursorHelper"
	dbusPath = "/com/deepin/api/CursorHelper"
	dbusIFC  = "com.deepin.api.CursorHelper"
)

var logger = log.NewLogger("api/cursor-helper")

func NewManager() *Manager {
	var m = new(Manager)
	m.running = false
	return m
}

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func (m *Manager) Set(name string) {
	m.locker.Lock()
	m.running = true
	go func() {
		defer m.locker.Unlock()
		setTheme(name)
		m.running = false
	}()
}

func main() {
	var name string
	if len(os.Args) == 2 {
		tmp := strings.ToLower(os.Args[1])
		if tmp == "-h" || tmp == "--help" {
			fmt.Println("Usage: cursor-theme-helper <Curosr theme>")
			return
		}
		name = os.Args[1]
	}

	if !lib.UniqueOnSession(dbusDest) {
		logger.Warning("There already has a cursor helper running...")
		return
	}

	if C.init_gtk() == -1 {
		logger.Warning("Init gtk or x11 thread environment failed")
		return
	}

	var m = NewManager()
	err := dbus.InstallOnSession(m)
	if err != nil {
		logger.Error("Install session dbus failed:", err)
		return
	}
	dbus.DealWithUnhandledMessage()

	if len(name) != 0 {
		setTheme(name)
		return
	}

	dbus.SetAutoDestroyHandler(time.Second*5, func() bool {
		if m.running {
			return false
		}
		return true
	})

	err = dbus.Wait()
	if err != nil {
		logger.Error("Lost cursor helper session:", err)
		os.Exit(-1)
	}
}

func setTheme(name string) {
	if name == themes.GetCursorTheme() {
		return
	}

	go doSetTheme(name)
	C.handle_gtk_cursor_changed()
}

func doSetTheme(name string) {
	err := themes.SetCursorTheme(name)
	if err != nil {
		logger.Warning("Set failed:", err)
	}

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.set_gtk_cursor(cName)
	C.set_qt_cursor(cName)
}
