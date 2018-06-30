package bot

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/g4stly/navi/common"
)

type database struct {
	location string
}

func (self *database) startup(location string) error {
	self.location = location
	return nil
}

func (self *database) open() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", self.location)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("sql.Open(): %v", err))
	}
	return db, nil
}

func (self *database) query(commandString string) (*sql.Rows, error) {
	db, err := self.open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare(commandString)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("sql.Open(): %v", err))
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("sql.Open(): %v", err))
	}

	return rows, nil
}

func (self *database) exec(commandString string, args ...interface{}) (sql.Result, error) {
	db, err := self.open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return db.Exec(commandString, args...)
}

func (self *database) LoadSlice(tableName string) ([]string, error) {
	common.Log("loading slice %v", tableName);

	var list []string
	commandString := fmt.Sprintf("SELECT value FROM %v;", tableName)
	rows, err := self.query(commandString)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var listItem string
		rows.Scan(&listItem)
		list = append(list, listItem)
	}
	common.Log("got a slice %v items long", len(list));

	return list, nil
}

func (self *database) SaveSlice(tableName string, slice []string) (error) {
	common.Log("saving slice %v", tableName)
	for index := range slice {
		commandString := fmt.Sprintf("REPLACE INTO %v (value) VALUES ('%v');", tableName, slice[index])
		_, err := self.exec(commandString)
		if err != nil {
			return err
		}
	}

	return nil
}
