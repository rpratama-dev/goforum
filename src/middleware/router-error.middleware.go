package middleware

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RouterErrorMiddleware (next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)

		// Check if the error is a router not found error
		if errors.Is(err, echo.ErrNotFound) {
			// Handle the router not found error
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"message": "Route not found",
				"data": nil,
			})
		}

		return err
	}
}
