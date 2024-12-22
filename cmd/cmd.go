package cmd

import (
	"context"

	"github.com/labstack/echo/v4"
	"quote-management-tech-task/pkg/health"
	"quote-management-tech-task/pkg/products"

	"quote-management-tech-task/config"
	"quote-management-tech-task/db/sqlc"

	"github.com/labstack/echo-contrib/echoprometheus"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

const versionAPI = "v1"

type Server struct {
	echo   *echo.Echo
	db     *sqlc.Queries
	config config.Config
}

func NewServer(appConfig config.Config, db *sqlc.Queries) Server {
	e := echo.New()
	addMiddleware(e)

	e.GET("/health", health.HealthHandler)
	e.GET("/metrics", echoprometheus.NewHandler())

	v1API := e.Group("/" + versionAPI)

	addAPIHandlers(
		v1API,
		products.NewHandler(products.NewService(db)),
	)

	return Server{
		echo:   e,
		db:     db,
		config: appConfig,
	}
}

func addMiddleware(echo *echo.Echo) {
	echo.Use(echoMiddleware.Logger())
	echo.Use(echoMiddleware.Recover())
	echo.Use(echoMiddleware.Secure())
	echo.Pre(echoMiddleware.RemoveTrailingSlash())
}

type Handler interface {
	Register(route *echo.Group)
}

func addAPIHandlers(v1 *echo.Group, handlers ...Handler) {
	for _, handler := range handlers {
		handler.Register(v1)
	}
}

func (s Server) Run() error {
	return s.echo.Start(":" + s.config.Addr)
}

func (s Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}
