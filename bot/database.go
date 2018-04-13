package bot

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
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
