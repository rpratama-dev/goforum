package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rpratama-dev/mymovie/src/controllers/v1"
	router "github.com/rpratama-dev/mymovie/src/routes/v1/question"
)

func QuestionRouter(route *echo.Group)  {
	route.GET("", controllers.UserQuestionIndex)
	route.GET("/:id", controllers.UserQuestionShow)
	router.AnswerRouter(route.Group("/:question_id/answers"))
	router.VoteRouter(route.Group("/:question_id/votes"))
}
