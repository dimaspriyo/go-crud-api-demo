package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
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

	exec, err := tx.Exec("INSERT into persons(name, sex, country) values(?,?,?)", person.Name, person.Sex, person.Country)
	if err != nil {
		tx.Rollback()
		return res, err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return res, err
	}

	id, err := exec.LastInsertId()
	if err != nil {
		return res, err
	}
	return Detail(db, id)
}

func Detail(db *sql.DB, id int64) (res Person, err error) {

	row := db.QueryRow("SELECT * from persons WHERE id=?", id)

	err = row.Scan(&res.ID, &res.Name, &res.Sex, &res.Country)
	switch err {
	case sql.ErrNoRows:
		return res, errors.New("No Row Found !")
	case nil:
		return res, err
	}

	return res, nil
}

func Update(ctx context.Context, db *sql.DB, person Person, personId int64) (res Person, err error) {

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return res, err
	}

	_, err = tx.Exec("UPDATE persons set name=?, sex=?, country=? WHERE id=?", person.Name, person.Sex, person.Country, personId)
	if err != nil {
		tx.Rollback()
		return res, err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return res, err
	}

	return Detail(db, personId)

}

func Delete(ctx context.Context, db *sql.DB, id int64) (err error) {

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE from persons where id=?", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
