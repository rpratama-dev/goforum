package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rpratama-dev/goforum/src/controllers/v1"
	router "github.com/rpratama-dev/goforum/src/routes/v1/question"
)

func QuestionRouter(route *echo.Group)  {
	route.GET("", controllers.UserQuestionIndex)
	route.GET("/:question_id", controllers.QuestionShow)
	router.AnswerRouter(route.Group("/:question_id/answers"))
	router.VoteRouter(route.Group("/:question_id/votes"))
	router.CommentRouter(route.Group("/:question_id/comments"))
}
