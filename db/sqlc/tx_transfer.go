package backend_masterclass

import "context"

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfers   Transfers `json:"transfer"`
	FromAccount Accounts  `json:"from_account"`
	ToAccount   Accounts  `json:"to_account"`
	FromEntry   Entries   `json:"from_entry"`
	ToEntry     Entries   `json:"to_entry"`
}

// This function executes the transfer transaction i.e the transfer of money from
// one account to another
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {

	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {

		var err error

		// txName := ctx.Value(txKey)
		// fmt.Println(txName, "create transfer")

		result.Transfers, err = q.CreateTransfer(
			ctx,
			CreateTransferParams{
				FromAccount: arg.FromAccountID,
				ToAccount:   arg.ToAccountID,
				Amout:       arg.Amount,
			})

		if err != nil {
			return err
		}

		// fmt.Println(txName, "create from entry")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		// fmt.Println(txName, "create to transfer")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		//Made changes after refactoring
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)

		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		return err
	})

	return result, err

}

// Code refactoring
func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Accounts, account2 Accounts, err error) {
	account1, err = q.AddAccountBalance(
		ctx,
		AddAccountBalanceParams{
			Amount: amount1,
			ID:     accountID1,
		})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(
		ctx,
		AddAccountBalanceParams{
			Amount: amount2,
			ID:     accountID2,
		})

	if err != nil {
		return
	}
	return account1, account2, nil
}
