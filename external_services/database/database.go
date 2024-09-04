package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
)

func NewDatabase() (*sqlx.DB, error) {
	mysqlURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT"), os.Getenv("MYSQL_DATABASE"))
	db, err := sqlx.Connect("mysql", mysqlURI)
	if err != nil {
		return nil, err
	}

	return db, nil
}
