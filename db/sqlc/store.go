package backend_masterclass

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
}

type SQLStore struct {
	//Using a queries struct like this is called composition. It is said to be a better
	//decision than inheritance
	//This line exactly is the composition, adding this pointer gives the Store struct the behaviour of the Queries struct. Make sesne pa
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// var txKey = struct{}{}

//This function creates a new DB transaction

// This function is not exported as we don't want other classes to be able to call
// the function directly
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	//[&sql.TxOptions{}] allows us to set a custom isolation level for this
	//transaction
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		rbError := tx.Rollback()
		if rbError != nil {
			//Combining errors to send to the user
			return fmt.Errorf("tx err: %v, rollback error: %v", err, rbError)
		}

		return err
	}

	return tx.Commit()

}
