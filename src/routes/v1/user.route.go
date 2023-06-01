package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rpratama-dev/mymovie/src/controllers/v1"
)

func Group(route *echo.Group)  {
	route.GET("", controllers.UserShow)
}
