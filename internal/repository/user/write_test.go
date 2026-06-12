package user

import (
	"context"
	"testing"

	"github.com/huypham67/bookmark-common/pkg/dbutils"
	"github.com/huypham67/user-service/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestRepository_Create(t *testing.T) {
	t.Parallel()

	type args struct {
		user model.User
	}

	testCases := []struct {
		name   string
		args   args
		verify func(*testing.T, *gorm.DB, error, args)
	}{
		{
			name: "should create user successfully",
			args: args{
				user: model.User{
					DisplayName: "Test User 4",
					Username:    "testuser4",
					Email:       "testuser4@gmail.com",
					Password:    "hashed-password",
				},
			},
			verify: func(t *testing.T, db *gorm.DB, err error, a args) {
				require.NoError(t, err)

				var actual model.User
				err = db.Where("username = ?", "testuser4").First(&actual).Error
				require.NoError(t, err)

				assert.Equal(t, a.user.DisplayName, actual.DisplayName)
				assert.Equal(t, a.user.Username, actual.Username)
				assert.Equal(t, a.user.Email, actual.Email)
				assert.Equal(t, a.user.Password, actual.Password)
			},
		},
		{
			name: "should return error when email already exists",
			args: args{
				user: model.User{
					DisplayName: "Test User 5",
					Username:    "testuser5",
					Email:       "testuser1@gmail.com", // duplicate email
					Password:    "hashed-password",
				},
			},
			verify: func(t *testing.T, db *gorm.DB, err error, a args) {
				require.Error(t, err)
				assert.ErrorIs(t, err, dbutils.ErrDuplicationType)

				var actual model.User
				err = db.Where("username = ?", "testuser5").First(&actual).Error
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
			},
		},
		{
			name: "should return error when username already exists",
			args: args{
				user: model.User{
					DisplayName: "Test User 6",
					Username:    "testuser1", // duplicate username
					Email:       "testuser6@gmail.com",
					Password:    "hashed-password",
				},
			},
			verify: func(t *testing.T, db *gorm.DB, err error, a args) {
				require.Error(t, err)
				assert.ErrorIs(t, err, dbutils.ErrDuplicationType)

				var actual model.User
				err = db.Where("username = ?", "testuser1").First(&actual).Error
				assert.NoError(t, err)
				assert.NotEqual(t, "testuser6@gmail.com", actual.Email)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			repo, testDB := newTestRepository(t)

			err := repo.Create(ctx, &tc.args.user)

			tc.verify(t, testDB, err, tc.args)
		})
	}
}

func TestRepository_Update(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		userID string
		user   model.User
		verify func(*testing.T, *gorm.DB, error)
	}{
		{
			name:   "should update user successfully",
			userID: "user-uuid-1",
			user: model.User{
				BaseModel: model.BaseModel{
					ID: "user-uuid-1",
				},
				DisplayName: "Updated User 1",
				Email:       "updated1@gmail.com",
			},
			verify: func(t *testing.T, db *gorm.DB, err error) {
				require.NoError(t, err)

				var actual model.User
				err = db.First(&actual, "id = ?", "user-uuid-1").Error
				require.NoError(t, err)

				assert.Equal(t, "user-uuid-1", actual.ID)
				assert.Equal(t, "Updated User 1", actual.DisplayName)
				assert.Equal(t, "updated1@gmail.com", actual.Email)
			},
		},
		{
			name:   "should return error when user does not exist",
			userID: "nonexistent-id",
			user: model.User{
				BaseModel: model.BaseModel{
					ID: "nonexistent-id",
				},
				DisplayName: "Nonexistent User",
				Email:       "nonexistent@gmail.com",
			},
			verify: func(t *testing.T, db *gorm.DB, err error) {
				require.NoError(t, err)

				var actual model.User
				err = db.First(&actual, "id = ?", "nonexistent-id").Error
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
			},
		},
		{
			name:   "should return error when updating to duplicate email",
			userID: "user-uuid-1",
			user: model.User{
				BaseModel: model.BaseModel{
					ID: "user-uuid-1",
				},
				DisplayName: "Updated User 1",
				Email:       "testuser2@gmail.com", // duplicate email
			},
			verify: func(t *testing.T, db *gorm.DB, err error) {
				require.Error(t, err)
				assert.ErrorIs(t, err, dbutils.ErrDuplicationType)

				var actual model.User
				err = db.First(&actual, "id = ?", "user-uuid-1").Error
				require.NoError(t, err)
				assert.Equal(t, "testuser1@gmail.com", actual.Email) // email unchanged
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			repo, testDB := newTestRepository(t)

			err := repo.Update(ctx, &tc.user)

			tc.verify(t, testDB, err)
		})
	}
}
