package controllers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	httpModels "github.com/rpratama-dev/mymovie/src/models/http"
	models "github.com/rpratama-dev/mymovie/src/models/table"
	"github.com/rpratama-dev/mymovie/src/services/database"
	"github.com/rpratama-dev/mymovie/src/utils"
)

func AnswerStore(c echo.Context) error {
	defer utils.DeferHandler(c)
	session := c.Get("session").(*models.Session)
	apiKey := c.Get("apiKey").(*string)

	// Bind user input & validate questionId
	var answerPayload models.AnswerPayload
	c.Bind(&answerPayload)
	questionId, err := uuid.Parse(c.Param("questionId"))
	if err != nil {
		panic(utils.PanicPayload{
			Message: "Param must be a uuid",
			HttpStatus: http.StatusBadRequest,
		})
	}
	answerPayload.QuestionID = questionId

	// Start validation input
	errValidation := answerPayload.Validate()
	if (errValidation != nil) {
		panic(utils.PanicPayload{
			Message: "Validation Error",
			Data: errValidation,
			HttpStatus: http.StatusBadRequest,
		})
	}

	// Start to validate if user already answered this question
	var answerExist models.Answer
	var count int64 = 0
	database.Conn.Where(map[string]interface{}{
		"question_id": answerPayload.QuestionID,
		"user_id": session.User.ID,
	}).First(&answerExist).Count(&count)
	if (count > 0) {
		panic(utils.PanicPayload{
			Message: "You have already provided an answer to this question",
			HttpStatus: http.StatusBadRequest,
		})
	}

	// Start to validate question is exist & active
	var question models.Question
	result := database.Conn.Where(map[string]interface{}{
		"id": answerPayload.QuestionID,
		"is_active": true,
	}).First(&question)
	if (result.Error != nil) {
		panic(utils.PanicPayload{
			Message: result.Error.Error(),
			HttpStatus: http.StatusInternalServerError,
		})
	}

	// Create answer
	var answer models.Answer
	answer.Append(answerPayload, *session, *apiKey)
	result = database.Conn.Create(&answer)
	if (result.Error != nil) {
		panic(utils.PanicPayload{
			Message: result.Error.Error(),
			HttpStatus: http.StatusInternalServerError,
		})
	}

	response := make(map[string]interface{})
	response["content"] = answer.Content;
	response["question_id"] = answer.QuestionID;
	response["user_id"] = answer.UserID;

	return c.JSON(http.StatusCreated, httpModels.BaseResponse{
		Message: "Success create answer",
		Data: response,
	})
}
