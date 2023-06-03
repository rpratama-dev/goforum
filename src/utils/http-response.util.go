package utils

import (
	"net/http"
	"runtime"

	"github.com/labstack/echo/v4"
	models "github.com/rpratama-dev/goforum/src/models/http"
)

type PanicPayload struct {
	Data	 			interface{}
	Message			string
	HttpStatus 	int
}

func DeferHandler(c echo.Context)  {
	if r := recover(); r != nil {
		var message string = "Unexpected types of errors"
		var data interface{} = nil
		var httpStatus int = http.StatusInternalServerError

		switch payload := r.(type) {
		case PanicPayload:
				// The err is an instance of utils.PanicPayload
			message = payload.Message
			httpStatus = payload.HttpStatus
			if payload.Data != nil {
				data = payload.Data
			}
		case *runtime.TypeAssertionError:
			message = payload.Error()
		case string:
			message = payload;
		}

		// Handle the panic here
		c.JSON(httpStatus, models.BaseResponse{
			Message: message,
			Data: data,
		})
	}
}
