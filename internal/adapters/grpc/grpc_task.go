package grpc

import (
	"gitlab.com/g6834/team31/analytics/internal/domain/models"

	"gitlab.com/g6834/team31/analytics/internal/ports"

	"context"
	"net"

	"gitlab.com/g6834/team31/tasks/pkg/grpc_task"

	"gitlab.com/g6834/team31/auth/pkg/logging"

	"google.golang.org/grpc"
)

type TaskServer struct {
	servive ports.Analytics
	grpc_task.UnimplementedGrpcTaskServer
	l *logging.Logger
}

func (t *TaskServer) PushTask(ctx context.Context, in *grpc_task.TaskMessage) (*grpc_task.TaskResponse, error) {
	task := models.Task{
		TaskId:      in.TaskId,
		CreateTs:    in.CreateTs.AsTime(),
		Email:       in.Email,
		Name:        in.Name,
		Description: in.Description,
		Status:      models.Decision(in.Status.Number()),
		EndTime:     in.EndTime.AsTime(),
		EmailList:   in.EmailList,
		Actual:      in.Action.Number() != 2,
	}
	err := t.servive.HandleTaskEvent(ctx, task, models.Kind(in.Kind.Number()), models.Action(in.Status.Number()))
	if err != nil {
		t.l.Debug().Msgf("TaskServer.PushTask err %v", err)
		return nil, err
	}
	t.l.Debug().Msgf("TaskServer.PushTask  task event handeled successfully %+v", in)
	return &grpc_task.TaskResponse{}, err
}

func (t *TaskServer) PushMail(ctx context.Context, in *grpc_task.Mail) (*grpc_task.TaskResponse, error) {
	err := t.servive.HandleMailEvent(ctx, models.Mail{
		Header:    in.Header,
		Body:      in.Body,
		CreateTS:  in.CreateTs.AsTime(),
		EmailList: in.EmailList,
	})
	if err != nil {
		t.l.Debug().Msgf("TaskServer.PushMail err %v", err)
		return nil, err
	}
	t.l.Debug().Msgf("TaskServer.PushMail  mail added successfully %+v", in)
	return &grpc_task.TaskResponse{}, err
}

func LaunchGRPCServer(port string, service ports.Analytics, l *logging.Logger) chan error {
	chanErr := make(chan error)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		l.Fatal().Err(err)
	}
	s := grpc.NewServer()
	server := &TaskServer{servive: service, l: l}
	grpc_task.RegisterGrpcTaskServer(s, server)
	go func() {
		chanErr <- s.Serve(lis)
	}()
	return chanErr
}
