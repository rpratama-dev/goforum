package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	httpModels "github.com/rpratama-dev/mymovie/src/models/http"
	"github.com/rpratama-dev/mymovie/src/utils"
)

// Middleware function to check Bearer token and verify using JWT
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the Authorization header value
		authHeader := c.Request().Header.Get("Authorization")
		// Check if the header value is empty or does not start with "Bearer "
		authHeaders := strings.Split(authHeader, " ")
		if (len(authHeaders) < 2 || authHeaders[0] != "Bearer" || len(strings.Split(authHeaders[1], ".")) < 3) {
			return c.JSON(http.StatusUnauthorized, httpModels.BaseResponse{
				Message: "Access Denied",
				Data: nil,
			})
		}

		// Extract the token from the header value & Verify the token
		claims, err := utils.VerifyJWT(authHeaders[1])
		if err != nil {
			return c.JSON(http.StatusUnauthorized, httpModels.BaseResponse{
				Message: "Access Denied, token has invalid / expired",
				Data: nil,
			})
		}

		// Validate session status
		

		// Store the claims in the context for access in subsequent handlers
		c.Set("claims", claims)
		return next(c)
	}
}
