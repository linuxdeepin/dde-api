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
	"syscall"
	"time"
)

const (
	selfID  uint64 = 1
	logfile        = "/var/log/deepin.log"
)

var (
	loggerID = selfID
	logimpl  *golog.Logger
)

// ProcessInfo store the process information which will be
// used to restart if application fataled.
type ProcessInfo struct {
	UID     int               // Real user ID
	Environ map[string]string // Environment variables
	Args    []string          // Command-line arguments, starting with the program name
}

// A Logger represents an active logging object that will provides a
// dbus service to write log message.
type Logger struct {
	names map[uint64]string
}

// NewLogger creates a new Logger object.
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

func (logger *Logger) getName(id uint64) (name string) {
	if id == selfID {
		name = "<logger>"
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

// Debug write a log message with 'DEBUG' as prefix.
func (logger *Logger) Debug(id uint64, msg string) {
	logger.doLog(id, "DEBUG", msg)
}

// Info write a log message with 'INFO' as prefix.
func (logger *Logger) Info(id uint64, msg string) {
	logger.doLog(id, "INFO", msg)
}

// Warning write a log message with 'WARNING' as prefix.
func (logger *Logger) Warning(id uint64, msg string) {
	logger.doLog(id, "WARNING", msg)
}

// Error write a log message with 'ERROR' as prefix.
func (logger *Logger) Error(id uint64, msg string) {
	logger.doLog(id, "ERROR", msg)
}

// Fatal write a log message with 'FATAL' as prefix.
func (logger *Logger) Fatal(id uint64, msg string) {
	logger.doLog(id, "FATAL", msg)
}

// NotifyRestart restart target process, need the following arguments:
// FIXME security issue
//
// - id, logger ID
// - uid, user ID
// - dir, working directory
// - environ, environment variables, in the form "key=value"
// - exefile, program file to execute
// - args, arguments for the program
func (logger *Logger) NotifyRestart(id uint64, uid int32, dir string, environ []string, exefile string, args []string) {
	msg := fmt.Sprintf("uid=%d\nworking directory=%s\nenvironment variables=%v\nprograme=%s\narguments=%v",
		uid, dir, environ, exefile, args)
	logger.doLog(id, "RESTART", msg)
	// TODO
	// logger.doRestart(id, uid, dir, environ, exefile, args)
}

func (logger *Logger) doRestart(id uint64, uid int32, dir string, environ []string, exefile string, args []string) {
	err := syscall.Setreuid(-1, int(uid))
	if err != nil {
		logger.doLog(id, "ERROR", err.Error())
		return
	}
	fmt.Printf("current euid=%d\n", syscall.Geteuid())
	os.StartProcess(exefile, args, &os.ProcAttr{Dir: dir, Env: environ,
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			golog.Fatal(err)
		}
	}()

	// open log file
	f, err := os.OpenFile(logfile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	logimpl = golog.New(f, "", 0)
	logger := NewLogger()
	err = dbus.InstallOnSystem(logger)
	if err != nil {
		golog.Printf("register dbus interface failed: %v\n", err)
		os.Exit(1)
	}
	dbus.DealWithUnhandledMessage()

	if err := dbus.Wait(); err != nil {
		golog.Printf("lost dbus session: %v\n", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
