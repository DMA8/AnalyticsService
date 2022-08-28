//go:build integration
// +build integration

package int_tests

import (
	"analytics/config"
	"analytics/internal/adapters/database/postgres"
	"analytics/internal/domain/models"
	"analytics/pkg/logging"
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose"
	"github.com/stretchr/testify/suite"
)

type integraTestSuite struct {
	suite.Suite

	cfg *config.Config
	db  *postgres.Database
	ctx context.Context
}

func TestIntegraTestSuite(t *testing.T) {
	suite.Run(t, &integraTestSuite{})
}

func testConfig() *config.Config {
	return &config.Config{
		App: config.App{
			Name:    "test",
			Version: "1",
		},
		Auth: config.Auth{
			HOST: "localhost",
			PORT: ":9000",
		},
		HTTP: config.HTTP{
			Port:       ":3003",
			ApiVersion: "/analytics/v1",
		},
		Log: config.Log{
			Level: "trace",
		},
		PG: config.PG{
			URL: "postgresql://test:test@127.0.0.1:5432/test",
		},
		Queue: config.Queue{
			HOST: "localhost",
			PORT: ":9001",
		},
	}
}

func (s *integraTestSuite) SetupSuite() {
	cfg := testConfig()
	s.cfg = cfg
	ctx, _ := context.WithCancel(context.Background())
	l := logging.New(cfg.Log.Level)
	pgConnStr := os.Getenv("PG_CONNSTR")
	if pgConnStr != "" {
		l.Info().Msg("found address to postgres in env")
		cfg.PG.URL = pgConnStr
	}

	l.Info().Msg("connecting to db...")
	dbCtx, _ := context.WithTimeout(ctx, time.Second*10)
	db, err := postgres.New(dbCtx, cfg.PG)
	if err != nil {
		l.Fatal().Msg(err.Error())
	}

	// Migrations
	l.Info().Msg("starting migrations...")
	pgxCfg, err := pgx.ParseConfig(cfg.PG.URL)
	if err != nil {
		l.Fatal().Msgf("pgx.ParseConfig: %+v", err)
	}
	pgxCfg.PreferSimpleProtocol = true

	gooseDb := stdlib.OpenDB(*pgxCfg)
	err = goose.SetDialect("postgres")
	if err != nil {
		l.Fatal().Msg(err.Error())
	}
	if err := goose.Up(gooseDb, "migrations"); err != nil {
		l.Fatal().Msg(err.Error())
	}
	s.ctx = ctx
	s.db = db
}

func (s *integraTestSuite) TestNewTask() {
	defer DeleteTaskAll(s.ctx, s.db.Pool)
	testTask := models.Task{
		Email:       "test1@email.com",
		Name:        "test1",
		Description: "test1descr",
		Actual:      true,
		EmailList:   []string{"email1", "email2", "email3"},
		Status:      models.NoDecision,
	}
	err := s.db.NewTask(s.ctx, testTask)
	s.NoError(err)
	result, err := Task(s.ctx, s.db.Pool, testTask.Name)
	s.NoError(err)
	s.Equal(testTask.CreateTs, result.CreateTs)
}

func (s *integraTestSuite) TestUpdateTask() {
	defer DeleteTaskAll(s.ctx, s.db.Pool)
	testTask := models.Task{
		Email:       "test2@email.com",
		Name:        "test2",
		Description: "test2descr",
		Actual:      true,
		EmailList:   []string{"2email1", "2email2", "2email3"},
		Status:      models.NoDecision,
	}
	err := s.db.NewTask(s.ctx, testTask)
	s.NoError(err)
	result, err := Task(s.ctx, s.db.Pool, testTask.Name)
	s.NoError(err)
	s.Equal(testTask.CreateTs, result.CreateTs)
	testUpdatedInput := testTask
	testUpdatedInput.Description = "updatedDescr"
	err = s.db.UpdateTask(s.ctx, testUpdatedInput)
	s.NoError(err)
	resTask, err := Task(s.ctx, s.db.Pool, testTask.Name)
	s.NoError(err)
	s.Equal(resTask.Description, testUpdatedInput.Description)
}

func (s *integraTestSuite) TestNewTaskStatus() {
	defer DeleteTaskStatus(s.ctx, s.db.Pool)
	testTask := models.Task{
		TaskId:      "TEST1",
		Email:       "test2@email.com",
		Name:        "test2",
		Description: "test2descr",
		Actual:      true,
		EmailList:   []string{"2email1", "2email2", "2email3"},
		Status:      models.NoDecision,
	}
	err := s.db.NewTask(s.ctx, testTask)
	s.NoError(err)
	result, err := Task(s.ctx, s.db.Pool, testTask.Name)
	s.NoError(err)
	s.Equal(testTask.CreateTs, result.CreateTs)
	newTaskStatus := models.TaskStatus{
		TaskId:   testTask.TaskId,
		Email:    testTask.EmailList[0],
		Status:   models.Approve,
		CreateTs: time.Now(),
	}
	err = s.db.NewTaskStatus(s.ctx, newTaskStatus)
	s.NoError(err)
	resultStatus, err := Status(s.ctx, s.db.Pool, newTaskStatus.TaskId, newTaskStatus.Email)
	s.NoError(err)
	s.Equal(resultStatus.Status, newTaskStatus.Status)
	s.Equal(resultStatus.Email, newTaskStatus.Email)
}

func (s *integraTestSuite) TestSummaryTime() {
	defer DeleteTaskStatus(s.ctx, s.db.Pool)
	defer DeleteTaskAll(s.ctx, s.db.Pool)
	create1 := time.Now().Add(-24 * time.Hour)
	create2 := time.Now().Add(-12 * time.Hour)
	create3 := time.Now().Add(-6 * time.Hour)
	testTask := models.Task{
		TaskId: "TEST",
	}
	testTaskStatus1 := models.TaskStatus{
		TaskId:   testTask.TaskId,
		Email:    "test1@email.com",
		Status:   models.Approve,
		CreateTs: create1,
		DiffTime: create2.Sub(create1).Seconds(),
	}
	testTaskStatus2 := models.TaskStatus{
		TaskId:   testTaskStatus1.TaskId,
		Email:    "test2@email.com",
		Status:   models.Approve,
		CreateTs: create2,
		DiffTime: create3.Sub(create2).Seconds(),
	}
	testTaskStatus3 := models.TaskStatus{
		TaskId:   testTaskStatus1.TaskId,
		Email:    "test3@email.com",
		Status:   models.Approve,
		CreateTs: create3,
		DiffTime: create3.Sub(time.Now()).Seconds(),
	}
	err := s.db.NewTask(s.ctx, testTask)
	s.NoError(err)
	err = s.db.NewTaskStatus(s.ctx, testTaskStatus1)
	s.NoError(err)
	err = s.db.NewTaskStatus(s.ctx, testTaskStatus2)
	s.NoError(err)
	err = s.db.NewTaskStatus(s.ctx, testTaskStatus3)
	s.NoError(err)
	_, err = s.db.SummaryTime(s.ctx)
	s.NoError(err)
	// s.Equal(durs[0].Duration, create1.Sub(time.Now()).Seconds())
	// DeleteTaskStatus(s.ctx, s.db.Pool)

}

func (s *integraTestSuite) TestApprovedDeclinedTasks() {
	testTask := models.Task{
		TaskId:      "TEST1",
		Email:       "test1@email.com",
		Name:        "test1",
		Description: "test1descr",
		Actual:      true,
		EmailList:   []string{"1email1", "1email2", "1email3"},
		Status:      models.Approve,
	}
	testTask2 := models.Task{
		TaskId:      "TEST2",
		Email:       "test2@email.com",
		Name:        "test2",
		Description: "test2descr",
		Actual:      true,
		EmailList:   []string{"2email1", "2email2", "2email3"},
		Status:      models.Decline,
	}
	testTask3 := models.Task{
		TaskId:      "TEST3",
		Email:       "test3@email.com",
		Name:        "test3",
		Description: "test3descr",
		Actual:      true,
		EmailList:   []string{"3email1", "3email2", "3email3"},
		Status:      models.NoDecision,
	}
	testTask4 := models.Task{
		TaskId:      "TEST4",
		Email:       "test4@email.com",
		Name:        "test4",
		Description: "test4descr",
		Actual:      true,
		EmailList:   []string{"4email1", "4email2", "4email3"},
		Status:      models.Approve,
	}
	DeleteTaskStatus(s.ctx, s.db.Pool)
	err := s.db.NewTask(s.ctx, testTask)
	s.NoError(err)
	err = s.db.NewTask(s.ctx, testTask2)
	s.NoError(err)
	err = s.db.NewTask(s.ctx, testTask3)
	s.NoError(err)
	err = s.db.NewTask(s.ctx, testTask4)
	s.NoError(err)
	counter, err := s.db.ApprovedTasks(s.ctx)
	s.NoError(err)
	s.Equal(counter, models.Counter{Count: int32(2)})
	counter2, err := s.db.DeclinedTasks(s.ctx)
	s.NoError(err)
	s.Equal(counter2, models.Counter{Count: int32(2)})
	// DeleteTaskStatus(s.ctx, s.db.Pool)
}

func (s *integraTestSuite)TestNewMain(){
	mail := models.Mail{
		Header: "test",
		Body: "test",
		CreateTS: time.Now(),
		EmailList: []string{"test1", "test2"},
	}
	err := s.db.NewMail(s.ctx, mail)
	s.NoError(err)
	returnedMail, err := Mail(s.ctx, s.db.Pool, mail.Header)
	s.NoError(err)
	s.Equal(returnedMail.Body, mail.Body)
	DeleteMails(s.ctx, s.db.Pool)
}