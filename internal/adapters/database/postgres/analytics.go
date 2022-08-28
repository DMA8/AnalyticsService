package postgres

import (
	"gitlab.com/g6834/team31/analytics/internal/domain/models"
	"bytes"
	"context"
	"encoding/gob"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"time"
)

var namespace = uuid.Nil

func (d *Database) ApprovedTasks(ctx context.Context) (models.Counter, error) {
	var counter int32
	err := d.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM t_tasks WHERE status = 1`).Scan(&counter)
	if err != nil {
		return models.Counter{}, err
	}
	return models.Counter{Count: counter}, nil
}

func (d *Database) DeclinedTasks(ctx context.Context) (models.Counter, error) {
	var counter int32
	// Под несогласованными задачами считаем: несогласованные + еще нет решения
	err := d.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM t_tasks WHERE status != 1`).Scan(&counter)
	if err != nil {
		return models.Counter{}, err
	}
	return models.Counter{Count: counter}, nil
}

func (d *Database) SummaryTime(ctx context.Context) ([]models.SummaryTime, error) {
	query := `SELECT task_id, COALESCE(SUM(diff_time), 0.0) FROM t_task_status GROUP BY task_id`
	rows, err := d.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	var summaryTime []models.SummaryTime
	for rows.Next() {
		var t models.SummaryTime
		err = rows.Scan(&t.TaskId, &t.Duration)
		if err != nil {
			return nil, err
		}
		summaryTime = append(summaryTime, t)
	}
	return summaryTime, nil
}

func (d *Database) NewTask(ctx context.Context, task models.Task) error {
	// Вычисляем хэш от структуры, для идентификации содержимого
	hash := Hash(task)
	keyHash := Hash(task.TaskId)
	err := d.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(
			ctx,
			`INSERT INTO t_tasks (key_hash, task_id, create_ts, owner_login, task_name, description, status, end_ts, email_list, actual, hash)
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ON CONFLICT DO NOTHING`,
			keyHash, task.TaskId, task.CreateTs, task.Email, task.Name, task.Description, task.Status, task.EndTime, task.EmailList, true, hash)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) UpdateTask(ctx context.Context, task models.Task) error {
	// Вычисляем хэш от структуры, для идентификации содержимого
	hash := Hash(task)
	keyHash := Hash(task.TaskId)
	err := d.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(
			ctx,
			`INSERT INTO t_tasks AS t (key_hash, task_id, create_ts, owner_login, task_name, description, status, end_ts, email_list, actual, hash) 
			 values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ON CONFLICT (task_id) DO UPDATE
				SET task_name = EXCLUDED.task_name,
					description = EXCLUDED.description,
					status = EXCLUDED.status,
					end_ts = EXCLUDED.end_ts,
					email_list = EXCLUDED.email_list,
					actual = EXCLUDED.actual,
					hash = EXCLUDED.hash
             WHERE t.hash != EXCLUDED.hash ON CONFLICT DO NOTHING`,
			//WHERE t.task_id = EXCLUDED.task_id AND t.hash != EXCLUDED.hash`,
			keyHash, task.TaskId, task.CreateTs, task.Email, task.Name, task.Description, task.Status, task.EndTime, task.EmailList, task.Actual, hash)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) NewTaskStatus(ctx context.Context, status models.TaskStatus) error {
	var lastTime time.Time
	err := d.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		err := tx.QueryRow(
			ctx,
			`SELECT max_time
				FROM (
					(
						SELECT MAX(create_ts) as max_time
						 FROM t_task_status tts
						 WHERE tts.task_id = $1
						 HAVING MAX(create_ts) NOTNULL
						 LIMIT 1
					)
					UNION ALL
					(
						SELECT create_ts
						FROM t_tasks tt
						WHERE tt.task_id = $2
						LIMIT 1
					)
					LIMIT 1
					) AS source`, status.TaskId, status.TaskId).Scan(&lastTime)
		if err != nil {
			return err
		}
		status.DiffTime = status.CreateTs.Sub(lastTime).Seconds() //не будет ли отрицательным?
		keyHash := Hash(status.TaskId + status.Email)
		_, err = tx.Exec(ctx,
			`INSERT INTO t_task_status (key_hash, task_id, create_ts, email, status, diff_time) values ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING`,
			keyHash, status.TaskId, status.CreateTs, status.Email, status.Status, status.DiffTime)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) NewMail(ctx context.Context, mail models.Mail) error {
	// Вычисляем хэш от структуры, для идентификации содержимого
	hash := Hash(mail)
	err := d.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(
			ctx,
			`INSERT INTO t_mails (hash, header, body, create_ts, email_list)
			values ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`,
			hash, mail.Header, mail.Body, mail.CreateTS, mail.EmailList)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

func Hash(s interface{}) uuid.UUID {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(s)
	if err != nil {
		return uuid.UUID{}
	}
	hash := uuid.NewSHA1(namespace, b.Bytes())
	return hash
}
