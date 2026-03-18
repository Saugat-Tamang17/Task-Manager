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

func (r *postgresTaskRepo) Create(ctx context.Context, task *model.Task) error {
	err := r.db.QueryRowContext(ctx, "INSERT into tasks(title,description,status,created_at) VALUES ($1,$2,$3,$4) RETURNING id)", task.Title, task.Description, task.Status, task.CreatedAt).Scan(&task.Id)
	if err != nil {
		return fmt.Errorf("taskRepo.Create: %w", err)
	}
	return nil
}

func (r *postgresTaskRepo) Update(ctx context.Context, id string, task *model.Task) error {
	_, err := r.db.ExecContext(ctx, "UPDATE tasks SET title=$1, description=$2, status=$3, created_at=$4 WHERE id=$5", task.Title, task.Description, task.Status, task.CreatedAt, id)
	if err != nil {
		return fmt.Errorf("TaskRepo.Update =%w", err)
	}
	return nil
}

func (r *postgresTaskRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM tasks WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("taskRepo.Delete: %w", err)
	}
	return nil
}

// now that all of the task are over , we moving to the authenticating functions //

func (r *postgresAuthRepo) Register(ctx context.Context, user *model.AuthUser) error {
	err := r.db.QueryRowContext(ctx, "INSERT INTO users_jwt(username,password) VALUES($1,$2) RETURNING id", user.Username, user.Password).Scan(&user.Id)
	if err != nil {
		return fmt.Errorf("AuthRepo.Register : %w", err)
	}
	return nil
}
