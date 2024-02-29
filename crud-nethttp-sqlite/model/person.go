package model

import (
	"database/sql"
	"errors"
	"fmt"
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

func (p *Person) Detail(db *sql.DB) (*Person, error) {
	row := db.QueryRow("SELECT id, name FROM persons WHERE id = ?", p.ID)
	var person Person
	err := row.Scan(&person.ID, &person.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("person with id %d not found", p.ID)
		}
		return nil, err
	}
	return &person, nil
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

func (p *Person) Update(db *sql.DB) error {
	_, err := db.Exec("UPDATE persons SET name = ? WHERE id = ?", p.Name, p.ID)
	if err != nil {
		return err
	}
	return nil
}

func (p *Person) Delete(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM persons WHERE id = ?", p.ID)
	if err != nil {
		return err
	}
	return nil
}
