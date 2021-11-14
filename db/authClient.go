package db

import "gopkg.in/guregu/null.v3"

type AuthClient struct {
	ID               int       `json:"id" db:"id"`
	ClientRef        string    `json:"client_ref" db:"client_ref"`
	ClientSecret     string    `json:"client_secret" db:"client_secret"`
	CreatedTimestamp null.Time `json:"created_timestamp" db:"created_timestamp"`
	UpdatedAt        null.Time `json:"last_updated_timestamp" db:"last_updated_timestamp"`
}
