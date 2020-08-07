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
	"strconv"
	"time"

	"github.com/jinzhu/gorm"

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
	const baseDBFile = "./huangliBase.db"
	flag.Parse()
	if *test {
		doTest()
		return
	}

	if !*festival && ((*start == 0 && *end == 0) || *end-*start < 0 || *start < 2008 || *end > (time.Now().Year()+20)) {
		fmt.Printf("Invalid start year and end year: %d - %d\n", *start, *end)
		return
	}

	err := huangli.Init(*dbFile)
	if err != nil {
		panic(err)
	}
	defer huangli.Finalize()

	db, err := gorm.Open("sqlite3", baseDBFile)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = db.Close()
	}()

	if *festival {
		genFestivalData()
		return
	}

	// generated db data
	genHuangLiData(db, *start, *end)
}

type Huangli struct {
	ID int
	Y  int    // 年
	M  int    // 月
	D  int    // 日
	Yi string // 宜
	Ji string // 忌
}

func genHuangLiData(db *gorm.DB, start, end int) {
	var baseHuangliList []*Huangli
	var genHuangliList huangli.HuangLiList
	err := db.Where("Y >= ? AND Y <= ?", start, end).Find(&baseHuangliList).Error
	if err != nil {
		fmt.Println("Failed to get db data:", err)
	}
	genHuangliList = make(huangli.HuangLiList, 0, 100)
	for _, item := range baseHuangliList {
		if len(genHuangliList) >= 100 {
			err := genHuangliList.Create()
			if err != nil {
				fmt.Println("Failed to create db data:", err)
				return
			}
			genHuangliList = genHuangliList[0:0:100]
		}
		temp := &huangli.HuangLi{}
		temp.Avoid = item.Ji
		temp.Suit = item.Yi
		temp.ID, _ = strconv.ParseInt(fmt.Sprintf("%d%02d%02d", item.Y, item.M, item.D), 10, 64)
		genHuangliList = append(genHuangliList, temp)
	}
	err = genHuangliList.Create()
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

func newBaiduFestivalByDate(year, month int) (*baiduFestival, error) {
	data, err := doGet(makeURL(year, month))
	if err != nil {
		return nil, err
	}
	return newBaiduFestival(data)
}
