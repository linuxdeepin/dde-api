// Copyright (c) 2015 Deepin Ltd. All rights reserved.
// Use of this source is govered by General Public License that can be found
// in the LICENSE file.
package main

import (
	"os"
	"time"

	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger(DBusName)

func main() {
	logger.BeginTracing()
	defer logger.EndTracing()

	if !lib.UniqueOnSession(DBusName) {
		logger.Warning("Validator daemon is already running.")
		return
	}

	validator := &Validator{}
	err := dbus.InstallOnSession(validator)
	if err != nil {
		logger.Errorf("Failed to register dbus interface: %v", err)
		os.Exit(1)
	}
	dbus.DealWithUnhandledMessage()

	dbus.SetAutoDestroyHandler(30*time.Second, nil)
	if err := dbus.Wait(); err != nil {
		logger.Errorf("lost dbus session: %v", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
