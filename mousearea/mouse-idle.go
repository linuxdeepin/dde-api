/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
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

// #cgo amd63 386 CFLAGS: -g -Wall
// #cgo pkg-config: x11 xtst glib-2.0
// #include "mouse-record.h"
import "C"

import (
	"dlib/logger"
	"time"
)

type TimerInfo struct {
	runFlag   bool
	cookie    int32
	timeout   int32
	name      string
	closeChan chan bool
	timer     *time.Timer
}

var (
	cookieTimerMap map[int32]*TimerInfo
	idle           = &IdleTick{}
)

//export startAllTimer
func startAllTimer() {
	if len(cookieTimerMap) <= 0 {
		logger.Println("cookieTimerMap is nil in startTimer")
		return
	}

	for cookie, _ := range cookieTimerMap {
		go startTimer(cookie)
	}
}

//export endAllTimer
func endAllTimer() {
	if len(cookieTimerMap) <= 0 {
		logger.Println("cookieTimerMap is nil in endTimer")
		return
	}

	for cookie, _ := range cookieTimerMap {
		endTimer(cookie, false)
	}
}

func startTimer(cookie int32) {
	info, ok := cookieTimerMap[cookie]
	if !ok {
		logger.Println("Get Timer Info Failed In StartTimer. Cookie:",
			cookie)
		return
	}
	timer := time.NewTimer(time.Second * time.Duration(info.timeout))
	info.timer = timer
	info.runFlag = true

	select {
	case <-info.timer.C:
		logger.Println("Timer Timeout...")
		idle.IdleTimeOut(info.name, info.cookie)

		//info.timer.Stop()
		delete(cookieTimerMap, cookie)
		return
	case <-info.closeChan:
		return
	}
}

func endTimer(cookie int32, deleteFlag bool) {
	info, ok := cookieTimerMap[cookie]
	if !ok {
		logger.Println("Get Timer Info Failed In EndTimer. Cookie:",
			cookie)
		return
	}

	if info.runFlag == false {
		logger.Println("Timer has been end. Cookie:", cookie)
		return
	}
	info.timer.Stop()
	info.closeChan <- true
	if deleteFlag {
		delete(cookieTimerMap, cookie)
	}
}

func newTimerInfo(name string, cookie, timeout int32) *TimerInfo {
	info := &TimerInfo{}

	info.runFlag = false
	info.name = name
	info.cookie = cookie
	info.timeout = timeout
	info.closeChan = make(chan bool)

	return info
}
