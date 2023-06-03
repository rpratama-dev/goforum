package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	httpModels "github.com/rpratama-dev/goforum/src/models/http"
	models "github.com/rpratama-dev/goforum/src/models/table"
	"github.com/rpratama-dev/goforum/src/services/database"
	"github.com/rpratama-dev/goforum/src/utils"
)

func QuestionVoteStore(c echo.Context) error {
	defer utils.DeferHandler(c)
	session := c.Get("session").(*models.Session)

	// Bind user input & validate
	var votePayload models.QuestionVotePayload
	c.Bind(&votePayload)
	votePayload.QuestionID = c.Param("question_id")

	// Start validation input
	errValidation := votePayload.Validate()
	if errValidation != nil {
		panic(utils.PanicPayload{
			Message: "Validation Error",
			Data: errValidation,
			HttpStatus: http.StatusBadRequest,
		})
	}

	// Validate question is exists
	var question models.Question
	err := database.Conn.
		Where(map[string]interface{}{
			"id": votePayload.QuestionID,
			"is_active": true,
		}).First(&question).Error

	if err != nil {
		panic(utils.PanicPayload{
			Message: err.Error(),
			HttpStatus: http.StatusNotFound,
		})
	}

	var total int64 = 0
	// Start to validate if user already answered this question
	var questionVote models.QuestionVote
	questionVote.VoteType = votePayload.VoteType
	err = database.Conn.
		Preload("Question").
		Where(map[string]interface{}{
			"question_id": votePayload.QuestionID,
			"user_id": session.User.ID,
		}).First(&questionVote).Count(&total).Error

	if err != nil {
		total = 0
	}

	if questionVote.Question != nil && !questionVote.Question.IsActive {
		panic(utils.PanicPayload{
			Message: "Unable to vote, your selected question has been archived",
			HttpStatus: http.StatusInternalServerError,
		})
	} else if total > 0 && questionVote.Question.IsActive {
		// Update record if already vote
		result := database.Conn.Model(&questionVote).Updates(map[string]interface{}{
			"vote_type": votePayload.VoteType,
			"updated_by": session.User.ID,
			"updated_name": session.User.FullName,
			"updated_from": *c.Get("apiKey").(*string),
		})
		if result.Error != nil {
			panic(utils.PanicPayload{
				Message: result.Error.Error(),
				HttpStatus: http.StatusInternalServerError,
			})
		}
		
	} else {
		// Create new record if already vote
		questionVote.Append(votePayload, *session, *c.Get("apiKey").(*string))
		result := database.Conn.Create(&questionVote)

		if result.Error != nil {
			panic(utils.PanicPayload{
				Message: result.Error.Error(),
				HttpStatus: http.StatusInternalServerError,
			})
		}
	}

	response := make(map[string]interface{})
	response["id"] = questionVote.ID;
	response["vote"] = questionVote.VoteType;
	response["question_id"] = questionVote.QuestionID;

	return c.JSON(http.StatusCreated, httpModels.BaseResponse{
		Message: "Success add vote to user",
		Data: questionVote,
	})
}
