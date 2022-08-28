package int_tests

import (
	"gitlab.com/g6834/team31/analytics/internal/domain/models"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func Task(ctx context.Context, Pool *pgxpool.Pool, taskName string) (*models.Task, error) {
	var t models.Task
	row, err := Pool.Query(ctx, `SELECT key_hash, owner_login, description, create_ts FROM t_tasks WHERE task_name = $1`, taskName)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		err = row.Scan(&t.TaskId, &t.Name, &t.Description, &t.CreateTs)
		if err != nil {
			return nil, err
		}
	}
	return &t, nil
}

func Status(ctx context.Context, Pool *pgxpool.Pool, taskID string, email string) (*models.TaskStatus, error) {
	var t models.TaskStatus

	row, err := Pool.Query(ctx, `SELECT task_id, email, status FROM t_task_status WHERE task_id = $1 AND email = $2`, taskID, email)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		err = row.Scan(&t.TaskId, &t.Email, &t.Status)
		if err != nil {
			return nil, err
		}
	}
	return &t, nil
}

func Mail(ctx context.Context, Pool *pgxpool.Pool, mailHeader string) (*models.Mail, error) {
	var t models.Mail

	row, err := Pool.Query(ctx, `SELECT header, body FROM t_mails WHERE header = $1`, mailHeader)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		err = row.Scan(&t.Header, &t.Body)
		if err != nil {
			return nil, err
		}
	}
	return &t, nil
}


func DeleteTask(ctx context.Context, Pool *pgxpool.Pool, taskName string) error {
	_, err := Pool.Exec(ctx, `DELETE FROM t_tasks WHERE owner_login = $1`, taskName)
	return err
}
func DeleteTaskAll(ctx context.Context, Pool *pgxpool.Pool) error {
	_, err := Pool.Exec(ctx, `TRUNCATE t_tasks;`)
	return err
}
func DeleteTaskStatus(ctx context.Context, Pool *pgxpool.Pool) error {
	_, err := Pool.Exec(ctx, `TRUNCATE t_task_status;`)
	return err
}

func DeleteMails(ctx context.Context, Pool *pgxpool.Pool) error {
	_, err := Pool.Exec(ctx, `TRUNCATE t_mails;`)
	return err
}
