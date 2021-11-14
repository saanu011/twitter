package db

import "fmt"

func (authClient *AuthClient) Create(txo *TxO) error {

	q := fmt.Sprintf(`
		INSERT INTO auth_client
		(
			client_ref,
			client_secret
		)
		VALUES
			(
			:client_ref,
			:client_secret
		)`,
	)

	r, err := txo.Tx.NamedExec(q, authClient)
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
	authClient.ID = int(lastID)

	return err
}

func GetAuthClientByClientRef(txo *TxO, email string) (*AuthClient, error) {

	q := fmt.Sprintf(`
		SELECT
			id,
			client_ref,
			client_secret,
			created_timestamp,
			last_updated_timestamp
		FROM
			auth_client
		WHERE
			client_ref = ?`,
	)

	var authClient []*AuthClient
	if err := txo.Select(&authClient, q, email); err != nil {
		return nil, fmt.Errorf("error executing sql query: %v", err)
	}

	if len(authClient) == 0 {
		return nil, fmt.Errorf("no rows found")
	}

	return authClient[0], nil
}
