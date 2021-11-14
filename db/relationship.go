package db

import (
	"time"
)

type Relationship struct {
	ID                   int       `json:"id" db:"id"`
	FollowedID           int       `json:"followed_id" db:"followed_id"`
	FollowerID           int       `json:"follower_id" db:"follower_id"`
	LastUpdatedTimestamp time.Time `json:"last_updated_timestamp" db:"last_updated_timestamp"`
	CreatedTimestamp     time.Time `json:"created_timestamp" db:"created_timestamp"`
}
