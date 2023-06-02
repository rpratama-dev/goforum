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
	var tagPayload models.TagPayload
	c.Bind(&tagPayload)

	var tag models.Tag
	session := c.Get("session").(*models.Session)
	tag.Name = tagPayload.Name
	tag.IsActive = true
	tag.CreatedBy = &session.User.ID
	tag.CreatedName = session.User.FullName
	tag.CreatedFrom = c.Request().Header.Get("x-api-key")
	result := database.Conn.Create(&tag)
	if (result.Error != nil) {
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
