package user

import (
	"database/sql"
)

type DBRepo struct {
	db *sql.DB
}

func NewUserDBRepo(db *sql.DB) *DBRepo {
	return &DBRepo{db: db}
}

func (u *DBRepo) CreateUser(user User) error {
	query := `
			INSERT INTO users (username,password)
			VALUES ($1, $2)
			RETURNING id
		`
	return u.db.QueryRow(query, user.Username, user.Password).Scan(&user.ID)
}
