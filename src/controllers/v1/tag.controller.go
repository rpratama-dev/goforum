package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	httpModels "github.com/rpratama-dev/mymovie/src/models/http"
	models "github.com/rpratama-dev/mymovie/src/models/table"
	"github.com/rpratama-dev/mymovie/src/services/database"
)

func TagIndex(c echo.Context) error {
	var tags []models.Tag
	database.Conn.Where("is_active = ?", true).Find(&tags)
	
	// Return response
	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Success get all tags",
		Data: tags,
	})
}

func TagStore(c echo.Context) error {
	var tag models.Tag
	c.Bind(&tag)
	session := c.Get("session").(*models.Session)
	tag.IsActive = true
	tag.CreatedBy = &session.ID
	tag.CreatedName = session.User.FullName
	tag.CreatedFrom = c.Request().Header.Get("x-api-key")
	result := database.Conn.Create(&tag)
	if (result.Error != nil) {
		// Return response
		return c.JSON(http.StatusBadRequest, httpModels.BaseResponse{
			Message: "Failed create tag",
			Data: result.Error.Error(),
		})	
	}
	// Return response
	return c.JSON(http.StatusCreated, httpModels.BaseResponse{
		Message: "Success add new tag",
		Data: tag,
	})
}

func TagShow(c echo.Context) error {
	var tag models.Tag
	err := database.Conn.First(&tag, "id = ?", c.Param("id"))
	if (err.Error != nil) {
		return c.JSON(http.StatusNotFound, httpModels.BaseResponse{
			Message: "Tag you'r looking for doesn't exist",
			Data: err.Error.Error(),
		})	
	} 
	// Return response
	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Success get tag",
		Data: tag,
	})
}
