package model

import (
	"database/sql"
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}

func (u *User) GetUser(db *sql.DB) error {
	return db.QueryRow("SELECT username, password_hash FROM users WHERE id=$1", u.ID).Scan(&u.Username, &u.PasswordHash)
}

func (u *User) UpdateUser(db *sql.DB) error {
	_, err := db.Exec("UPDATE users SET username=$1, password_hash=$2 WHERE id=$3", u.Username, u.PasswordHash, u.ID)
	return err
}

func (u *User) DeleteUser(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM users WHERE id=$1", u.ID)
	return err
}

func (u *User) CreateUser(db *sql.DB) error {
	err := db.QueryRow("INSERT INTO users (username, password_hash) values ($1, $2) RETURNING id", u.Username, u.PasswordHash).Scan(&u.ID)
	return err
}

func GetUsers(db *sql.DB, start, count int) ([]User, error) {
	rows, err := db.Query(
		"SELECT id, username, password_hash FROM users LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.PasswordHash); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
