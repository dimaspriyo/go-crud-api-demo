package repository

import (
	"context"
	"database/sql"
	"go-crud-api-demo/config"

	_ "github.com/go-sql-driver/mysql"
)

func NewConn(ctx context.Context) (*sql.DB, error) {

	var connString string

	read, err := config.ReadYAMLConfig()
	if err != nil {
		return nil, err
	}

	connString = read.DB.Username + ":" + read.DB.Password + "@tcp(" + read.DB.Host + ":" + read.DB.Port + ")/" + read.DB.Name
	conn, err := sql.Open("mysql", connString)
	if err != nil {
		panic(err)
	}

	return conn, nil

}
