package models

import (
	"errors"
	"time"
)

type (
	Decision int
	Action int
	Kind int
)

var (
	ErrUnexpectedEventKind   error    = errors.New("unexpected event kind")
	ErrUnexpectedEventAction error    = errors.New("unexpected event action")
)
const (
	Approve                  Decision = 1
	Decline                  Decision = -1
	NoDecision               Decision = 0
	CreateAction             Action   = 0
	UpdateAction             Action   = 1
	TaskKind                 Kind     = 0
	StatusKind               Kind     = 1
)

type SummaryTime struct {
	// TODO поменять на реальный пример
	//TaskId   uuid.UUID     `json:"task_id" example:"test123"`
	TaskId   string  `json:"task_id" example:"test123"`
	Duration float64 `json:"duration" swaggertype:"primitive,integer" example:"1005"`
}

type Counter struct {
	Count int32 `json:"count" example:"5"`
}

type Mail struct {
	Header    string    `json:"header"`
	Body      string    `json:"body"`
	CreateTS  time.Time `json:"create_ts"`
	EmailList []string  `json:"email_list"`
}


type Task struct {
	//TaskId      uuid.UUID
	TaskId      string
	Email       string
	Name        string
	Description string
	CreateTs    time.Time
	Status      Decision
	EndTime     time.Time
	EmailList   []string
	Actual      bool
}

type TaskStatus struct {
	//TaskId   uuid.UUID
	TaskId   string
	Email    string
	Status   Decision
	CreateTs time.Time
	DiffTime float64
}

type TaskResponse struct {
	Success bool `json:"success"`
}
