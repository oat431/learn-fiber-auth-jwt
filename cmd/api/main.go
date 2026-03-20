package main

import (
	"log"
	"oat431/learn-fiber-auth-jwt/internal/bootstrap"
	"oat431/learn-fiber-auth-jwt/internal/config"
	"oat431/learn-fiber-auth-jwt/internal/router"
	"os"

	"github.com/gofiber/fiber/v3"
)

func main() {
	config.LoadEnvConfig()
	db := config.StartDatabase()
	defer db.Close()

	apiContainer := bootstrap.NewAPIContainer(db)

	app := fiber.New()
	router.SetupRoutes(app, apiContainer)

	port := os.Getenv("PORT")
	err := app.Listen(":" + port)
	if err != nil {
		log.Fatal("port :" + port + " is already in use")
	}
}
