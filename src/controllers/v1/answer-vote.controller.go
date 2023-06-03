package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	httpModels "github.com/rpratama-dev/goforum/src/models/http"
	models "github.com/rpratama-dev/goforum/src/models/table"
	"github.com/rpratama-dev/goforum/src/services/database"
	"github.com/rpratama-dev/goforum/src/utils"
)

func AnswerVoteStore(c echo.Context) error {
	defer utils.DeferHandler(c)
	session := c.Get("session").(*models.Session)

	// Bind user input & validate
	var votePayload models.AnswerVotePayload
	c.Bind(&votePayload)
	votePayload.QuestionID = c.Param("question_id")
	votePayload.AnswerID = c.Param("answer_id")

	// Start validation input
	errValidation := votePayload.Validate()
	if (errValidation != nil) {
		panic(utils.PanicPayload{
			Message: "Validation Error",
			Data: errValidation,
			HttpStatus: http.StatusBadRequest,
		})
	}

	// Validate question & answer is exists
	var answer models.Answer
	err := database.Conn.
		Preload("Question").
		Where(map[string]interface{}{
			"id": votePayload.AnswerID,
			"question_id": votePayload.QuestionID,
			"is_active": votePayload.QuestionID,
		}).First(&answer).Error

	if err != nil || !answer.Question.IsActive {
		message := "Your selected question has ben archived"
		if err != nil {
			message = err.Error()
		}
		panic(utils.PanicPayload{
			Message: message,
			HttpStatus: http.StatusNotFound,
		})
	}

	// Validate is question still active
	var total int64 = 0
	// Start to validate if user already answered this question
	var answerVote models.AnswerVote
	answerVote.VoteType = votePayload.VoteType
	database.Conn.
		Preload("Answer").
		Where(map[string]interface{}{
			"answer_id": votePayload.AnswerID,
			"user_id": session.User.ID,
		}).First(&answerVote).Count(&total)

	if total > 0 {
		if !answerVote.Answer.IsActive {
			panic(utils.PanicPayload{
				Message: "Unable to vote, your selected answer has been archived",
				HttpStatus: http.StatusInternalServerError,
			})
		}
		// Update record if already vote
		result := database.Conn.Model(&answerVote).Updates(map[string]interface{}{
			"vote_type": votePayload.VoteType,
			"updated_by": session.User.ID,
			"updated_name": session.User.FullName,
			"updated_from": *c.Get("apiKey").(*string),
		})
		if (result.Error != nil) {
			panic(utils.PanicPayload{
				Message: result.Error.Error(),
				HttpStatus: http.StatusInternalServerError,
			})
		}
	} else {
		// Create new record if already vote
		answerVote.Append(votePayload, *session, *c.Get("apiKey").(*string))
		result := database.Conn.Create(&answerVote)
		if (result.Error != nil) {
			panic(utils.PanicPayload{
				Message: result.Error.Error(),
				HttpStatus: http.StatusInternalServerError,
			})
		}
	}

	response := make(map[string]interface{})
	response["id"] = answerVote.ID;
	response["vote"] = answerVote.VoteType;
	response["answer_id"] = answerVote.AnswerID;

	return c.JSON(http.StatusCreated, httpModels.BaseResponse{
		Message: "Success add vote to user",
		Data: answerVote,
	})
}
