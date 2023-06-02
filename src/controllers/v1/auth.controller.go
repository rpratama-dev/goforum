package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	httpModels "github.com/rpratama-dev/mymovie/src/models/http"
	models "github.com/rpratama-dev/mymovie/src/models/table"
	"github.com/rpratama-dev/mymovie/src/services/database"
	"github.com/rpratama-dev/mymovie/src/utils"
)

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

	// Save model user
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

	// Create Session
	var sessionPayload models.SessionPayload
	sessionPayload.UserID = user.ID
	sessionPayload.IPAddress = c.RealIP()
	sessionPayload.UserAgent = c.Request().UserAgent()
	sessionPayload.FullName = user.FullName
	var sessionModel models.Session
	sessionModel.Append(sessionPayload)
	database.Conn.Create(&sessionModel)

	// Generate access token
	var claim = utils.ClaimPayload{}
	claim.Name = user.FullName
	claim.UserID = user.ID
	claim.SessionID = sessionModel.ID
	claim.UserName = user.Email
	claim.ExpiresAt = sessionModel.ExpiredAt.Unix()
	token, _, _ := utils.GenerateJWT(claim)

	// Convert Unix timestamp to time.Time in UTC
	response := make(map[string]interface{})
	response["access_token"] = token;
	response["expired_at"] = sessionModel.ExpiredAt;

	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Sign in success",
		Data: response,
	})
}

func AuthSignOut(c echo.Context) error  {
	session := c.Get("session").(*models.Session)
	session.DeletedBy = &session.User.ID;
	session.DeletedName = session.User.FullName;
	session.DeletedFrom = "User Sign Out";
	session.SoftDelete()

	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Sign out success",
		Data: session,
	})
}

func AuthChangePassword() {
	//
}

func AuthGetSession() {
	//
}
