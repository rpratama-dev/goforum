package controllers

import (
	"net/http"

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
	answerPayload.QuestionID = c.Param("question_id")

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

func AnswerUpdate(c echo.Context) error {
	defer utils.DeferHandler(c)
	session := c.Get("session").(*models.Session)

	// Bind user input & validate questionId
	var answerPayload models.AnswerPayloadUpdate
	c.Bind(&answerPayload)
	answerPayload.QuestionID = c.Param("question_id")
	answerPayload.AnswerID = c.Param("answer_id")

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
	var answer models.Answer
	result := database.Conn.Preload("Question").Where(map[string]interface{}{
		"question_id": answerPayload.QuestionID,
		"user_id": session.User.ID,
		"id": answerPayload.AnswerID,
	}).First(&answer)
	if (result.Error != nil || !answer.Question.IsActive || answer.IsTheBest) {
		message := "Can't change answer for inactive question"
		if (answer.IsTheBest) {
			message = "You'r answer mark as the best, so you can't edit"
		} else if result.Error != nil {
			message = result.Error.Error() 
		}
		panic(utils.PanicPayload{
			Message: message,
			HttpStatus: http.StatusBadRequest,
		})
	}

	// Update record
	database.Conn.Model(&answer).Updates(map[string]interface{}{
		"content": answerPayload.Content,
		"updated_by": session.User.ID,
		"updated_name": session.User.FullName,
		"updated_from": c.Get("apiKey").(*string),
	})

	response := make(map[string]interface{})
	response["content"] = answer.Content;
	response["question_id"] = answer.QuestionID;
	response["user_id"] = answer.UserID;

	return c.JSON(http.StatusCreated, httpModels.BaseResponse{
		Message: "Success create answer",
		Data: response,
	})
}
