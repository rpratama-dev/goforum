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

func UserQuestionIndex(c echo.Context) error {
	defer utils.PanicHandler(c, nil)
	session := c.Get("session").(*models.Session)
	var questions []models.Question
	result := database.Conn.Preload("Tags").Preload("User").Where(map[string]interface{}{
		"user_id": session.UserID.String(),
		"is_active": true,
	}).Find(&questions)

	if (result.Error != nil) {
		panic(result.Error.Error())
	}

	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Success get all questions by user",
		Data: questions,
	})
}

func UserQuestionStore(c echo.Context) error {
	var errorData interface{}
	defer utils.PanicHandler(c, &errorData)
	session := c.Get("session").(*models.Session)
	// Bind user input
	var questionPayload models.QuestionPayload
	c.Bind(&questionPayload)

	// Start validation input
	errValidation := questionPayload.Validate()
	if (errValidation != nil) {
		errorData = errValidation
		panic("Validation Error")
	}

	// Validate input tags is exist
	var tags []models.Tag
	resultTag := database.Conn.Where("id IN (?)", questionPayload.Tags).Find(&tags)
	if (resultTag.Error != nil || (len(tags) != len(questionPayload.Tags))) {
		message := "Contains invalid tags, please check tag id"
		if (resultTag.Error != nil) {
			message = resultTag.Error.Error()
		}
		panic(message)
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
		panic(result.Error.Error())
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
	defer utils.PanicHandler(c, nil)
	session := c.Get("session").(*models.Session)
	_, err := uuid.Parse(c.Param("id"))
	if err != nil {
		panic("Param must be a uuid")
	}

	var question models.Question
	result := database.Conn.Preload("Tags").Preload("User").Where(map[string]interface{}{
		"id": c.Param("id"),
		"user_id": session.UserID.String(),
		"is_active": true,
	}).First(&question)

	// Check if failed to create question
	if (result.Error != nil) {
		panic(result.Error.Error())
	}

	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Success get question by user",
		Data: question,
	})
}

func UserQuestionUpdate(c echo.Context) error {
	var errorData interface{}
	defer utils.PanicHandler(c, &errorData)
	session := c.Get("session").(*models.Session)
	apiKey := c.Get("apiKey").(*string)
	questionId := c.Param("id")
	_, err := uuid.Parse(questionId)
	if err != nil {
		panic("Param must be a uuid")
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
		panic(result.Error.Error())
	}

	// Start validation input
	errValidation := questionPayload.Validate()
	if (errValidation != nil) {
		errorData = errValidation
		panic("Validation Error")
	}

	// Validate input tags is exist
	var tags []models.Tag
	resultTag := database.Conn.Where("id IN (?)", questionPayload.Tags).Find(&tags)
	if (resultTag.Error != nil || (len(tags) != len(questionPayload.Tags))) {
		message := "Contains invalid tags, please check tag id"
		if (resultTag.Error != nil) {
			message = resultTag.Error.Error()
		}
		panic(message)
	}

	// Remove the existing tags from the question
	if (len(question.Tags) > 0) {
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
	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Success delete question by user",
		Data: nil,
	})
}
