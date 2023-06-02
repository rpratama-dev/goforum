package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rpratama-dev/mymovie/src/controllers/v1"
	"github.com/rpratama-dev/mymovie/src/middleware"
)

func AuthRouter(route *echo.Group)  {
	route.POST("/sign-up", controllers.AuthSignUp)
	route.POST("/sign-in", controllers.AuthSignIn)
	route.POST("/verify/:token", controllers.AuthVerify)
	route.Use(middleware.AuthMiddleware)
	route.POST("/sign-out", controllers.AuthSignOut)
}
