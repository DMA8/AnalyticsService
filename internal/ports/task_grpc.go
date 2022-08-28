package ports

import (
	"gitlab.com/g6834/team31/analytics/internal/domain/models"
	"context"
)

// ServerTask Интерфейс grpc Сервера
type ServerTask interface {
	PushTask(ctx context.Context, task models.Task, action, kind int) (models.TaskResponse, error)
	PushMail(ctx context.Context, mail models.Mail) (models.TaskResponse, error)
}
