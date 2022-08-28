package analytics

import (
	"context"
	"fmt"

	"gitlab.com/g6834/team31/analytics/internal/domain/models"
	"gitlab.com/g6834/team31/analytics/internal/ports"
	"gitlab.com/g6834/team31/auth/pkg/logging"
)

type Service struct {
	db ports.DbInterface
	l  *logging.Logger
}

func New(db ports.DbInterface, l *logging.Logger) *Service {
	return &Service{
		db: db,
		l:  l,
	}
}

func (s *Service) ApprovedTasks(ctx context.Context) (models.Counter, error) {
	counter, err := s.db.ApprovedTasks(ctx)
	if err != nil {
		return models.Counter{}, err
	}
	return counter, nil
}

func (s *Service) DeclinedTasks(ctx context.Context) (models.Counter, error) {
	counter, err := s.db.DeclinedTasks(ctx)
	if err != nil {
		return models.Counter{}, err
	}
	return counter, nil
}

func (s *Service) SummaryTime(ctx context.Context) ([]models.SummaryTime, error) {
	result, err := s.db.SummaryTime(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) HandleTaskEvent(ctx context.Context, task models.Task, kind models.Kind, action models.Action) error {
	s.l.Debug().Msgf("service.HandleTaskEvent kind %d, action %d task %+v", kind, action, task)
	if kind == models.TaskKind {
		switch action {
		case models.CreateAction:
			return s.db.NewTask(ctx, task)
		case models.UpdateAction:
			return s.db.UpdateTask(ctx, task)
		default:
			return fmt.Errorf("service.TaskEventHandler there are CreateAction and UpdateAction only %w", models.ErrUnexpectedEventAction)
		}
	} else if kind == models.StatusKind {
		taskStatus := models.TaskStatus{
			TaskId:   task.TaskId,
			CreateTs: task.CreateTs,
			Email:    task.Email,
			Status:   task.Status,
		}
		return s.db.NewTaskStatus(ctx, taskStatus)
	} else {
		return fmt.Errorf("service.TaskEventHandler there are TaskKind and StatusKind only %w", models.ErrUnexpectedEventKind)
	}
}

func (s *Service) HandleMailEvent(ctx context.Context, mail models.Mail) error {
	s.l.Debug().Msgf("service.HandleMailEvent mail: %+v", mail)
	return s.db.NewMail(ctx, mail)
}

