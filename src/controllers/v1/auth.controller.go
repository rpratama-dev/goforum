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
	defer utils.DeferHandler(c)
	var userInput models.UserRegister
	c.Bind(&userInput)

	// Start Validation
	errValidation := userInput.Validate()
	if (errValidation != nil) {
		panic(utils.PanicPayload{
			Message: "Validation Error",
			Data: errValidation,
			HttpStatus: http.StatusBadRequest,
		})
	}

	// Save model user
	var userModel = models.User{}
	userModel.Append(userInput)
	result := database.Conn.Create(&userModel)
	if result.Error != nil {
		panic(utils.PanicPayload{
			Message: result.Error.Error(),
			HttpStatus: http.StatusInternalServerError,
		})
	}

	// Return response
	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "success",
		Data: userModel,
	})
}

func AuthSignIn(c echo.Context) error {
	defer utils.DeferHandler(c)
	// Bind input user
	var userInput models.UserLogin
	c.Bind(&userInput)

	// Start Validation
	errValidation := userInput.Validate()
	if (errValidation != nil) {
		panic(utils.PanicPayload{
			Message: "Validation Error",
			HttpStatus: http.StatusBadRequest,
		})
	}
	
	// Start find and validate user
	var user = models.User{}
	user.Email = userInput.Email 
	err := user.GetByEmail()
	if (err != nil) {
		panic(utils.PanicPayload{
			Message: "Invalid email / password",
			HttpStatus: http.StatusBadRequest,
		})
	}

	// Check if password is match
	isMatch := userInput.IsPasswordMatch(user.Password)
	if (!isMatch) {
		panic(utils.PanicPayload{
			Message: "Invalid email / password",
			HttpStatus: http.StatusBadRequest,
		})
	}

	// Only verified and active user can sign-in
	if (!user.IsVerified || !user.IsActive) {
		message := "You'r account is inactive please contact web administrator"
		if (!user.IsVerified) {
			message = "Please verified you'r account first, before try sign-in"
		}
		panic(utils.PanicPayload{
			Message: message,
			HttpStatus: http.StatusBadRequest,
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
	defer utils.DeferHandler(c)
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
	defer utils.DeferHandler(c)
	verifyToken := c.Param("token")
	var user models.User
	err := user.GetByToken(verifyToken);
	if (err != nil) {
		panic(utils.PanicPayload{
			Message: "Failed to verified user registration",
			HttpStatus: http.StatusBadRequest,
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
