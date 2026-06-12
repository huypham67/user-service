package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/huypham67/user-service/internal/handler/auth"
	"github.com/huypham67/user-service/internal/handler/health"
	"github.com/huypham67/user-service/internal/handler/profile"
)

// Router wraps the Gin engine and application server configuration.
type Router struct {
	engine *gin.Engine
}

// NewRouter creates and configures a new HTTP router with all API endpoints.
func NewRouter() *Router {
	engine := gin.Default()

	engine.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	return &Router{
		engine: engine,
	}
}

// GroupAPI returns a router group for endpoints (no /api prefix - handled by API Gateway).
func (r *Router) GroupAPI() *gin.RouterGroup {
	return r.engine.Group("/api/user_service")
}

// GroupV1 returns a router group for API version 1 endpoints.
func (r *Router) GroupV1() *gin.RouterGroup {
	return r.GroupAPI().Group("/v1")
}

// RegisterHealthRoutes registers all health check routes.
func RegisterHealthRoutes(
	apiGroup *gin.RouterGroup,
	handler health.Handler,
) {
	apiGroup.GET(
		"/health-check",
		handler.GetHealthCheck,
	)
}

// RegisterAuthRoutes registers all authentication routes (registration and login).
func RegisterAuthRoutes(
	routerGroup *gin.RouterGroup,
	handler auth.Handler,
) {
	routerGroup.POST(
		"/users/register",
		handler.Register,
	)

	routerGroup.POST(
		"/users/login",
		handler.Login,
	)
}

// RegisterProfileRoutes registers all user profile routes.
func RegisterProfileRoutes(
	routerGroup *gin.RouterGroup,
	handler profile.Handler,
	jwtMiddleware gin.HandlerFunc,
) {
	routerGroup.GET(
		"/self/info",
		jwtMiddleware,
		handler.GetUserInfo,
	)

	routerGroup.PUT(
		"/self/info",
		jwtMiddleware,
		handler.UpdateUserInfo,
	)
}

// ServeHTTP implements the http.Handler interface.
func (r *Router) ServeHTTP(
	writer http.ResponseWriter,
	request *http.Request,
) {
	r.engine.ServeHTTP(writer, request)
}

// Engine exposes underlying Gin engine
func (r *Router) Engine() *gin.Engine {
	return r.engine
}
