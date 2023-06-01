package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

/**
 * used to show user data
 */
func UserShow(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success get post",
		"data": "posts",
	}) 
}

/** 
 * used to update user data 
 */
func UserUpdate(c echo.Context) {

}
