/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

import (
        "dlib/dbus"
        liblogger "dlib/logger"
        "os"
)

var logger = liblogger.NewLogger("dde-api/sound")

func main() {
        defer func() {
                if err := recover(); err != nil {
                        logger.Fatal("%v", err)
                }
        }()

        s := &Sound{}
        err := dbus.InstallOnSession(s)
        if err != nil {
                logger.Error("register dbus interface failed: %v", err)
                os.Exit(1)
        }

        dbus.DealWithUnhandledMessage()
        if err := dbus.Wait(); err != nil {
                logger.Error("lost dbus session: %v", err)
                os.Exit(1)
        } else {
                os.Exit(0)
        }
}
