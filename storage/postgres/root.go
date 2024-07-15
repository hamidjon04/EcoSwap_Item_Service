package postgres

import "database/sql"

type ItemRepo struct{
	Db *sql.DB
}

func NewItemRepo(db *sql.DB)*ItemRepo{
	return &ItemRepo{
		Db: db,
	}
}