package db

import "fmt"

func (relation *Relationship) Create(txo *TxO) error {

	q := fmt.Sprintf(`
		INSERT INTO relationship
		(
			followed_id,
			follower_id
		)
		VALUES
			(
			:followed_id,
			:follower_id
		)`,
	)

	r, err := txo.Tx.NamedExec(q, relation)
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
	relation.ID = int(lastID)

	return err
}
//
//func GetAllRelationWithFollowerId(txo *TxO, id int) (*Relationship, error) {
//
//	q := fmt.Sprintf(`
//		SELECT
//			user_id,
//			name,
//			email,
//			address,
//			created_by,
//			created_timestamp,
//			last_updated_by,
//			last_updated_timestamp
//		FROM
//			user
//		WHERE
//			user_id = ?`,
//	)
//
//	var users []*User
//	if err := txo.Select(&users, q, id); err != nil {
//		return nil, fmt.Errorf("error executing sql query")
//	}
//
//	if len(users) == 0 {
//		return nil, fmt.Errorf("no rows found")
//	}
//
//	return users[0], nil
//}
