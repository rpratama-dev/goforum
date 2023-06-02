package middleware

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rpratama-dev/mymovie/src/utils"
)

func RouterErrorMiddleware (next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		defer utils.PanicHandler(c)

		// Check if the error is a router not found error
		if errors.Is(err, echo.ErrNotFound) {
			// Handle the router not found error
			panic(utils.PanicPayload{
				Message: "You'r route destination doesn't exist",
				HttpStatus: http.StatusNotFound,
			})
		}
		return err
	}
}
