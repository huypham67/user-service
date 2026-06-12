package profile

import (
	"context"
	"errors"
	"testing"

	"github.com/huypham67/user-service/internal/model"
	"github.com/huypham67/user-service/internal/repository/user/mocks"
	"github.com/stretchr/testify/assert"
)

func TestService_GetUserInfo(t *testing.T) {
	t.Parallel()

	type args struct {
		userID string
	}

	testCases := []struct {
		name           string
		args           args
		setupMocks     func(context.Context, *mocks.Repository)
		verifyResponse func(*testing.T, *model.User, error)
	}{
		{
			name: "should get user info successfully",
			args: args{
				userID: "user-id",
			},
			setupMocks: func(
				ctx context.Context,
				mockRepo *mocks.Repository,
			) {
				mockRepo.
					On(
						"GetByID",
						ctx,
						"user-id",
					).
					Return(
						&model.User{
							BaseModel: model.BaseModel{
								ID: "user-id",
							},
							DisplayName: "John Doe",
							Username:    "johndoe",
							Email:       "john@example.com",
						},
						nil,
					).
					Once()
			},
			verifyResponse: func(
				t *testing.T,
				user *model.User,
				err error,
			) {
				assert.NoError(t, err)

				assert.Equal(
					t,
					"user-id",
					user.ID,
				)
			},
		},
		{
			name: "should return error when repository fails",
			args: args{
				userID: "user-id",
			},
			setupMocks: func(
				ctx context.Context,
				mockRepo *mocks.Repository,
			) {
				mockRepo.
					On(
						"GetByID",
						ctx,
						"user-id",
					).
					Return(
						nil,
						errors.New("database error"),
					).
					Once()
			},
			verifyResponse: func(
				t *testing.T,
				user *model.User,
				err error,
			) {
				assert.Error(t, err)
				assert.Nil(t, user)
			},
		},
		{
			name: "should return error when context is cancelled",
			args: args{
				userID: "user-id",
			},
			setupMocks: func(
				ctx context.Context,
				mockRepo *mocks.Repository,
			) {
				mockRepo.
					On(
						"GetByID",
						ctx,
						"user-id",
					).
					Return(
						nil,
						context.Canceled,
					).
					Once()
			},
			verifyResponse: func(
				t *testing.T,
				user *model.User,
				err error,
			) {
				assert.Error(t, err)
				assert.Nil(t, user)
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var ctx context.Context
			mockRepo := new(mocks.Repository)

			// For context cancellation test, create a cancelled context
			if tc.name == "should return error when context is cancelled" {
				cancelledCtx, cancel := context.WithCancel(context.Background())
				cancel()
				ctx = cancelledCtx
			} else {
				ctx = context.Background()
			}

			tc.setupMocks(ctx, mockRepo)

			service := NewService(mockRepo)

			user, err := service.GetUserInfo(ctx, tc.args.userID)

			tc.verifyResponse(t, user, err)

			mockRepo.AssertExpectations(t)
		})
	}
}
