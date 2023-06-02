package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rpratama-dev/mymovie/src/controllers/v1"
)

func TagRouter(route *echo.Group)  {
	route.GET("", controllers.TagIndex)
	route.POST("", controllers.TagStore)
	route.GET("/:id", controllers.TagShow)
}
