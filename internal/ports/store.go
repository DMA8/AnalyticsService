package ports

import (
	"gitlab.com/g6834/team31/analytics/internal/domain/models"
	"context"
)

// DbInterface . Абстракция от СУБД
type DbInterface interface {
	ApprovedTasks(ctx context.Context) (models.Counter, error)
	DeclinedTasks(ctx context.Context) (models.Counter, error)
	SummaryTime(ctx context.Context) ([]models.SummaryTime, error)
	NewTask(ctx context.Context, task models.Task) error
	UpdateTask(ctx context.Context, task models.Task) error
	NewTaskStatus(ctx context.Context, status models.TaskStatus) error
	NewMail(ctx context.Context, mail models.Mail) error
}
