package utils

import (
	"fmt"

	"github.com/labstack/echo/v4"
	models "github.com/rpratama-dev/mymovie/src/models/http"
)

type PanicPayload struct {
	Data	 			interface{}
	Message			string
	HttpStatus 	int
}

func DeferHandler(c echo.Context)  {
	if r := recover(); r != nil {
		pyd := r.(PanicPayload)
		var errorData interface{}
		if (pyd.Data != nil) {
			errorData = pyd.Data
		}
		// Handle the panic here
		c.JSON(pyd.HttpStatus, models.BaseResponse{
			Message: pyd.Message,
			Data: errorData,
		})
	}
	fmt.Println("Exec last")
}
