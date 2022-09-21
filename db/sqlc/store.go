package backend_masterclass

import (
	"context"
	"database/sql"
	"fmt"
)

//Store provides all functions to execute db queries and transactions
type Store struct {
	//Using a queries struct like this is called composition. It is said to be a better
	//decision than inheritance
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

//This function creates a new DB transaction

//This function is not exported as we don't want other classes to be able to call
//the function directly
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	//[&sql.TxOptions{}] allows us to set a custom isolation level for this
	//transaction
	tx, err := store.db.BeginTx(ctx, nil)
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

// TransferToParams contains the input parameters of the transfer transaction
type TransferToParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

//TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfers   Transfers `json:"transfer"`
	FromAccount Accounts  `json:"from_account"`
	ToAccount   Accounts  `json:"to_account"`
	FromEntry   Entries   `json:"from_entry"`
	ToEntry     Entries   `json:"to_entry"`
}

//This function executes the transfer transaction i.e the transfer of money from
//one account to another
func (store *Store) TransferTx(ctx context.Context, arg TransferToParams) (TransferTxResult, error) {

	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		return nil
	})

	return result, err

}
