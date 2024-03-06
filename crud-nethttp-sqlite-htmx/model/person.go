package model

import (
	"database/sql"
)

type Person struct {
	ID   int    `sql:"id" json:"id"`
	Name string `sql:"name" json:"name"`
}

func (p *Person) List(db *sql.DB) ([]Person, error) {
	rows, err := db.Query("SELECT id, name FROM persons")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var persons []Person
	for rows.Next() {
		var p Person
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			return nil, err
		}
		persons = append(persons, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return persons, nil
}

func (p *Person) Insert(db *sql.DB) (*Person, error) {
	result, err := db.Exec("INSERT INTO persons (name) VALUES (?)", p.Name)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	p.ID = int(id)
	return p, nil
}

func (p *Person) Delete(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM persons WHERE id = ?", p.ID)
	if err != nil {
		return err
	}
	return nil
}
