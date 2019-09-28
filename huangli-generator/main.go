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
	start    = flag.Int("s", 0, "The start year, the min value is 2008")
	end      = flag.Int("e", 0, "The end year, the max year is (now year) + 1")
	festival = flag.Bool("fest", false, "Generate the current year festival db data")
	test     = flag.Bool("t", false, "Test huangli api")
	dbFile   = flag.String("f", "huangli.db", "The huangli data sqlite db file")
)

func main() {
	flag.Parse()
	if *test {
		doTest()
		return
	}

	if !*festival && ((*start == 0 && *end == 0) || *end-*start < 0 || *start < 2008 || *end > (time.Now().Year()+1)) {
		fmt.Printf("Invalid start year and end year: %d - %d\n", *start, *end)
		return
	}

	err := huangli.Init(*dbFile)
	if err != nil {
		panic(err)
	}
	defer huangli.Finalize()

	if *festival {
		genFestivalData()
		return
	}

	// generated db data
	genHuangLiData(*start, *end)
}

func genHuangLiData(s, e int) {
	var list huangli.HuangLiList
	for i := s; i <= e; i++ {
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

	err := list.Create()
	if err != nil {
		fmt.Println("Failed to create db data:", err)
		return
	}
}

func genFestivalData() {
	t := time.Now()
	var list huangli.FestivalList
	for i := 1; i < 13; i++ {
		info, err := newBaiduFestivalByDate(t.Year(), i)
		if err != nil {
			fmt.Println("Failed to get festival data:", err, t.Year(), i)
			return
		}
		list = append(list, info.ToFestival(t.Year(), i)...)
	}
	err := list.Create(t.Year())
	if err != nil {
		fmt.Println("Failed to create festival:", err)
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

	fest, err := newBaiduFestival(data)
	if err != nil {
		fmt.Println("Failed to unmarshal festival:", err)
		return
	}
	fest.Dump()
}

func newBaiduHuangLiByDate(year, month int) (*baiduHuangLi, error) {
	data, err := doGet(makeURL(year, month))
	if err != nil {
		return nil, err
	}
	return newBaiduHuangLi(data)
}

func newBaiduFestivalByDate(year, month int) (*baiduFestival, error) {
	data, err := doGet(makeURL(year, month))
	if err != nil {
		return nil, err
	}
	return newBaiduFestival(data)
}
