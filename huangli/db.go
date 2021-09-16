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
package huangli

import (
	"database/sql"
	"fmt"

	"time"

	"encoding/json"

	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// HuangLi huang li info from baidu
type HuangLi struct {
	ID    int64  `json:"id"` // format: ("%s%02s%02s", year, month, day)
	Avoid string `json:"avoid"`
	Suit  string `json:"suit"`
}

// HuangLiList huang li info list
type HuangLiList []*HuangLi

type HolidayStatus int

const (
	HolidayStatusLeave HolidayStatus = iota + 1
	HolidayStatusWork
)

type Holiday struct {
	Date   string        `json:"date"`
	Status HolidayStatus `json:"status"`
}

type HolidayList []*Holiday

func (list HolidayList) Contain(year, month int) bool {
	str := fmt.Sprintf("%d-%d-", year, month)
	for _, info := range list {
		if strings.Contains(info.Date, str) {
			return true
		}
	}
	return false
}

type Festival struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Rest        string `json:"rest"`
	list        string

	Month int `json:"month"`

	Holidays HolidayList `json:"list"`
}

type FestivalList []*Festival

var (
	_db *sql.DB
)

// Init open db and create table
func Init(filename string) error {
	var err error
	_db, err = sql.Open("sqlite3", filename)
	if err != nil {
		return err
	}
	tableStmt := `
CREATE TABLE IF NOT EXISTS huangli (id INTEGER NOT NULL PRIMARY KEY, avoid TEXT, suit TEXT);
`
	_, err = _db.Exec(tableStmt)
	if err != nil {
		return err
	}
	return initFestival()
}

func initFestival() error {
	var year = time.Now().Year()
	var table = fmt.Sprintf("festival_%d", year)
	var tableStmt = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s  ", table)
	tableStmt += `
(id TEXT NOT NULL PRIMARY KEY,month INTEGER NOT NULL,name TEXT,description TEXT, rest TEXT, list TEXT)`
	_, err := _db.Exec(tableStmt)
	return err
}

// Finalize close db
func Finalize() {
	_ = _db.Close()
}

// Create insert to sqlite, if exists, ignore
func (list HuangLiList) Create() error {
	if len(list) == 0 {
		return nil
	}
	tx, err := _db.Begin()
	if err != nil {
		return err
	}

	for _, info := range list {
		tmp, _ := txQueryHuangLi(tx, info.ID)
		if tmp != nil {
			fmt.Println("Has exists:", tmp.ID, tmp.Avoid, tmp.Suit, info.ID, info.Avoid, info.Suit)
			continue
		}
		err = txCreateHuangLi(tx, info)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (list FestivalList) Create(year int) error {
	var table = fmt.Sprintf("festival_%d", year)
	tx, err := _db.Begin()
	if err != nil {
		return err
	}

	for _, info := range list {
		tmp, _ := txQueryFestival(tx, table, info.ID)
		if tmp != nil {
			fmt.Println("Has exists:", tmp.ID, tmp.Name, tmp.Description)
			continue
		}
		err = txCreateFestival(tx, table, info)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (list FestivalList) String() string {
	data, _ := json.Marshal(list)
	return string(data)
}

func (info *Festival) EncodeHolidayList() {
	info.list = ""
	if len(info.Holidays) == 0 {
		return
	}
	data, _ := json.Marshal(info.Holidays)
	info.list = string(data)
}

func (info *Festival) DecodeHolidayList() {
	info.Holidays = HolidayList{}
	if len(info.list) == 0 {
		return
	}
	_ = json.Unmarshal([]byte(info.list), &info.Holidays)
}

// NewHuangLi query by id
func NewHuangLi(id int64) (*HuangLi, error) {
	return txQueryHuangLi(nil, id)
}

// NewHuangLiList query by id list
func NewHuangLiList(idList []int64) (HuangLiList, error) {
	if len(idList) == 0 {
		return nil, nil
	}

	tx, err := _db.Begin()
	if err != nil {
		return nil, err
	}

	var list HuangLiList
	for _, id := range idList {
		info, err := txQueryHuangLi(tx, id)
		if err != nil {
			// TODO(jouyouyun): warning?
			fmt.Println("Failed to query huangli by id:", id, err)
			info = &HuangLi{}
		}
		list = append(list, info)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func NewFestivalList(year, month int) (FestivalList, error) {
	table := fmt.Sprintf("festival_%d", year)
	return txQueryFestivalList(nil, table, month)
}

func txQueryHuangLi(tx *sql.Tx, id int64) (*HuangLi, error) {
	var (
		stmt *sql.Stmt
		err  error
	)

	if tx != nil {
		stmt, err = tx.Prepare("SELECT id, avoid, suit FROM huangli WHERE id = ?")
	} else {
		stmt, err = _db.Prepare("SELECT id, avoid, suit FROM huangli WHERE id = ?")
	}
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	var info HuangLi
	err = stmt.QueryRow(id).Scan(&info.ID, &info.Avoid, &info.Suit)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func txCreateHuangLi(tx *sql.Tx, info *HuangLi) error {
	stmt, err := tx.Prepare("INSERT INTO huangli (id,avoid,suit) VALUES (?,?,?)")
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.Exec(info.ID, info.Avoid, info.Suit)
	return err
}

func txQueryFestival(tx *sql.Tx, table, id string) (*Festival, error) {
	var (
		stmt *sql.Stmt
		err  error
	)
	str := fmt.Sprintf("SELECT id,month,name,description,rest,list FROM %s WHERE id = ?",
		table)
	if tx != nil {
		stmt, err = tx.Prepare(str)
	} else {
		stmt, err = _db.Prepare(str)
	}
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	var info Festival
	err = stmt.QueryRow(id).Scan(&info.ID, &info.Month, &info.Name, &info.Description,
		&info.Rest, &info.list)
	if err != nil {
		return nil, err
	}
	info.DecodeHolidayList()
	return &info, nil
}

func txQueryFestivalList(tx *sql.Tx, table string, month int) (FestivalList, error) {
	var (
		rows *sql.Rows
		err  error
	)
	str := fmt.Sprintf("SELECT id,month,name,description,rest,list FROM %s WHERE month = %d",
		table, month)
	if tx != nil {
		rows, err = tx.Query(str)
	} else {
		rows, err = _db.Query(str)
	}
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var list FestivalList
	for rows.Next() {
		var info Festival
		err := rows.Scan(&info.ID, &info.Month, &info.Name, &info.Description,
			&info.Rest, &info.list)
		if err != nil {
			return nil, err
		}
		info.DecodeHolidayList()
		list = append(list, &info)
	}
	return list, nil
}

func txCreateFestival(tx *sql.Tx, table string, info *Festival) error {
	str := fmt.Sprintf("INSERT INTO %s (id,month,name,description,rest,list) VALUES (?,?,?,?,?,?)",
		table)
	stmt, err := tx.Prepare(str)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	info.EncodeHolidayList()
	_, err = stmt.Exec(info.ID, info.Month, info.Name, info.Description,
		info.Rest, info.list)
	return err
}
