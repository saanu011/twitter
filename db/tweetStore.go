package db

import "fmt"

func (tweet *Tweet) Create(txo *TxO) error {

	if err := tweet.validate(); err != nil {
		return err
	}

	q := fmt.Sprintf(`
		INSERT INTO tweet
		(
			content,
			user_id
		)
		VALUES
			(
			:content,
			:user_id
		)`,
	)

	r, err := txo.Tx.NamedExec(q, tweet)
	if err != nil {
		return err
	}

	rows, err := r.RowsAffected()
	if err != nil {
		return err
	}

	if rows < 1 {
		return err
	}

	lastID, err := r.LastInsertId()
	if err != nil {
		return err
	}
	tweet.ID = int(lastID)

	return err
}

func GetUserTweets(txo *TxO, id int) ([]*FollowedUserTweet, error) {

	q := fmt.Sprintf(`
		SELECT
			t.user_id,
			t.content,
			t.created_timestamp,
			t.last_updated_timestamp,
			u.name user_name,
			u.email user_email
		FROM
			tweet t
				JOIN user u ON t.user_id = u.user_id
		WHERE
			u.user_id = ?`,
	)

	var tweetFeed []*FollowedUserTweet
	if err := txo.Select(&tweetFeed, q, id); err != nil {
		return nil, fmt.Errorf("error executing sql query: %v", err)
	}

	return tweetFeed, nil
}

func GetFollowedUserTweets(txo *TxO, id int) ([]*FollowedUserTweet, error) {

	q := fmt.Sprintf(`
		SELECT
			t.id,
			t.user_id,
			t.content,
			t.created_timestamp,
			t.last_updated_timestamp,
			u.name user_name,
			u.email user_email
		FROM
			tweet t
				JOIN relationship r ON t.user_id = r.followed_id
				JOIN user u ON t.user_id = u.user_id
		WHERE
			r.follower_id = ?`,
	)

	var tweetFeed []*FollowedUserTweet
	if err := txo.Select(&tweetFeed, q, id); err != nil {
		return nil, fmt.Errorf("error executing sql query: %v", err)
	}

	return tweetFeed, nil
}
