package backend_masterclass

import (
	"backend_masterclass/util"
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	arg := CreateTransferParams{
		FromAccount: sql.NullInt64{
			Int64: util.RandomInt(1, 10),
			Valid: true,
		},
		ToAccount: sql.NullInt64{
			Int64: util.RandomInt(1, 10),
			Valid: true,
		},
		Amout: util.RandomMoney(),
	}

	transfers, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)
	require.Equal(t, transfers.Amout, arg.Amout)
	require.Equal(t, transfers.FromAccount, arg.FromAccount)
	require.Equal(t, transfers.ToAccount, arg.ToAccount)
}
