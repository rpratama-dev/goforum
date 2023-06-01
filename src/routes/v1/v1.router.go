package routes

import (
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(g *echo.Group) {
	AuthRouter(g.Group("/auth"));
	UserRouter(g.Group("/users"));
}
