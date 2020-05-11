package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Person struct {
	ID      int64  `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Sex     string `json:"sex,omitempty"`
	Country string `json:"country,omitempty"`
}

func List(db *sql.DB) (res []Person, err error) {

	rows, err := db.Query("select * from persons")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var temp Person

		err := rows.Scan(&temp.ID, &temp.Name, &temp.Sex, &temp.Country)
		if err != nil {
			return nil, err
		}

		res = append(res, temp)
	}

	return res, nil

}

func Insert(ctx context.Context, db *sql.DB, person Person) (res Person, err error) {

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return res, err
	}

	fmt.Println("Starting To Insert Data")

	insert, err := tx.Exec("insert into persons(name, sex, country) values(?,?,?)", person.Name, person.Sex, person.Country)
	if err != nil {
		tx.Rollback()
		return res, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return res, err
	}

	id, err := insert.LastInsertId()
	if err != nil {
		return res, err
	}
	return Detail(db, id)
}

func Detail(db *sql.DB, id int64) (res Person, err error) {

	row := db.QueryRow("select * from persons where id=?", id)

	err = row.Scan(&res.ID, &res.Name, &res.Sex, &res.Country)
	switch err {
	case sql.ErrNoRows:
		return res, errors.New("No Row Found !")
	case nil:
		return res, err
	}

	return res, nil
}
