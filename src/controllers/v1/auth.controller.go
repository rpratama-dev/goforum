package controllers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	httpModels "github.com/rpratama-dev/mymovie/src/models/http"
	models "github.com/rpratama-dev/mymovie/src/models/table"
	"github.com/rpratama-dev/mymovie/src/services/database"
	"github.com/rpratama-dev/mymovie/src/utils"
)

func AuthSignIn(c echo.Context) error {
	// Bind input user
	var userInput models.UserLogin
	c.Bind(&userInput)

	// Start Validation
	errValidation := userInput.Validate()
	if (errValidation != nil) {
		return c.JSON(http.StatusBadRequest, httpModels.BaseResponse{
			Message: "Validation Error",
			Data: errValidation,
		})
	}
	
	// Start find and validate user
	var user = models.User{}
	user.Email = userInput.Email 
	err := user.FindByEmail()
	if (err != nil) {
		return c.JSON(http.StatusBadRequest, httpModels.BaseResponse{
			Message: "Invalid email / password",
			Data: nil,
		})
	}

	// Check if password is match
	isMatch := userInput.IsPasswordMatch(user.Password)
	if (!isMatch) {
		return c.JSON(http.StatusBadRequest, httpModels.BaseResponse{
			Message: "Invalid email / password b",
			Data: nil,
		})
	}

	// Generate access token
	var claim = utils.ClaimPayload{}
	claim.Name = user.FullName;
	claim.UserID = user.ID.String();
	claim.SessionID = "test session id";
	claim.UserName = user.Email;
	token, claims, _ := utils.GenerateJWT(claim);

	// Convert Unix timestamp to time.Time in UTC
	t := time.Unix(claims.ExpiresAt, 0).UTC()
	response := map[string]string{
		"access_token": token,
		"expired_at": t.Format("2006-01-02T15:04:051Z"),
	}

	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "success",
		Data: response,
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
