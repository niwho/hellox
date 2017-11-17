package common

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DBClient struct {
	*gorm.DB
	Name     string
	Server   string
	User     string
	Password string
	filename string
}

func NewMysqlDBClient(name, server, user, passwd string) *DBClient {
	c := &DBClient{
		Name:     name,
		Server:   server,
		User:     user,
		Password: passwd,
	}
	if err := c.initMysqlDb(); err != nil {
		fmt.Errorf("db init error=%v", err)
	}
	return c
}

func NewSqliteDBClient(filename string) *DBClient {
	c := &DBClient{
		filename: filename,
	}
	if err := c.initSqliteDb(); err != nil {
		fmt.Errorf("db init error=%v", err)
	}
	return c
}

func (db *DBClient) initMysqlDb() error {
	var err error
	db.DB, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=100ms&readTimeout=50ms&writeTimeout=50ms", db.User, db.Password, db.Server, db.Name))
	return err
}

func (db *DBClient) initSqliteDb() error {
	var err error
	db.DB, err = gorm.Open("sqlite3", db.filename)
	return err
}
