package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	httpModels "github.com/rpratama-dev/goforum/src/models/http"
	models "github.com/rpratama-dev/goforum/src/models/table"
	"github.com/rpratama-dev/goforum/src/services/database"
	"github.com/rpratama-dev/goforum/src/utils"
)

func QuestionCommentStore(c echo.Context) error {
	defer utils.DeferHandler(c)
	session := c.Get("session").(*models.Session)

	// Bind user input & validate
	var questionCommentPayload models.QuestionCommentPayload
	c.Bind(&questionCommentPayload)
	questionCommentPayload.QuestionID = c.Param("question_id")

	// Start validation input
	errValidation := questionCommentPayload.Validate()
	if (errValidation != nil) {
		panic(utils.PanicPayload{
			Message: "Validation Error",
			Data: errValidation,
			HttpStatus: http.StatusBadRequest,
		})
	}

	// Check if question is exist & active
	var question models.Question
	result := database.Conn.Where(map[string]interface{}{
		"id": questionCommentPayload.QuestionID,
		"is_active": true,
	}).First(&question)
	if (result.Error != nil || !question.IsActive) {
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
	var questionComment models.QuestionComment
	questionComment.Append(questionCommentPayload, *session, *c.Get("apiKey").(*string))
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
	response["question_id"] = questionComment.QuestionID

	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Success add a comment for selected question",
		Data: response,
	})
}

func QuestionCommentUpdate(c echo.Context) error {
	defer utils.DeferHandler(c)
	session := c.Get("session").(*models.Session)

	// Bind user input & validate
	var questionCommentPayload models.QuestionCommentPayloadUpdate
	c.Bind(&questionCommentPayload)
	questionCommentPayload.QuestionID = c.Param("question_id")
	questionCommentPayload.CommentID = c.Param("comment_id")

	// Start validation input
	errValidation := questionCommentPayload.Validate()
	if (errValidation != nil) {
		panic(utils.PanicPayload{
			Message: "Validation Error",
			Data: errValidation,
			HttpStatus: http.StatusBadRequest,
		})
	}

	// Check if question is exist & active
	var questionComment models.QuestionComment
	result := database.Conn.Preload("Question").Where(map[string]interface{}{
		"id": questionCommentPayload.CommentID,
		"question_id": questionCommentPayload.QuestionID,
		"is_active": true,
	}).First(&questionComment)
	if (result.Error != nil || !questionComment.Question.IsActive) {
		message := "Unable to update comment for inactive question"
		if (result.Error != nil) {
			message = result.Error.Error()
		}
		panic(utils.PanicPayload{
			Message: message,
			HttpStatus: http.StatusNotFound,
		})
	}

	// Update comment
	database.Conn.Model(&questionComment).Updates(map[string]interface{}{
		"content": questionCommentPayload.Content,
		"updated_by": session.User.ID,
		"updated_name": session.User.FullName,
		"updated_from": *c.Get("apiKey").(*string),
	})

	response := make(map[string]interface{})
	response["id"] = questionComment.ID
	response["content"] = questionComment.Content
	response["question_id"] = questionComment.QuestionID
	response["questionComment"] = questionComment

	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Success update a comment for selected question",
		Data: response,
	})
}
