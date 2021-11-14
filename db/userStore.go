package db

import (
	"fmt"
)

func (user *User) Create(txo *TxO) error {

	q := fmt.Sprintf(`
		INSERT INTO user
		(
			name,
			email,
			address,
			last_updated_by,
			created_by
		)
		VALUES
			(
			:name,
			:email,
			:address,
			:last_updated_by,
			:created_by
		)`,
	)

	r, err := txo.Tx.NamedExec(q, user)
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
	user.UserID = int(lastID)

	return nil
}

func GetUserById(txo *TxO, id int) (*User, error) {

	q := fmt.Sprintf(`
		SELECT
			user_id,
			name,
			email,
			address,
			created_by,
			created_timestamp,
			last_updated_by,
			last_updated_timestamp
		FROM
			user
		WHERE
			user_id = ?`,
	)

	var users []*User
	if err := txo.Select(&users, q, id); err != nil {
		return nil, fmt.Errorf("error executing sql query")
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("no rows found")
	}

	return users[0], nil
}

func GetUserByEmail(txo *TxO, email string) (*User, error) {

	q := fmt.Sprintf(`
		SELECT
			user_id,
			name,
			email,
			address,
			created_by,
			created_timestamp,
			last_updated_by,
			last_updated_timestamp
		FROM
			user
		WHERE
			email = ?`,
	)

	var users []*User
	if err := txo.Select(&users, q, email); err != nil {
		return nil, fmt.Errorf("error executing sql query")
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("no rows found")
	}

	return users[0], nil
}

func GetUsers(txo *TxO) ([]*User, error) {

	q := fmt.Sprintf(`
		SELECT
			user_id,
			name,
			email,
			address,
			created_by,
			created_timestamp,
			last_updated_by,
			last_updated_timestamp
		FROM
			user;`,
	)

	var users []*User
	if err := txo.Select(&users, q); err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, nil
	}

	return users, nil
}
