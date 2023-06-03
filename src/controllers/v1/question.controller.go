package controllers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	httpModels "github.com/rpratama-dev/goforum/src/models/http"
	models "github.com/rpratama-dev/goforum/src/models/table"
	"github.com/rpratama-dev/goforum/src/services/database"
	"github.com/rpratama-dev/goforum/src/utils"
)

func UserQuestionIndex(c echo.Context) error {
	defer utils.DeferHandler(c)

	session := c.Get("session").(*models.Session)
	var questions []models.Question
	result := database.Conn.Preload("Tags").Preload("User").Where(map[string]interface{}{
		"user_id": session.UserID.String(),
		"is_active": true,
	}).Find(&questions)

	if (result.Error != nil) {
		panic(utils.PanicPayload{
			Message: result.Error.Error(),
			HttpStatus: http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Success get all questions by user",
		Data: questions,
	})
}

func UserQuestionStore(c echo.Context) error {
	defer utils.DeferHandler(c)
	session := c.Get("session").(*models.Session)
	// Bind user input
	var questionPayload models.QuestionPayload
	c.Bind(&questionPayload)

	// Start validation input
	errValidation := questionPayload.Validate()
	if (errValidation != nil) {
		panic(utils.PanicPayload{
			Message: "Validation Error",
			Data: errValidation,
			HttpStatus: http.StatusBadRequest,
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
		panic(utils.PanicPayload{
			Message: message,
			HttpStatus: http.StatusBadRequest,
		})
	}

	// Populate & create model
	var question models.Question
	question.Title = questionPayload.Title
	question.Content = questionPayload.Content
	question.UserID = session.User.ID
	question.Tags = &tags
	question.IsActive = true
	question.CreatedBy = &session.User.ID
	question.CreatedName = session.User.FullName
	question.CreatedFrom = c.Request().Header.Get("x-api-key")
	result := database.Conn.Create(&question)

	// Check if failed to create question
	if (result.Error != nil) {
		panic(utils.PanicPayload{
			Message: result.Error.Error(),
			HttpStatus: http.StatusInternalServerError,
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
	defer utils.DeferHandler(c)
	session := c.Get("session").(*models.Session)
	_, err := uuid.Parse(c.Param("id"))
	if err != nil {
		panic(utils.PanicPayload{
			Message: "Param must be a uuid",
			HttpStatus: http.StatusBadRequest,
		})
	}

	var question models.Question
	result := database.Conn.Preload("Tags").Preload("User").Where(map[string]interface{}{
		"id": c.Param("id"),
		"user_id": session.UserID.String(),
		"is_active": true,
	}).First(&question)

	// Check if failed to create question
	if (result.Error != nil) {
		panic(utils.PanicPayload{
			Message: result.Error.Error(),
			HttpStatus: http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Success get question by user",
		Data: question,
	})
}

func UserQuestionUpdate(c echo.Context) error {
	defer utils.DeferHandler(c)
	session := c.Get("session").(*models.Session)
	apiKey := c.Get("apiKey").(*string)
	questionId := c.Param("id")
	_, err := uuid.Parse(questionId)
	if err != nil {
		panic(utils.PanicPayload{
			Message: "Param must be a uuid",
			HttpStatus: http.StatusInternalServerError,
		})
	}

	// Bind user input
	var questionPayload models.QuestionPayload
	c.Bind(&questionPayload)

	// Find Question
	var question models.Question
	result := database.Conn.Preload("Tags").Preload("User").Where(map[string]interface{}{
		"id": questionId,
		"user_id": session.UserID.String(),
		"is_active": true,
	}).First(&question)
	if (result.Error != nil) {
		panic(utils.PanicPayload{
			Message: result.Error.Error(),
			HttpStatus: http.StatusInternalServerError,
		})
	}

	// Start validation input
	errValidation := questionPayload.Validate()
	if (errValidation != nil) {
		panic(utils.PanicPayload{
			Message: "Validation Error",
			Data: errValidation,
			HttpStatus: http.StatusBadRequest,
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
		panic(utils.PanicPayload{
			Message: message,
			HttpStatus: http.StatusBadRequest,
		})
	}

	// Remove the existing tags from the question
	if (len(*question.Tags) > 0) {
		database.Conn.Model(&question).Association("Tags").Delete(question.Tags)
	}

	// Append the new tags to the question & update question
	database.Conn.Model(&question).Association("Tags").Append(tags)
	database.Conn.Model(&question).Updates(map[string]interface{}{
		"title": questionPayload.Title,
		"content": questionPayload.Content,
		"updated_by": session.User.ID,
		"updated_name": session.User.FullName,
		"updated_from": apiKey,
	})

	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Success update question by user",
		Data: question,
	})
}

func UserQuestionDestroy(c echo.Context) error {
	defer utils.DeferHandler(c)
	session := c.Get("session").(*models.Session)
	apiKey := c.Get("apiKey").(*string)
	questionId := c.Param("id")
	_, err := uuid.Parse(questionId)
	if err != nil {
		panic(utils.PanicPayload{
			Message: "Param must be a uuid",
			HttpStatus: http.StatusBadRequest,
		})
	}

	// Find Question
	var question models.Question
	result := database.Conn.Where(map[string]interface{}{
		"id": questionId,
		"user_id": session.UserID.String(),
		"is_active": true,
	}).First(&question)
	if (result.Error != nil) {
		panic(utils.PanicPayload{
			Message: result.Error.Error(),
			HttpStatus: http.StatusInternalServerError,
		})
	}

	// Delete Question
	question.IsActive = false
	question.DeletedBy = &session.User.ID
	question.DeletedName = session.User.FullName
	question.DeletedFrom = *apiKey
	err = question.SoftDelete()
	if (err != nil) {
		panic(utils.PanicPayload{
			Message: err.Error(),
			HttpStatus: http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Success delete question by user",
		Data: nil,
	})
}

func QuestionShow(c echo.Context) error  {
	defer utils.DeferHandler(c)
	questionId, err := uuid.Parse(c.Param("question_id"))
	if err != nil {
		panic(utils.PanicPayload{
			Message: "Param must be a uuid",
			HttpStatus: http.StatusInternalServerError,
		})
	}
	var question models.Question
	err = database.Conn.
		Preload("Tags").
		Preload("User").
		Preload("Votes").
		Preload("Comments").
		Preload("Comments.User").
		Preload("Answers").
		Preload("Answers.User").
		Preload("Answers.Votes").
		Preload("Answers.Comments").
		Preload("Answers.Comments.User").
		First(&question, questionId).Error
	if err != nil {
		panic(utils.PanicPayload{
			Message: err.Error(),
			HttpStatus: http.StatusNotFound,
		})
	}

	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Success show question by id",
		Data: question,
	})
}
