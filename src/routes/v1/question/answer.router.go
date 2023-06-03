package question

import (
	"github.com/labstack/echo/v4"
	"github.com/rpratama-dev/goforum/src/controllers/v1"
)

func AnswerRouter(route *echo.Group)  {
	route.POST("", controllers.AnswerStore)
	route.PUT("/:answer_id", controllers.AnswerUpdate)
	route.PATCH("/:answer_id", controllers.AnswerPatch)
	route.POST("/:answer_id/votes", controllers.AnswerVoteStore)
	route.POST("/:answer_id/comments", controllers.AnswerCommentStore)
	route.PUT("/:answer_id/comments/:comment_id", controllers.AnswerCommentUpdate)
}
