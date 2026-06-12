package user

import (
	"context"
	"testing"

	"github.com/huypham67/bookmark-common/pkg/dbutils"
	"github.com/huypham67/user-service/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_GetByEmail(t *testing.T) {
	t.Parallel()

	type args struct {
		email string
	}

	testCases := []struct {
		name   string
		args   args
		verify func(*testing.T, *model.User, error)
	}{
		{
			name: "should return user when email exists",
			args: args{
				email: "testuser1@gmail.com",
			},
			verify: func(t *testing.T, user *model.User, err error) {
				require.NoError(t, err)
				require.NotNil(t, user)
				assert.Equal(t, "user-uuid-1", user.ID)
				assert.Equal(t, "Test User 1", user.DisplayName)
				assert.Equal(t, "testuser1", user.Username)
				assert.Equal(t, "testuser1@gmail.com", user.Email)
				assert.NotEmpty(t, user.Password)
			},
		},
		{
			name: "should return nil when email does not exist",
			args: args{
				email: "nonexistent@gmail.com",
			},
			verify: func(t *testing.T, user *model.User, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, dbutils.ErrRecordNotFoundType)
				assert.Nil(t, user)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			repo, _ := newTestRepository(t)

			user, err := repo.GetByEmail(ctx, tc.args.email)

			tc.verify(t, user, err)
		})
	}
}

func TestRepository_GetByUsername(t *testing.T) {
	t.Parallel()

	type args struct {
		username string
	}

	testCases := []struct {
		name   string
		args   args
		verify func(*testing.T, *model.User, error)
	}{
		{
			name: "should return user when username exists",
			args: args{
				username: "testuser1",
			},
			verify: func(t *testing.T, user *model.User, err error) {
				require.NoError(t, err)
				require.NotNil(t, user)
				assert.Equal(t, "user-uuid-1", user.ID)
				assert.Equal(t, "Test User 1", user.DisplayName)
				assert.Equal(t, "testuser1", user.Username)
				assert.Equal(t, "testuser1@gmail.com", user.Email)
				assert.NotEmpty(t, user.Password)
			},
		},
		{
			name: "should return nil when username does not exist",
			args: args{
				username: "nonexistent",
			},
			verify: func(t *testing.T, user *model.User, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, dbutils.ErrRecordNotFoundType)
				assert.Nil(t, user)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			repo, _ := newTestRepository(t)

			user, err := repo.GetByUsername(ctx, tc.args.username)

			tc.verify(t, user, err)
		})
	}
}

func TestRepository_GetByID(t *testing.T) {
	t.Parallel()

	type args struct {
		userID string
	}

	testCases := []struct {
		name   string
		args   args
		verify func(*testing.T, *model.User, error)
	}{
		{
			name: "should return user when ID exists",
			args: args{
				userID: "user-uuid-1",
			},
			verify: func(t *testing.T, user *model.User, err error) {
				require.NoError(t, err)
				require.NotNil(t, user)
				assert.Equal(t, "user-uuid-1", user.ID)
				assert.Equal(t, "Test User 1", user.DisplayName)
				assert.Equal(t, "testuser1", user.Username)
				assert.Equal(t, "testuser1@gmail.com", user.Email)
				assert.NotEmpty(t, user.Password)
			},
		},
		{
			name: "should return nil when ID does not exist",
			args: args{
				userID: "nonexistent-id",
			},
			verify: func(t *testing.T, user *model.User, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, dbutils.ErrRecordNotFoundType)
				assert.Nil(t, user)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			repo, _ := newTestRepository(t)

			user, err := repo.GetByID(ctx, tc.args.userID)

			tc.verify(t, user, err)
		})
	}
}
