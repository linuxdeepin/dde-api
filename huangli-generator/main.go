/*
 * Copyright (C) 2014 ~ 2019 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
package main

import (
	"flag"
	"fmt"
	"time"

	"pkg.deepin.io/dde/api/huangli"
)

var (
	start  = flag.Int("s", 0, "The start year, the min value is 2008")
	end    = flag.Int("e", 0, "The end year, the max year is (now year) + 1")
	test   = flag.Bool("t", false, "Test huangli api")
	dbFile = flag.String("f", "huangli.db", "The huangli data sqlite db file")
)

func main() {
	flag.Parse()
	if *test {
		doTest()
		return
	}

	if (*start == 0 && *end == 0) || *end-*start < 0 || *start < 2008 || *end > (time.Now().Year()+1) {
		fmt.Printf("Invalid start year and end year: %d - %d\n", *start, *end)
		return
	}

	err := huangli.Init(*dbFile)
	if err != nil {
		panic(err)
	}
	defer huangli.Finalize()

	// generated db data
	var list huangli.HuangLiList
	for i := *start; i <= *end; i++ {
		for j := 1; j < 13; j++ {
			if len(list) > 100 {
				err := list.Create()
				if err != nil {
					fmt.Println("Failed to create db data:", err)
					return
				}
				list = huangli.HuangLiList{}
			}
			info, err := newBaiduHuangLiByDate(i, j)
			if err != nil {
				fmt.Println("Failed to generate huangli info:", err)
				return
			}
			list = append(list, info.ToHuangLiList()...)
		}
	}

	if len(list) == 0 {
		return
	}

	err = list.Create()
	if err != nil {
		fmt.Println("Failed to create db data:", err)
		return
	}
}

func doTest() {
	n := time.Now()
	data, err := doGet(makeURL(n.Year(), int(n.Month())))
	if err != nil {
		fmt.Println("Failed to get huangli from api:", err)
		return
	}
	info, err := newBaiduHuangLi(data)
	if err != nil {
		fmt.Println("Failed to unmarshal:", err)
		return
	}
	info.Dump()
}

func newBaiduHuangLiByDate(year, month int) (*baiduHuangLi, error) {
	data, err := doGet(makeURL(year, month))
	if err != nil {
		return nil, err
	}
	return newBaiduHuangLi(data)
}
