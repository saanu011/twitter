package db

import (
	"fmt"
	"time"
)

type Tweet struct {
	ID                   int       `json:"id" db:"id"`
	Content              string    `json:"content" db:"content"`
	UserID               int       `json:"user_id" db:"user_id"`
	LastUpdatedTimestamp time.Time `json:"last_updated_timestamp" db:"last_updated_timestamp"`
	CreatedTimestamp     time.Time `json:"created_timestamp" db:"created_timestamp"`
}

type FollowedUserTweet struct {
	Tweet
	UserName  string `json:"user_name" db:"user_name"`
	UserEmail string `json:"user_email" db:"user_email"`
}

func (tweet Tweet) validate() error {
	if len(tweet.Content) == 0 {
		return fmt.Errorf("content can't be empty")
	}
	return nil
}
