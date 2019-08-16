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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"pkg.deepin.io/dde/api/huangli"
)

const (
	apiURL     = "https://sp0.baidu.com/8aQDcjqpAAV3otqbppnN2DJv/api.php"
	resourceID = 6018 // 黄历资源 ID
	apiCharset = "utf8"
)

type baiduHuangLi struct {
	Data []struct {
		Almanac []struct {
			Date  string `json:"date"`
			Avoid string `json:"avoid"`
			Suit  string `json:"suit"`
		} `json:"almanac"`
	} `json:"data"`
}

func (info *baiduHuangLi) ToHuangLiList() huangli.HuangLiList {
	var list huangli.HuangLiList
	for _, almanac := range info.Data {
		for _, value := range almanac.Almanac {
			id, err := convertDateToID(value.Date)
			if err != nil {
				fmt.Println("Failed to convert date to id:", err)
				continue
			}
			list = append(list, &huangli.HuangLi{
				ID:    id,
				Avoid: value.Avoid,
				Suit:  value.Suit,
			})
		}
	}
	return list
}

func (info *baiduHuangLi) Dump() {
	fmt.Println("Baidu huangli:")
	for _, almanac := range info.Data {
		for _, value := range almanac.Almanac {
			fmt.Printf("\tDate: %q, \tavoid: %q, \tsuit: %q\n",
				value.Date, value.Avoid, value.Suit)
		}
	}
	fmt.Println("Baidu huangli dump done")
}

func convertDateToID(date string) (int64, error) {
	list := strings.SplitN(date, "-", 3)
	if len(list) != 3 {
		return 0, fmt.Errorf("invalid baidu huangli date: %s", date)
	}
	return strconv.ParseInt(fmt.Sprintf("%s%02s%02s", list[0], list[1], list[2]), 10, 64)
}

func newBaiduHuangLi(data []byte) (*baiduHuangLi, error) {
	var info baiduHuangLi
	err := json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func doGet(uri string) ([]byte, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	if resp.Body == nil {
		return nil, fmt.Errorf("no data return")
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", string(data))
	}
	return data, nil
}

func makeURL(year, month int) string {
	var params = make(url.Values)
	params["resource_id"] = []string{fmt.Sprint(resourceID)}
	params["ie"] = []string{apiCharset}
	params["oe"] = []string{apiCharset}
	params["query"] = []string{fmt.Sprintf("%d年%d月", year, month)}

	return fmt.Sprintf("%s?%s", apiURL, params.Encode())
}
