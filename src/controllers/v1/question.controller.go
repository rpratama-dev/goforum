package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	httpModels "github.com/rpratama-dev/mymovie/src/models/http"
	models "github.com/rpratama-dev/mymovie/src/models/table"
	"github.com/rpratama-dev/mymovie/src/services/database"
)

func UserQuestionIndex(c echo.Context) error {
	session := c.Get("session").(*models.Session)
	var questions []models.Question
	result := database.Conn.Where(map[string]interface{}{
		"user_id": session.UserID.String(),
		"is_active": true,
	}).Find(&questions)

	if (result.Error != nil) {
		return c.JSON(http.StatusBadRequest, httpModels.BaseResponse{
			Message: "Failed find all question",
			Data: result.Error.Error(),
		})	
	}

	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Success get all questions by user",
		Data: questions,
	})
}

func UserQuestionStore(c echo.Context) error {
	session := c.Get("session").(*models.Session)
	// Bind user input
	var questionPayload models.QuestionPayload
	c.Bind(&questionPayload)

	// Start validation input
	errValidation := questionPayload.Validate()
	if (errValidation != nil) {
		return c.JSON(http.StatusBadRequest, httpModels.BaseResponse{
			Message: "Validation Error",
			Data: errValidation,
		})
	}

	// Validate input tags is exist
	var tags []models.Tag
	resultTag := database.Conn.Where("id IN (?)", questionPayload.Tags).Find(&tags)
	if (resultTag.Error != nil || (len(tags) != len(questionPayload.Tags))) {
		message := "Contains invalid tags, please check tag id"
		if (resultTag.Error != nil) {
			message = resultTag.Error.Error()
		}
		return c.JSON(http.StatusBadRequest, httpModels.BaseResponse{
			Message: message,
			Data: nil,
		})
	}

	// Populate & create model
	var question models.Question
	question.Title = questionPayload.Title
	question.Content = questionPayload.Content
	question.UserID = session.User.ID
	question.Tags = tags
	question.IsActive = true
	question.CreatedBy = &session.User.ID
	question.CreatedName = session.User.FullName
	question.CreatedFrom = c.Request().Header.Get("x-api-key")
	result := database.Conn.Create(&question)

	// Check if failed to create question
	if (result.Error != nil) {
		return c.JSON(http.StatusBadRequest, httpModels.BaseResponse{
			Message: "Failed create question",
			Data: result.Error.Error(),
		})	
	}

	// Output the created question with tags & user
	var savedQuestion models.Question
	database.Conn.Preload("Tags").Preload("User").First(&savedQuestion, question.ID)

	return c.JSON(http.StatusCreated, httpModels.BaseResponse{
		Message: "Success create question by user",
		Data: savedQuestion,
	})
}

func UserQuestionShow(c echo.Context) error {
	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Success get question by user",
		Data: nil,
	})
}

func UserQuestionUpdate(c echo.Context) error {
	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Success update question by user",
		Data: nil,
	})
}

func UserQuestionDestroy(c echo.Context) error {
	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Success delete question by user",
		Data: nil,
	})
}
