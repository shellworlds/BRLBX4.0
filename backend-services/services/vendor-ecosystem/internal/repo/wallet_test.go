package repo

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestApplyMealNetCredits_MockTx(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()
	s := &Store{DB: mock}
	vid := uuid.New()
	txid := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO vendor_wallets`).WithArgs(vid).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectExec(`UPDATE vendor_wallets SET balance`).WithArgs(vid, 9.5).WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectExec(`INSERT INTO wallet_ledger`).WithArgs(vid, 9.5, txid).WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()

	err = s.ApplyMealNetCredits(context.Background(), vid, 10.0, 500, txid)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
