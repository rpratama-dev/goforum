package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rpratama-dev/mymovie/src/middleware"
)

func RegisterRoutes(g *echo.Group) {
	AuthRouter(g.Group("/auth"))
	g.Use(middleware.ApiKeyMiddleWare, middleware.AuthMiddleware)
	UserRouter(g.Group("/users"))
	TagRouter(g.Group("/tags"))
	QuestionRouter(g.Group("/questions"))
}
