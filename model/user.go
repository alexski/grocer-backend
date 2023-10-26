package model

import (
	"database/sql"
	"errors"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (p *User) getUser(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *User) updateUser(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *User) deleteUser(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *User) createUser(db *sql.DB) error {
	return errors.New("Not implemented")
}

func getUsers(db *sql.DB, start, count int) ([]User, error) {
	return nil, errors.New("Not implemented")
}
