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

func AuthSignUp(c echo.Context) error {
	var userInput models.UserRegister
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
	err := user.GetByEmail()
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

	// Only verified and active user can sign-in
	if (!user.IsVerified || !user.IsActive) {
		message := "You'r account is inactive please contact web administrator"
		if (!user.IsVerified) {
			message = "Please verified you'r account first, before try sign-in"
		}
		return c.JSON(http.StatusBadRequest, httpModels.BaseResponse{
			Message: message,
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

func AuthVerify(c echo.Context) error {
	verifyToken := c.Param("token")
	var user models.User
	err := user.GetByToken(verifyToken);
	if (err != nil) {
		// Return response
		return c.JSON(http.StatusBadRequest, httpModels.BaseResponse{
			Message: "Failed to verified user registration",
			Data: nil,
		})
	}

	// Update user status verification
	currentTime := time.Now()
	user.VerifiedToken = ""
	user.VerifiedAt = &currentTime
	user.IsVerified = true
	user.IsActive = true
	user.Update("verified_token", "verified_at", "is_verified", "is_active")

	response := make(map[string]interface{})
	response["verified_at"] = user.VerifiedAt;
	response["is_verified"] = user.IsVerified;
	response["is_active"] = user.IsActive;

	// Return response
	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Success verified user registration",
		Data: response,
	})
}
