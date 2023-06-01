package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rpratama-dev/mymovie/src/configs"
	appMiddleware "github.com/rpratama-dev/mymovie/src/middleware"
	"github.com/rpratama-dev/mymovie/src/routes"
	"github.com/rpratama-dev/mymovie/src/services/database"
	"github.com/rpratama-dev/mymovie/src/services/migration"
	"github.com/rpratama-dev/mymovie/src/utils"
)

func init()  {
	utils.GenerateKeyPair()
	configs.InitConfig()
	database.StartDB()
	migration.Migrate()
}

func main() {
	app := echo.New()
	app.Pre(middleware.RemoveTrailingSlash())
	// Register router API
	routes.ApiRouter(app.Group("/api"))
	// Custom middleware to handle router not found error
	app.Use(appMiddleware.RouterErrorMiddleware)
	app.Logger.Fatal(app.Start(":"+ configs.Env.Port))
}
