package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-crud-api-demo/config"

	_ "github.com/go-sql-driver/mysql"
)

func NewConn(ctx context.Context) (*sql.DB, error) {

	var connString string

	read, err := config.ReadConfig()
	if err != nil {
		return nil, err
	}

	connString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", read.DB.Username, read.DB.Password, read.DB.Host, read.DB.Port, read.DB.Name)
	conn, err := sql.Open("mysql", connString)
	if err != nil {
		panic(err)
	}

	return conn, nil

}
