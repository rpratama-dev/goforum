package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	httpModels "github.com/rpratama-dev/mymovie/src/models/http"
	models "github.com/rpratama-dev/mymovie/src/models/table"
	"github.com/rpratama-dev/mymovie/src/services/database"
)

func AuthSignIn(c echo.Context) error {
	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "success",
		Data: "User Sign In",
	})
}

func AuthSignUp(c echo.Context) error {
	var userInput models.UserPayload
	c.Bind(&userInput)

	// Start Validation
	errValidation := userInput.Validate()
	if (errValidation != nil) {
		return c.JSON(http.StatusBadRequest, httpModels.BaseResponse{
			Message: "Validation Error",
			Data: errValidation,
		})
	}

	// Save model
	var userModel = models.User{}
	userModel.Append(userInput)
	result := database.Conn.Create(&userModel)
	if result.Error != nil {
		return c.JSON(http.StatusBadRequest, httpModels.BaseResponse{
			Message: result.Error.Error(),
			Data: nil,
		})
	}

	// Return response
	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "success",
		Data: userModel,
	})
}

func AuthChangePassword() {
	//
}

func AuthGetSession() {
	//
}
