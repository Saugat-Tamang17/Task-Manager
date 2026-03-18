package repository

import (
	"context"
	"database/sql"
	"fmt"
	"taskmanager/internal/model"
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

// this part will contain the task methods //

func (r *postgresTaskRepo) GetAll(ctx context.Context) ([]model.Task, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id,title,description,status,created_at FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("taskRepoGetAll :%w", err)
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.Id, &t.Title, &t.Description, &t.Status, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("taskRepo.GetAll :%w", err)

		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
