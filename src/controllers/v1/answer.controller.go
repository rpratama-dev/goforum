package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	models "github.com/rpratama-dev/mymovie/src/models/table"
	"github.com/rpratama-dev/mymovie/src/utils"
)

func AnswerStore(c echo.Context) {
	defer utils.DeferHandler(c)
	// session := c.Get("session").(*models.Session)
	// Bind user input
	var answerPayload models.AnswerPayload
	c.Bind(&answerPayload)

	// Start validation input
	errValidation := answerPayload.Validate()
	if (errValidation != nil) {
		panic(utils.PanicPayload{
			Message: "Validation Error",
			Data: errValidation,
			HttpStatus: http.StatusBadRequest,
		})
	}
}
