package user

import (
	"github.com/labstack/echo/v4"
	"github.com/rpratama-dev/goforum/src/controllers/v1"
)

func QuestionRouter(route *echo.Group)  {
	route.GET("", controllers.UserQuestionIndex)
	route.POST("", controllers.UserQuestionStore)
	route.GET("/:id", controllers.UserQuestionShow)
	route.PUT("/:id", controllers.UserQuestionUpdate)
	route.DELETE("/:id", controllers.UserQuestionDestroy)
}
