package db

import (
	"context"
	"database/sql"
	"simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	// err 必须为 nil
	require.NoError(t, err)

	// account 不能为空对象
	require.NotEmpty(t, account)

	// 账户的所有者、余额和币种是否与输入的一致
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	// 检查ID是否自动生成的，必须不为0
	require.NotZero(t, account.ID)

	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	// 查询 account, 参数为 account1 的 id，把结果给 account2
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)

	require.NotEmpty(t, account2)

	require.Equal(t, account2.ID, account1.ID)
	require.Equal(t, account2.Owner, account1.Owner)
	require.Equal(t, account2.Balance, account1.Balance)
	require.Equal(t, account2.Currency, account1.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	// 比较 account2 和 account1, 除了 Balance，其他的字段都应该相同
	require.Equal(t, account2.ID, account1.ID)
	require.Equal(t, account2.Owner, account1.Owner)
	// 这里使用 arg.Balance 和 account2.Balance 比较
	require.Equal(t, account2.Balance, arg.Balance)
	require.Equal(t, account2.Currency, account1.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	// 为了验证账户确实被删除了，再查找一次
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	// 因为已经删除掉了，这里必须有错误
	require.Error(t, err)
	// 更准确的说，错误应该是 sql.ErrNoRows
	require.EqualError(t, err, sql.ErrNoRows.Error())
	// account2 也应该是空的
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	// accounts 切片的长度为 5
	require.Len(t, accounts, 5)

	// 变量 accounts, 其中的每个 account 都不能为空
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}
