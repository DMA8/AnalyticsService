package mq

import (
	"context"

	"gitlab.com/g6834/team31/analytics/internal/config"
	"gitlab.com/g6834/team31/analytics/internal/domain/models"
	"gitlab.com/g6834/team31/analytics/internal/ports"
	"gitlab.com/g6834/team31/auth/pkg/logging"
	"gitlab.com/g6834/team31/tasks/pkg/grpc_task"
	"gitlab.com/g6834/team31/tasks/pkg/mq"
	"gitlab.com/g6834/team31/tasks/pkg/mq/types"
	"google.golang.org/protobuf/proto"
)

type MQClient struct {
	service           ports.Analytics
	consumerMailTopic mq.Consumer
	consumerTaskTopic mq.Consumer
	logger            *logging.Logger
}

func New(ConsumerMailTopic, ConsumerTaskTopic mq.Consumer, service ports.Analytics, logger *logging.Logger) *MQClient {
	return &MQClient{
		consumerMailTopic: ConsumerMailTopic,
		consumerTaskTopic: ConsumerTaskTopic,
		logger:            logger,
		service:           service,
	}
}

func BuildMQClient(cfg config.Kafka, service ports.Analytics, logger *logging.Logger) *MQClient {
	consumerMailTopic, err := mq.NewConsumer([]string{cfg.URL}, cfg.MailTopic, cfg.GroupID)
	if err != nil {
		logger.Fatal().Err(err).Msg("couldn't init mail topic consumer")
	}
	consumerTaslTopic, err := mq.NewConsumer([]string{cfg.URL}, cfg.TaskTopic, cfg.GroupID)
	if err != nil {
		logger.Fatal().Err(err).Msg("couldn't init task topic consumer")
	}
	return &MQClient{
		consumerMailTopic: consumerMailTopic,
		consumerTaskTopic: consumerTaslTopic,
		service:           service,
		logger:            logger,
	}
}

func (m *MQClient) RunMQ(ctx context.Context) chan error {
	errChan := make(chan error)
	go m.ConsumeMailTopic(ctx, errChan)
	go m.ConsumeTaskTopic(ctx, errChan)
	return errChan
}

func (m *MQClient) ConsumeMailTopic(ctx context.Context, errChan chan error) {
	for {
		_, err := m.consumerMailTopic.ReadAndCommit(ctx, m.operateMailEvent)
		if err != nil {
			errChan <- err
		}
	}
}

func (m *MQClient) ConsumeTaskTopic(ctx context.Context, errChan chan error) {
	for {
		_, err := m.consumerTaskTopic.ReadAndCommit(ctx, m.operateTaskEvent)
		if err != nil {
			errChan <- err
		}
	}
}

func (m *MQClient) operateTaskEvent(ctx context.Context, msg types.Message) error {
	var taskInp grpc_task.TaskMessage
	err := proto.Unmarshal(msg.Value, &taskInp)
	if err != nil {
		m.logger.Warn().Err(err).Msg("operateTaskEvent error while unmarshaling kafka event")
		return err
	}
	task := models.Task{
		TaskId:      taskInp.TaskId,
		CreateTs:    taskInp.CreateTs.AsTime(),
		Email:       taskInp.Email,
		Name:        taskInp.Name,
		Description: taskInp.Description,
		Status:      models.Decision(taskInp.Status.Number()),
		EndTime:     taskInp.EndTime.AsTime(),
		EmailList:   taskInp.EmailList,
		Actual:      taskInp.Action.Number() != 2,
	}
	return m.service.HandleTaskEvent(ctx, task, models.Kind(taskInp.Kind.Number()), models.Action(taskInp.Status.Number()))
}

func (m *MQClient) operateMailEvent(ctx context.Context, msg types.Message) error {
	var mailInp grpc_task.Mail
	err := proto.Unmarshal(msg.Value, &mailInp)
	if err != nil {
		m.logger.Warn().Err(err).Msg("operateMailEvent error while unmarshaling kafka event")
		return err
	}
	mail := models.Mail{
		Header:    mailInp.Header,
		Body:      mailInp.Body,
		CreateTS:  mailInp.CreateTs.AsTime(),
		EmailList: mailInp.EmailList,
	}
	return m.service.HandleMailEvent(ctx, mail)
}
