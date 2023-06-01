package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rpratama-dev/mymovie/src/configs"
	"github.com/rpratama-dev/mymovie/src/utils"
)

func init()  {
	utils.GenerateKeyPair()
	configs.InitConfig()
}

func main() {
	app := echo.New()
	app.Pre(middleware.RemoveTrailingSlash())
	app.Logger.Fatal(app.Start(":"+ configs.Env.Port))
}
