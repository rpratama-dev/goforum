package question

import (
	"github.com/labstack/echo/v4"
	"github.com/rpratama-dev/mymovie/src/controllers/v1"
)

func AnswerRouter(route *echo.Group)  {
	route.POST("", controllers.AnswerStore)
	route.PUT("/:answer_id", controllers.AnswerUpdate)
	route.POST("/:answer_id/votes", controllers.AnswerVoteStore)
}
