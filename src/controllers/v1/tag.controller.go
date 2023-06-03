package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	httpModels "github.com/rpratama-dev/goforum/src/models/http"
	models "github.com/rpratama-dev/goforum/src/models/table"
	"github.com/rpratama-dev/goforum/src/services/database"
	"github.com/rpratama-dev/goforum/src/utils"
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
	defer utils.DeferHandler(c)
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
		panic(utils.PanicPayload{
			Message: result.Error.Error(),
			HttpStatus: http.StatusBadRequest,
		})
	}
	// Return response
	return c.JSON(http.StatusCreated, httpModels.BaseResponse{
		Message: "Success add new tag",
		Data: tag,
	})
}

func TagShow(c echo.Context) error {
	defer utils.DeferHandler(c)
	var tag models.Tag
	err := database.Conn.First(&tag, "id = ?", c.Param("id"))
	if (err.Error != nil) {
		panic(utils.PanicPayload{
			Message: "Tag you'r looking for doesn't exist",
			HttpStatus: http.StatusNotFound,
		})
	} 
	// Return response
	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Success get tag",
		Data: tag,
	})
}
