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
	return nil
}

// Finalize close db
func Finalize() {
	_ = _db.Close()
}

// Create insert to sqlite, if exists, ignore
func (list HuangLiList) Create() error {
	tx, err := _db.Begin()
	if err != nil {
		return err
	}

	for _, info := range list {
		tmp, _ := txQuery(tx, info.ID)
		if tmp != nil {
			fmt.Println("Has exists:", tmp.ID, tmp.Avoid, tmp.Suit, info.ID, info.Avoid, info.Suit)
			continue
		}
		err = txCreate(tx, info)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// NewHuangLi query by id
func NewHuangLi(id int64) (*HuangLi, error) {
	return txQuery(nil, id)
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
		info, err := txQuery(tx, id)
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

func txQuery(tx *sql.Tx, id int64) (*HuangLi, error) {
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
	defer stmt.Close()

	var info HuangLi
	err = stmt.QueryRow(id).Scan(&info.ID, &info.Avoid, &info.Suit)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func txCreate(tx *sql.Tx, info *HuangLi) error {
	stmt, err := tx.Prepare("INSERT INTO huangli (id,avoid,suit) VALUES (?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(info.ID, info.Avoid, info.Suit)
	return err
}
