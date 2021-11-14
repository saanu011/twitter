package db

import (
	"gopkg.in/guregu/null.v3"
	"time"
)

type User struct {
	UserID               int         `json:"user_id" db:"user_id"`
	Name                 string      `json:"name" db:"name"`
	Email                string      `json:"email" db:"email"`
	Address              null.String `json:"address" db:"address"`
	LastUpdatedBy        string      `json:"last_updated_by" db:"last_updated_by"`
	LastUpdatedTimestamp time.Time   `json:"last_updated_timestamp" db:"last_updated_timestamp"`
	CreatedBy            string      `json:"created_by" db:"created_by"`
	CreatedTimestamp     time.Time   `json:"created_timestamp" db:"created_timestamp"`

	// Extra fields that may be included, depending on query params
}
