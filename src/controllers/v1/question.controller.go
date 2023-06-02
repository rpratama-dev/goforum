package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	httpModels "github.com/rpratama-dev/mymovie/src/models/http"
)

func UserQuestionIndex(c echo.Context) error {
	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Success get all questions by user",
		Data: nil,
	})
}

func UserQuestionStore(c echo.Context) error {
	return c.JSON(http.StatusCreated, httpModels.BaseResponse{
		Message: "Success create question by user",
		Data: nil,
	})
}

func UserQuestionShow(c echo.Context) error {
	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Success get question by user",
		Data: nil,
	})
}

func UserQuestionUpdate(c echo.Context) error {
	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Success update question by user",
		Data: nil,
	})
}

func UserQuestionDestroy(c echo.Context) error {
	return c.JSON(http.StatusOK, httpModels.BaseResponse{
		Message: "Success delete question by user",
		Data: nil,
	})
}
