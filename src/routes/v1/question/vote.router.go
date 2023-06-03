package question

import (
	"github.com/labstack/echo/v4"
	"github.com/rpratama-dev/mymovie/src/controllers/v1"
)

func VoteRouter(route *echo.Group)  {
	route.POST("", controllers.QuestionVoteStore)
}
