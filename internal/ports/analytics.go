package ports

import (
	"context"

	"gitlab.com/g6834/team31/analytics/internal/domain/models"
)

type Analytics interface {
	ApprovedTasks(ctx context.Context) (models.Counter, error)
	DeclinedTasks(ctx context.Context) (models.Counter, error)
	SummaryTime(ctx context.Context) ([]models.SummaryTime, error)
	HandleTaskEvent(ctx context.Context, task models.Task, kind models.Kind, action models.Action) error
	HandleMailEvent(ctx context.Context, mail models.Mail) error
}
