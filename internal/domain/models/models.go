package models

import (
	"time"

	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
)

type Decisioner struct {
	TaskID         uuid.UUID
	Email          string
	EmailOrder     int
	EmailSentTS    time.Time
	DecisionMadeTS time.Time
	Decision       int
}
