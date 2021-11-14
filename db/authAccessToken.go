package db

import (
	"gopkg.in/guregu/null.v3"
	"time"
)

// AuthAccessToken is an object representing the database table.
type AuthAccessToken struct {
	ID           int       `json:"id" db:"id"`
	AuthClientID int       `json:"auth_client_id" db:"auth_client_id"`
	AccessToken  string    `json:"access_token" db:"access_token"`
	ExpiresIn    string    `json:"expires_in" db:"expires_in"`
	ExpiresOn    time.Time `json:"expires_on" db:"expires_on"`
	CreatedAt    null.Time `json:"created_at" db:"expires_in"`
	UpdatedAt    null.Time `json:"updated_at" db:"expires_in"`
}
