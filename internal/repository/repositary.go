package repository

import (
	"database/sql"
)

type postgresTaskRepo struct {
	db *sql.DB
}

type postgresAuthRepo struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) TaskRepositary {
	return &postgresTaskRepo{db: db}
}

func NewAuthRepo(db *sql.DB) AuthRepositary {
	return &postgresAuthRepo{db: db}
}
