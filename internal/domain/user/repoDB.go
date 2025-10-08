package user

import (
	"database/sql"
)

type UserDBRepo struct {
	db *sql.DB
}

func NewUserDBRepo(db *sql.DB) *UserDBRepo {
	return &UserDBRepo{db: db}
}

func (u *UserDBRepo) CreateUser(user User) error {
	query := `
			INSERT INTO users (username,password)
			VALUES ($1, $2)
			RETURNING id
		`
	return u.db.QueryRow(query, user.Username, user.Password).Scan(&user.ID)
}
