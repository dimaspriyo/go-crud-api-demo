package repository

import (
	"database/sql"
)

type Person struct {
	ID      int64 `json:"id"`
	Name    int64 `json:"name"`
	Sex     int64 `json:"sex"`
	Country int64 `json:"country"`
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
