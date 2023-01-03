package backend_masterclass

import "context"

// TransferTxParams contains the input parameters of the transfer transaction
type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user Users) error
}

// TransferTxResult is the result of the transfer transaction
type CreateUserTxResult struct {
	User Users
}

// This function executes the transfer transaction i.e the transfer of money from
// one account to another
func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {

		var err error

		user, err := q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		result.User = user

		err = arg.AfterCreate(result.User)
		if err != nil {
			return err
		}

		return nil

	})

	return result, err

}
