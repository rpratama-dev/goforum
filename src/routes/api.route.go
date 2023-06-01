package routes

import (
	"github.com/labstack/echo/v4"
	v1Router "github.com/rpratama-dev/mymovie/src/routes/v1"
)

func ApiRouter(router *echo.Group) {
	// Register here all API router
	v1Router.RegisterRoutes(router.Group("/v1"))
}
