/**
 * Copyright (c) 2013 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 Xu FaSheng
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
	"fmt"
	golog "log"
	"os"
	"strings"
	"time"
)

const (
	logfile = "/var/log/deepin.log"
)

var (
	loggerID uint64
	logimpl  *golog.Logger
)

// A Logger represents an active logging object that will provides a
// dbus service to write log message.
type Logger struct {
	names map[uint64]string
}

// NewLogger creates a new Logger.
func NewLogger() *Logger {
	logger := &Logger{}
	logger.names = make(map[uint64]string)
	return logger
}

// GetDBusInfo implement interface of dbus.DBusObject
func (logger *Logger) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.api.Logger",
		"/com/deepin/api/Logger",
		"com.deepin.api.Logger",
	}
}

// NewLogger register a new logger source with name, and return a
// uniquely id which will be used in following operator.
func (logger *Logger) NewLogger(name string) (id uint64, err error) {
	loggerID++
	id = loggerID
	logger.names[id] = name
	logger.doLog(id, "NEW", fmt.Sprintf("id=%d", id))
	return
}

// DeleteLogger unregister a logger source. TODO[remove]
func (logger *Logger) DeleteLogger(id uint64) {
	logger.doLog(id, "DELETE", fmt.Sprintf("id=%d", id))
	delete(logger.names, id)
}

func (logger *Logger) getName(id uint64) (name string) {
	if id == 0 {
		name = "<common>"
		return
	}
	name = logger.names[id]
	if len(name) == 0 {
		name = "<unknown>"
	}
	return
}

func (logger *Logger) doLog(id uint64, level, msg string) {
	now := time.Now()
	date := fmt.Sprintf("%d-%d-%d %d:%d:%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	prefix := fmt.Sprintf("%s %s: [%s] ", date, logger.getName(id), level)
	fmtMsg := prefix + msg
	fmtMsg = strings.Replace(fmtMsg, "\n", "\n"+prefix, -1)
	logimpl.Println(fmtMsg)
	return
}

// Debug will write a log message with 'DEBUG' as prefix.
func (logger *Logger) Debug(id uint64, msg string) {
	logger.doLog(id, "DEBUG", msg)
}

// Info will write a log message with 'INFO' as prefix.
func (logger *Logger) Info(id uint64, msg string) {
	logger.doLog(id, "INFO", msg)
}

// Warning will write a log message with 'WARNING' as prefix.
func (logger *Logger) Warning(id uint64, msg string) {
	logger.doLog(id, "WARNING", msg)
}

// Error will write a log message with 'ERROR' as prefix.
func (logger *Logger) Error(id uint64, msg string) {
	logger.doLog(id, "ERROR", msg)
}

// Fatal will write a log message with 'FATAL' as prefix.
func (logger *Logger) Fatal(id uint64, msg string) {
	logger.doLog(id, "FATAL", msg)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			golog.Fatal(err)
		}
	}()

	// open log file
	logfile, err := os.OpenFile(logfile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer logfile.Close()

	logimpl = golog.New(logfile, "", 0)

	logger := NewLogger()
	err = dbus.InstallOnSystem(logger)
	if err != nil {
		panic(err)
	}
	dbus.DealWithUnhandledMessage()

	select {}
}
