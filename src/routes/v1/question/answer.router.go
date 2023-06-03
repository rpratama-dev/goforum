package question

import (
	"github.com/labstack/echo/v4"
	"github.com/rpratama-dev/mymovie/src/controllers/v1"
)

func AnswerRouter(route *echo.Group)  {
	// route.GET("", controllers.AnswerStore)
	route.POST("", controllers.AnswerStore)
	// route.GET("/:id", controllers.UserQuestionShow)
	// route.PUT("/:id", controllers.UserQuestionUpdate)
	// route.DELETE("/:id", controllers.UserQuestionDestroy)
}
