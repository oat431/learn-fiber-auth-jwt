package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
)

func RegisterHealthRoutes(router fiber.Router) {
	route := router.Group("/health-check")

	route.Get("/health", healthcheck.New())
}
