package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rpratama-dev/mymovie/src/routes/v1/user"
)

func UserRouter(route *echo.Group)  {
	user.QuestionRouter(route.Group("/questions"))
}
