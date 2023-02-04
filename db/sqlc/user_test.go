package sqlc

import (
	"bank/util"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FullName, user.FullName)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

// func TestUpdateUser(t *testing.T) {
// 	account1 := createRandomAccount(t)
// 	arg := UpdateAccountParams{
// 		ID:      account1.ID,
// 		Balance: util.RandomMoney(),
// 	}

// 	account2, err := testQueries.UpdateAccount(context.Background(), arg)

// 	require.NoError(t, err)
// 	require.NotEmpty(t, account2)
// 	require.Equal(t, arg.Balance, account2.Balance)
// 	require.Equal(t, account1.ID, account2.ID)
// 	require.Equal(t, account1.Currency, account2.Currency)
// 	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
// }

// func TestDeleteUser(t *testing.T) {
// 	account1 := createRandomAccount(t)
// 	err := testQueries.DeleteAccount(context.Background(), account1.ID)
// 	require.NoError(t, err)

// 	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

// 	require.Error(t, err)
// 	require.EqualError(t, err, sql.ErrNoRows.Error())
// 	require.Empty(t, account2)
// }

// func TestListUsers(t *testing.T) {
// 	for i := 0; i < 10; i++ {
// 		createRandomAccount(t)
// 	}

// 	arg := ListAccountParams{
// 		Limit:  5,
// 		Offset: 5,
// 	}

// 	accounts, err := testQueries.ListAccount(context.Background(), arg)

// 	require.NoError(t, err)
// 	require.Len(t, accounts, 5)

// 	for _, account := range accounts {
// 		require.NotEmpty(t, account)
// 	}
// }

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)
	newFullName := util.RandomOwner()
	updateUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, updateUser.FullName, oldUser.FullName)
	require.Equal(t, updateUser.Username, oldUser.Username)
	require.Equal(t, oldUser.Email, updateUser.Email)
	require.Equal(t, oldUser.HashedPassword, updateUser.HashedPassword)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)

	newEmail := util.RandomEmail()
	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := util.RandomOwner()
	newEmail := util.RandomEmail()
	newPassword := util.RandomString(6)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	updatedUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
	})

	require.NoError(t, err)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, newFullName, updatedUser.FullName)
}
