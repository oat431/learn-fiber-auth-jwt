package bootstrap

import (
	"oat431/learn-fiber-auth-jwt/internal/controller"
	"oat431/learn-fiber-auth-jwt/internal/repository"
	"oat431/learn-fiber-auth-jwt/internal/service"

	"github.com/gofiber/fiber/v3/log"
	"github.com/jmoiron/sqlx"
)

type APIContainer struct {
	AuthController *controller.AuthController
}

func NewAPIContainer(db *sqlx.DB) *APIContainer {
	log.Info("Registering Auth Repository")
	authRepository := repository.NewAuthRepository(db)

	log.Info("Registering Auth Service")
	authService := service.NewAuthService(authRepository)

	log.Info("Registering Auth Controller")
	authController := controller.NewAuthController(authService)

	log.Info("Registered All API")
	return &APIContainer{
		AuthController: authController,
	}
}
