package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

// TxO is a transaction object that should always be passed along
// to database-interacting functions
type TxO struct {
	*sqlx.Tx
	Email string
}

func BeginTx(db *sqlx.DB) (*TxO, bool) {
	tx, err := db.Beginx()
	if err == nil {
		txo := &TxO{
			Tx: tx,
		}

		return txo, false
	}

	return nil, true
}

func TerminateTxIfError(txo *TxO, err error) bool {
	if err == nil {
		return false
	}

	// Because Tx errors have higher precedence than other errors, this code
	// block does not need to worry about sending the original error back
	// to the client.  It merely needs to log the original to the console
	// and send the error message for tx rollback failure back to the client.
	if txErr := txo.Rollback(); txErr != nil {
		return true
	}

	return true
}

func TerminateTx(txo *TxO, err error) bool {
	if err == nil {

		// by having tx.Commit() write to the original error (which was nil),
		// the downstream parts of this function outside of this if block do
		// not have to be modified to account for any potential errors here.
		if err = txo.Commit(); err == nil {
			return false
		}

		// If execution reaches this point, then the error was originally nil,
		// but the transaction could not be committed.  In this case, passively
		// leave this if block and roll back as if the error were originally
		// non-nil
		err = fmt.Errorf("Terminate with Tx", err)
	}

	// Because Tx errors have higher precedence than other errors, this code
	// block does not need to worry about sending the original error back
	// to the client.  It merely needs to log the original to the console
	// and send the error message for tx rollback failure back to the client.
	if txErr := txo.Rollback(); txErr != nil {
		return true
	}

	return true
}
