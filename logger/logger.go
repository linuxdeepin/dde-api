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
	"bufio"
	"crypto/rand"
	"dlib"
	"dlib/dbus"
	"fmt"
	"io/ioutil"
	golog "log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	selfName    = "<logger>"
	selfID      = "0000000"
	unknownName = "<unknown>"
	logfile     = "/var/log/deepin.log"
)

var logimpl *golog.Logger

// A Logger represents an active logging object that will provides a
// dbus service to write log message.
type Logger struct{}

// NewLogger creates a new Logger object.
func NewLogger() *Logger {
	logger := &Logger{}
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
func (logger *Logger) NewLogger(name string) (id string, err error) {
	id = randString(len(selfID))
	logger.doLog(id, name, "NEW", fmt.Sprintf("id=%s", id))
	return
}

func randString(n int) string {
	const alphanum = "0123456789abcdef"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func (logger *Logger) doLog(id, name string, level, msg string) {
	now := time.Now()
	date := fmt.Sprintf("%02d-%02d-%02d %02d:%02d:%02d.%03d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), int(float64(now.Nanosecond())/float64(999999999)*1000))
	prefix := fmt.Sprintf("%s %s %s: [%s] ", id, date, name, level)
	fmtMsg := prefix + msg
	fmtMsg = strings.Replace(fmtMsg, "\n", "\n"+prefix, -1) // format multi-lines message
	logimpl.Println(fmtMsg)
	return
}

// Debug write a log message with 'DEBUG' as prefix.
func (logger *Logger) Debug(id, name, msg string) {
	logger.doLog(id, name, "DEBUG", msg)
}

// Info write a log message with 'INFO' as prefix.
func (logger *Logger) Info(id, name, msg string) {
	logger.doLog(id, name, "INFO", msg)
}

// Warning write a log message with 'WARNING' as prefix.
func (logger *Logger) Warning(id, name, msg string) {
	logger.doLog(id, name, "WARNING", msg)
}

// Error write a log message with 'ERROR' as prefix.
func (logger *Logger) Error(id, name, msg string) {
	logger.doLog(id, name, "ERROR", msg)
}

// Fatal write a log message with 'FATAL' as prefix.
func (logger *Logger) Fatal(id, name, msg string) {
	logger.doLog(id, name, "FATAL", msg)
}

// GetLog return all log messages that wrote by target ID.
func (logger *Logger) GetLog(id string) (msg string) {
	// get all deepin log files
	logfiles, err := filepath.Glob(logfile + "*")
	if err != nil {
		msg = "get log files failed"
		golog.Println(msg)
		return
	}

	// open log file in order
	sort.Sort(sort.Reverse(sort.StringSlice(logfiles)))
	for _, f := range logfiles {
		m, err := logger.doGetLog(id, f)
		if err != nil {
			golog.Println(err)
			continue
		}
		if len(m) != 0 {
			msg = msg + m
		}
	}
	return
}

func (logger *Logger) doGetLog(id string, file string) (msg string, err error) {
	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	s := bufio.NewScanner(strings.NewReader(string(fileContent)))
	s.Split(bufio.ScanLines)
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, id) {
			msg = msg + line + "\n"
		}
	}

	return
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			golog.Fatal(err)
		}
	}()

	if !dlib.UniqueOnSystem("com.deepin.api.Logger") {
		golog.Println("There already has an Logger daemon running.")
		return
	}

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

	dbus.SetAutoDestroyHandler(5*time.Second, nil)

	if err := dbus.Wait(); err != nil {
		golog.Printf("lost dbus session: %v\n", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
