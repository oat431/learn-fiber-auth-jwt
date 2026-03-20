package router

import (
	"oat431/learn-fiber-auth-jwt/internal/bootstrap"
	"oat431/learn-fiber-auth-jwt/internal/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func init() {
	log.Info("Initializing routes...")
}

func SetupRoutes(app *fiber.App, apiContainer *bootstrap.APIContainer) {
	app.Use(middleware.GlobalLogger)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	RegisterHealthRoutes(v1)
	RegisterAuthRoutes(v1, apiContainer.AuthController)
}
