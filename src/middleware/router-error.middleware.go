package middleware

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	httpModels "github.com/rpratama-dev/goforum/src/models/http"
)

func RouterErrorMiddleware (next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		// Check if the error is a router not found error
		if errors.Is(err, echo.ErrNotFound) {
			// Handle the router not found error
			return c.JSON(http.StatusCreated, httpModels.BaseResponse{
				Message: "You'r route destination doesn't exist",
				Data: nil,
			})
		}
		return err
	}
}
