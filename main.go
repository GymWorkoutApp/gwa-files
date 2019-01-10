package main

import (
	"fmt"
	"github.com/GymWorkoutApp/gwap-files/handlers"
	"github.com/GymWorkoutApp/gwap-files/middlewares"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	echolog "github.com/labstack/gommon/log"
	"go.elastic.co/apm/module/apmecho"
	"os"
)

func main() {
	// Http handlers
	e := echo.New()
	e.Use(apmecho.Middleware())
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339_nano} - [${uri} - ${method}] - ${status} - ${remote_ip}\n",
	}))

	v1 := e.Group("v1")

	v1.Use(middlewares.MiddlewareAuth)
	v1.POST("/files", handlers.HandleFilesCreateRequest)
	v1.DELETE("/files/:id", handlers.HandleFilesDeleteRequest)
	v1.GET("/files/:id", handlers.HandleFilesGetRequest)
	v1.GET("/files/:id/download", handlers.HandleFilesDownloadGetRequest)

	e.Logger.SetLevel(echolog.INFO)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", os.Getenv("PORT"))))
}