package integration

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/huypham67/bookmark-common/middleware"
	"github.com/huypham67/bookmark-common/pkg/jwt"
	"github.com/huypham67/bookmark-common/pkg/security"
	"github.com/huypham67/bookmark-common/pkg/sqldb"
	"github.com/huypham67/user-service/internal/api"
	authHandler "github.com/huypham67/user-service/internal/handler/auth"
	healthHandler "github.com/huypham67/user-service/internal/handler/health"
	profileHandler "github.com/huypham67/user-service/internal/handler/profile"
	"github.com/huypham67/user-service/internal/repository/ping"
	userRepo "github.com/huypham67/user-service/internal/repository/user"
	authSvc "github.com/huypham67/user-service/internal/service/auth"
	healthSvc "github.com/huypham67/user-service/internal/service/health"
	profileSvc "github.com/huypham67/user-service/internal/service/profile"
	"github.com/huypham67/user-service/internal/test/fixtures"
	"gorm.io/gorm"
)

const (
	testIssuer   = "test-issuer"
	testAudience = "test-audience"
)

// TestApp represents the test application with its dependencies.
type TestApp struct {
	Router *api.Router
	MockDB *gorm.DB
}

type AuthenticatedTestApp struct {
	*TestApp
	TokenGenerator jwt.TokenGenerator
}

func createTestJWT(t *testing.T) (
	jwt.TokenGenerator,
	jwt.TokenValidator,
) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(
		rand.Reader,
		2048,
	)
	require.NoError(t, err)

	tokenGenerator, err := jwt.NewTokenGenerator(
		privateKey,
		testIssuer,
		testAudience,
		time.Hour,
	)
	require.NoError(t, err)

	tokenValidator, err := jwt.NewTokenValidator(
		&privateKey.PublicKey,
		testIssuer,
		testAudience,
	)
	require.NoError(t, err)

	return tokenGenerator, tokenValidator
}

func setupHealthCheckTestApp(t *testing.T, serviceName string, instanceID string) *TestApp {
	t.Helper()

	mockDB := sqldb.NewMock(t)

	pinger := ping.NewSQLDB(mockDB)

	healthService := healthSvc.NewService(serviceName, instanceID, pinger)

	healthHandlerInstance := healthHandler.NewHandler(healthService)

	router := api.NewRouter()

	api.RegisterHealthRoutes(
		router.GroupAPI(),
		healthHandlerInstance,
	)

	return &TestApp{
		Router: router,
		MockDB: mockDB,
	}
}

func createTestTokenGenerator(t *testing.T) jwt.TokenGenerator {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	generator, err := jwt.NewTokenGenerator(
		privateKey,
		"test-issuer",
		"test-audience",
		time.Hour,
	)
	require.NoError(t, err)

	return generator
}

func setupAuthTestApp(t *testing.T) *TestApp {
	t.Helper()

	mockDB := fixtures.NewTestDB(t, &fixtures.UserTestDB{})

	userRepository := userRepo.NewRepository(mockDB)

	passwordHasher := security.NewBcryptPasswordHasher()

	tokenGenerator := createTestTokenGenerator(t)

	authService := authSvc.NewService(userRepository, passwordHasher, tokenGenerator)

	authHandlerInstance := authHandler.NewHandler(authService)

	router := api.NewRouter()

	api.RegisterAuthRoutes(
		router.GroupV1(),
		authHandlerInstance,
	)

	return &TestApp{
		Router: router,
	}
}

func setupProfileTestApp(t *testing.T) *AuthenticatedTestApp {
	t.Helper()

	mockDB := fixtures.NewTestDB(t, &fixtures.UserTestDB{})

	userRepository := userRepo.NewRepository(mockDB)
	profileService := profileSvc.NewService(userRepository)
	profileHandlerInstance := profileHandler.NewHandler(profileService)

	tokenGenerator, tokenValidator := createTestJWT(t)

	router := api.NewRouter()

	jwtMiddleware := middleware.JWTAuth(tokenValidator)

	api.RegisterProfileRoutes(router.GroupV1(), profileHandlerInstance, jwtMiddleware)

	return &AuthenticatedTestApp{
		TestApp: &TestApp{
			Router: router,
		},
		TokenGenerator: tokenGenerator,
	}
}
