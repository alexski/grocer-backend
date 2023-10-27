package model

import (
	"database/sql"
)

type Recipe struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Image    string `json:"image"`
	Filename string `json:"filename"`
	Type     string `json:"type"`
}

func (r *Recipe) GetRecipe(db *sql.DB) error {
	return db.QueryRow("SELECT name, image_filename, filename, type FROM recipes WHERE id=$1", r.ID).Scan(&r.Title, &r.Image, &r.Filename, &r.Type)
}

func (r *Recipe) UpdateRecipe(db *sql.DB) error {
	_, err := db.Exec("UPDATE recipes SET name=$1, image_filename=$2, filename=$3, type=$4 WHERE id=$5", r.Title, r.Image, r.Filename, r.Type, r.ID)
	return err
}

func (r *Recipe) DeleteRecipe(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM recipes WHERE id=$1", r.ID)
	return err
}

func (r *Recipe) CreateRecipe(db *sql.DB) error {
	err := db.QueryRow("INSERT INTO recipes (name, image_filename, filename, type) values ($1, $2, $3, $4) RETURNING id", r.Title, r.Image, r.Filename, r.Type).Scan(&r.ID)
	return err
}

func GetRecipes(db *sql.DB, start, count int) ([]Recipe, error) {
	rows, err := db.Query(
		"SELECT id, name, image_filename, filename, type FROM recipes LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	recipes := []Recipe{}
	for rows.Next() {
		var r Recipe
		if err := rows.Scan(&r.ID, &r.Title, &r.Image, &r.Filename, &r.Type); err != nil {
			return nil, err
		}
		recipes = append(recipes, r)
	}

	return recipes, nil
}
