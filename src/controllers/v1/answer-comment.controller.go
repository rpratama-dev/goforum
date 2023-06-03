package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	httpModels "github.com/rpratama-dev/mymovie/src/models/http"
	models "github.com/rpratama-dev/mymovie/src/models/table"
	"github.com/rpratama-dev/mymovie/src/services/database"
	"github.com/rpratama-dev/mymovie/src/utils"
)

func AnswerCommentStore(c echo.Context) error {
	defer utils.DeferHandler(c)
	session := c.Get("session").(*models.Session)

	// Bind user input & validate
	var answerCommentPayload models.AnswerCommentPayload
	c.Bind(&answerCommentPayload)
	answerCommentPayload.QuestionID = c.Param("question_id")
	answerCommentPayload.AnswerID = c.Param("answer_id")

	// Start validation input
	errValidation := answerCommentPayload.Validate()
	if (errValidation != nil) {
		panic(utils.PanicPayload{
			Message: "Validation Error",
			Data: errValidation,
			HttpStatus: http.StatusBadRequest,
		})
	}

	// Check if question is exist & active
	var answer models.Answer
	result := database.Conn.Preload("Question").Where(map[string]interface{}{
		"id": answerCommentPayload.AnswerID,
		"is_active": true,
	}).First(&answer)
	if (result.Error != nil || !answer.Question.IsActive) {
		message := "Unable to add comment for inactive question"
		if (result.Error != nil) {
			message = result.Error.Error()
		}
		panic(utils.PanicPayload{
			Message: message,
			HttpStatus: http.StatusNotFound,
		})
	}

	// Create new record
	var questionComment models.AnswerComment
	questionComment.Append(answerCommentPayload, *session, *c.Get("apiKey").(*string))
	result = database.Conn.Create(&questionComment)
	if (result.Error != nil) {
		panic(utils.PanicPayload{
			Message: result.Error.Error(),
			HttpStatus: http.StatusInternalServerError,
		})
	}

	response := make(map[string]interface{})
	response["id"] = questionComment.ID
	response["content"] = questionComment.Content
	response["answer_id"] = questionComment.AnswerID

	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Success add a comment for selected question",
		Data: response,
	})
}

func AnswerCommentUpdate(c echo.Context) error {
	defer utils.DeferHandler(c)
	session := c.Get("session").(*models.Session)

	// Bind user input & validate
	var answerCommentPayload models.AnswerCommentPayloadUpdate
	c.Bind(&answerCommentPayload)
	answerCommentPayload.QuestionID = c.Param("question_id")
	answerCommentPayload.AnswerID = c.Param("answer_id")
	answerCommentPayload.CommentID = c.Param("comment_id")

	// Start validation input
	errValidation := answerCommentPayload.Validate()
	if (errValidation != nil) {
		panic(utils.PanicPayload{
			Message: "Validation Error",
			Data: errValidation,
			HttpStatus: http.StatusBadRequest,
		})
	}

	// Check if question is exist & active
	var answerComment models.AnswerComment
	result := database.Conn.Preload("Answer").Preload("Answer.Question").Where(map[string]interface{}{
		"id": answerCommentPayload.CommentID,
		"answer_id": answerCommentPayload.AnswerID,
		"is_active": true,
	}).First(&answerComment)
	if (result.Error != nil || !answerComment.Answer.IsActive || !answerComment.Answer.Question.IsActive) {
		message := "Unable to update comment for inactive question"
		if (result.Error != nil) {
			message = result.Error.Error()
		} else if !answerComment.Answer.IsActive {
			message = "Unable to update comment for inactive answer"
		}
		panic(utils.PanicPayload{
			Message: message,
			HttpStatus: http.StatusNotFound,
		})
	}

	// Update comment
	database.Conn.Model(&answerComment).Updates(map[string]interface{}{
		"content": answerCommentPayload.Content,
		"updated_by": session.User.ID,
		"updated_name": session.User.FullName,
		"updated_from": *c.Get("apiKey").(*string),
	})

	response := make(map[string]interface{})
	response["id"] = answerComment.ID
	response["content"] = answerComment.Content
	response["answer_id"] = answerComment.AnswerID

	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Success update a comment for selected question",
		Data: response,
	})
}
