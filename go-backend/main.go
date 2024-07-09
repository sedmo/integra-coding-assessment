package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/sedmo/integra-coding-assessment/go-backend/db"
	"github.com/sedmo/integra-coding-assessment/go-backend/db/connectors"
	_ "github.com/sedmo/integra-coding-assessment/go-backend/docs"
	"github.com/sedmo/integra-coding-assessment/go-backend/handlers"
)

// @title Integra API
// @version 1.0
// @description This is a server for the integra coding assessment.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:1323
// @BasePath /
func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Enable CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:4200"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	pgConnector := &connectors.PostgresConnector{}
    db.InitDB(pgConnector)

	// Routes
	e.GET("/users", handlers.GetUsers)
	e.POST("/users", handlers.CreateUser)
	e.PUT("/users/:id", handlers.UpdateUser)
	e.DELETE("/users/:id", handlers.DeleteUser)

	e.Logger.Fatal(e.Start(":1323"))
}
