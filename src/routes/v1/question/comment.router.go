package question

import (
	"github.com/labstack/echo/v4"
	"github.com/rpratama-dev/mymovie/src/controllers/v1"
)

func CommentRouter(route *echo.Group)  {
	route.POST("", controllers.QuestionCommentStore)
	route.PUT("/:comment_id", controllers.QuestionCommentUpdate)
}
