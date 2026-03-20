package controller

import (
	"oat431/learn-fiber-auth-jwt/internal/payload/request"
	"oat431/learn-fiber-auth-jwt/internal/payload/response"
	"oat431/learn-fiber-auth-jwt/internal/service"
	"oat431/learn-fiber-auth-jwt/pkg/common"

	"github.com/gofiber/fiber/v3"
)

type AuthController struct {
	service service.AuthService
}

func NewAuthController(service service.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (auth *AuthController) RegisterNewUser(c fiber.Ctx) error {
	req := c.Locals("payload").(*request.RegisterRequest)
	authDto, err := auth.service.Register(c.Context(), *req)
	var response = common.ResponseDTO[response.AuthResponse]{}
	if err != nil {
		response.Data = nil
		response.Status = common.ERROR
		response.Error = &common.ResponseDTOError{
			HttpCode:  fiber.ErrBadRequest.Code,
			ErrorCode: "REGISTER-01",
			Message:   err.Error(),
		}
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	response.Data = authDto
	response.Status = common.SUCCESS
	response.Error = nil
	return c.Status(fiber.StatusCreated).JSON(response)
}
