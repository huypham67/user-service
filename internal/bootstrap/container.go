package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/huypham67/bookmark-common/middleware"
	pkgRedis "github.com/huypham67/bookmark-common/pkg/redis"
	"github.com/huypham67/bookmark-common/pkg/sqldb"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/huypham67/bookmark-common/pkg/jwt"
	jwtprovider "github.com/huypham67/bookmark-common/pkg/jwt/provider"
	ratelimitprovider "github.com/huypham67/bookmark-common/pkg/ratelimit/provider"
	"github.com/huypham67/bookmark-common/pkg/security"
	authHandler "github.com/huypham67/user-service/internal/handler/auth"
	healthHandler "github.com/huypham67/user-service/internal/handler/health"
	profileHandler "github.com/huypham67/user-service/internal/handler/profile"
	"github.com/huypham67/user-service/internal/repository/ping"
	"github.com/huypham67/user-service/internal/repository/user"
	authSvc "github.com/huypham67/user-service/internal/service/auth"
	healthSvc "github.com/huypham67/user-service/internal/service/health"
	profileSvc "github.com/huypham67/user-service/internal/service/profile"
)

// Container holds the application's dependencies and initialized services.
// It serves as the single source of truth for all infrastructure and business logic components.
type Container struct {
	// Infrastructure
	Config *Config
	DB     *gorm.DB
	Redis  *redis.Client

	// Handlers
	HealthHandler  healthHandler.Handler
	AuthHandler    authHandler.Handler
	ProfileHandler profileHandler.Handler

	// Middleware
	JWTMiddleware       gin.HandlerFunc
	RateLimitMiddleware gin.HandlerFunc
}

// NewContainer initializes the application container by loading configuration,
// setting up infrastructure clients, and initializing all handlers with their dependencies.
func NewContainer() (*Container, error) {
	cfg, err := NewConfig()
	if err != nil {
		log.Error().Err(err).Msg("failed to load config")
		return nil, err
	}

	db, err := sqldb.NewClient("")
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize postgres client")
		return nil, err
	}

	db, err = sqldb.RunMigration(db, "migrations")
	if err != nil {
		log.Error().Err(err).Msg("failed to run database migrations")
		return nil, err
	}

	tokenGenerator, err := jwtprovider.NewIssuer("")
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize jwt issuer")
		return nil, err
	}

	tokenValidator, err := jwtprovider.NewValidator("")
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize jwt validator")
		return nil, err
	}

	jwtMiddleware := middleware.JWTAuth(tokenValidator)

	// Redis backs the rate limiter only; the health check pings the database, not Redis.
	rdb, err := pkgRedis.NewClient("")
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize redis client")
		return nil, err
	}

	rateLimiter, err := ratelimitprovider.New(rdb, "")
	if err != nil {
		log.Error().Err(err).Msg("failed to initialize rate limiter")
		return nil, err
	}
	rateLimitMiddleware := middleware.RateLimit(rateLimiter)

	healthHandlerInstance := initHealthHandler(cfg, db)
	authHandlerInstance, err := initAuthHandler(db, tokenGenerator)
	if err != nil {
		return nil, err
	}
	profileHandlerInstance := initProfileHandler(db)

	return &Container{
		Config:              cfg,
		DB:                  db,
		Redis:               rdb,
		HealthHandler:       healthHandlerInstance,
		AuthHandler:         authHandlerInstance,
		ProfileHandler:      profileHandlerInstance,
		JWTMiddleware:       jwtMiddleware,
		RateLimitMiddleware: rateLimitMiddleware,
	}, nil
}

func initAuthHandler(db *gorm.DB, tokenGenerator jwt.TokenGenerator) (authHandler.Handler, error) {
	userRepository := user.NewRepository(db)
	passwordHasher := security.NewBcryptPasswordHasher()

	authService := authSvc.NewService(userRepository, passwordHasher, tokenGenerator)
	return authHandler.NewHandler(authService), nil
}

func initProfileHandler(db *gorm.DB) profileHandler.Handler {
	userRepository := user.NewRepository(db)
	service := profileSvc.NewService(userRepository)
	return profileHandler.NewHandler(service)
}

func initHealthHandler(cfg *Config, database *gorm.DB) healthHandler.Handler {
	pinger := ping.NewSQLDB(database)
	healthService := healthSvc.NewService(cfg.ServiceName, cfg.InstanceID, pinger)
	return healthHandler.NewHandler(healthService)
}

// Close gracefully shuts down the database and Redis clients, ensuring that all resources are properly released.
func (c *Container) Close() error {
	if c.Redis != nil {
		_ = c.Redis.Close()
	}

	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}

	return nil
}
