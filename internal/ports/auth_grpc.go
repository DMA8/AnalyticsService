package ports

import (
	"gitlab.com/g6834/team31/analytics/internal/domain/models"
	"context"
)

// ClientAuth Интерфейс grpc Клиента
type ClientAuth interface {
	Validate(ctx context.Context, in models.JWTTokens) (models.ValidateResponse, error)
}
