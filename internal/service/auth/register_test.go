package auth

import (
	"context"
	"testing"

	"github.com/huypham67/bookmark-common/pkg/dbutils"
	jwtMocks "github.com/huypham67/bookmark-common/pkg/jwt/mocks"
	securityMocks "github.com/huypham67/bookmark-common/pkg/security/mocks"
	authDTO "github.com/huypham67/user-service/internal/dto/auth"
	"github.com/huypham67/user-service/internal/model"
	userMocks "github.com/huypham67/user-service/internal/repository/user/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func expectedAuthRegisteredUser() *model.User {
	return &model.User{
		DisplayName: "Test Display Name",
		Username:    "testuser",
		Email:       "testuser@gmail.com",
		Password:    "$2a$10$hashedpassword123456789",
	}
}

func matchAuthUser(expected *model.User) interface{} {
	return mock.MatchedBy(func(actual *model.User) bool {
		return actual.DisplayName == expected.DisplayName &&
			actual.Username == expected.Username &&
			actual.Email == expected.Email &&
			actual.Password == expected.Password
	})
}

func TestService_RegisterUser(t *testing.T) {
	t.Parallel()

	type args struct {
		request authDTO.RegisterUserRequest
	}

	testCases := []struct {
		name           string
		args           args
		setupMocks     func(context.Context, *userMocks.Repository, *securityMocks.PasswordHasher, *jwtMocks.TokenGenerator)
		verifyResponse func(*testing.T, *model.User, error)
	}{
		{
			name: "should register user successfully",
			args: args{
				request: authDTO.RegisterUserRequest{
					DisplayName: "Test Display Name",
					Username:    "testuser",
					Email:       "testuser@gmail.com",
					Password:    "password123",
				},
			},
			setupMocks: func(
				ctx context.Context,
				userRepo *userMocks.Repository,
				passwordHasher *securityMocks.PasswordHasher,
				tokenGenerator *jwtMocks.TokenGenerator,
			) {
				passwordHasher.
					On("Hash", "password123").
					Return("$2a$10$hashedpassword123456789", nil).
					Once()

				expectedUser := expectedAuthRegisteredUser()

				userRepo.
					On(
						"Create",
						ctx,
						matchAuthUser(expectedUser),
					).
					Return(nil).
					Once()
			},
			verifyResponse: func(
				t *testing.T,
				user *model.User,
				err error,
			) {
				assert.NoError(t, err)
				require.NotNil(t, user)

				assert.Equal(t, "Test Display Name", user.DisplayName)
				assert.Equal(t, "testuser", user.Username)
				assert.Equal(t, "testuser@gmail.com", user.Email)
				assert.Equal(t, "$2a$10$hashedpassword123456789", user.Password)
			},
		},
		{
			name: "should return error when user already exists",
			args: args{
				request: authDTO.RegisterUserRequest{
					DisplayName: "Test Display Name",
					Username:    "testuser",
					Email:       "testuser@gmail.com",
					Password:    "password123",
				},
			},
			setupMocks: func(
				ctx context.Context,
				userRepo *userMocks.Repository,
				passwordHasher *securityMocks.PasswordHasher,
				tokenGenerator *jwtMocks.TokenGenerator,
			) {
				passwordHasher.
					On("Hash", "password123").
					Return("$2a$10$hashedpassword123456789", nil).
					Once()

				expectedUser := expectedAuthRegisteredUser()

				userRepo.
					On(
						"Create",
						ctx,
						matchAuthUser(expectedUser),
					).
					Return(dbutils.ErrDuplicationType).
					Once()
			},
			verifyResponse: func(
				t *testing.T,
				user *model.User,
				err error,
			) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.ErrorIs(t, err, ErrUserAlreadyExists)
			},
		},
		{
			name: "should return error when password hashing fails",
			args: args{
				request: authDTO.RegisterUserRequest{
					DisplayName: "Test Display Name",
					Username:    "testuser",
					Email:       "testuser@gmail.com",
					Password:    "password123",
				},
			},
			setupMocks: func(
				ctx context.Context,
				userRepo *userMocks.Repository,
				passwordHasher *securityMocks.PasswordHasher,
				tokenGenerator *jwtMocks.TokenGenerator,
			) {
				passwordHasher.
					On("Hash", "password123").
					Return("", assert.AnError).
					Once()
			},
			verifyResponse: func(
				t *testing.T,
				user *model.User,
				err error,
			) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.ErrorIs(t, err, ErrInternalServerError)
			},
		},
		{
			name: "should return error when creating user in database fails",
			args: args{
				request: authDTO.RegisterUserRequest{
					DisplayName: "Test Display Name",
					Username:    "testuser",
					Email:       "testuser@gmail.com",
					Password:    "password123",
				},
			},
			setupMocks: func(
				ctx context.Context,
				userRepo *userMocks.Repository,
				passwordHasher *securityMocks.PasswordHasher,
				tokenGenerator *jwtMocks.TokenGenerator,
			) {
				passwordHasher.
					On("Hash", "password123").
					Return("$2a$10$hashedpassword123456789", nil).
					Once()

				expectedUser := expectedAuthRegisteredUser()

				userRepo.
					On(
						"Create",
						ctx,
						matchAuthUser(expectedUser),
					).
					Return(assert.AnError).
					Once()
			},
			verifyResponse: func(
				t *testing.T,
				user *model.User,
				err error,
			) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.ErrorIs(t, err, ErrInternalServerError)
			},
		},
		{
			name: "should return error when context is cancelled",
			args: args{
				request: authDTO.RegisterUserRequest{
					DisplayName: "Test Display Name",
					Username:    "testuser",
					Email:       "testuser@gmail.com",
					Password:    "password123",
				},
			},
			setupMocks: func(
				ctx context.Context,
				userRepo *userMocks.Repository,
				passwordHasher *securityMocks.PasswordHasher,
				tokenGenerator *jwtMocks.TokenGenerator,
			) {
				passwordHasher.
					On("Hash", "password123").
					Return("$2a$10$hashedpassword123456789", nil).
					Once()

				expectedUser := expectedAuthRegisteredUser()

				userRepo.
					On(
						"Create",
						ctx,
						matchAuthUser(expectedUser),
					).
					Return(context.Canceled).
					Once()
			},
			verifyResponse: func(
				t *testing.T,
				user *model.User,
				err error,
			) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.ErrorIs(t, err, ErrInternalServerError)
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var ctx context.Context
			userRepo := new(userMocks.Repository)
			passwordHasher := securityMocks.NewPasswordHasher(t)
			tokenGenerator := jwtMocks.NewTokenGenerator(t)

			// For context cancellation test, create a cancelled context
			if tc.name == "should return error when context is cancelled" {
				cancelledCtx, cancel := context.WithCancel(context.Background())
				cancel()
				ctx = cancelledCtx
			} else {
				ctx = context.Background()
			}

			tc.setupMocks(ctx, userRepo, passwordHasher, tokenGenerator)

			authService := NewService(userRepo, passwordHasher, tokenGenerator)

			user, err := authService.RegisterUser(
				ctx,
				tc.args.request,
			)

			tc.verifyResponse(t, user, err)

			userRepo.AssertExpectations(t)
		})
	}
}
