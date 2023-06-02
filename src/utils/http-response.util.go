package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
	models "github.com/rpratama-dev/mymovie/src/models/http"
)

func PanicHandler(c echo.Context, data *interface{})  {
	var errorData interface{}
	if (data != nil) {
		errorData = *data
	}
	if r := recover(); r != nil {
		// Handle the panic here
		c.JSON(http.StatusBadRequest, models.BaseResponse{
			Message: r.(string),
			Data: errorData,
		})
	}
}
