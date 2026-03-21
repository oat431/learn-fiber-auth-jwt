package service

import (
	"context"
	"oat431/learn-fiber-auth-jwt/internal/payload/request"
	"oat431/learn-fiber-auth-jwt/internal/payload/response"
	"oat431/learn-fiber-auth-jwt/internal/repository"
	"oat431/learn-fiber-auth-jwt/pkg/utils"

	"github.com/gofiber/fiber/v3/log"
)

type authService struct {
	repo repository.AuthRepository
}

type AuthService interface {
	Register(ctx context.Context, request request.RegisterRequest) (*response.AuthResponse, error)
	LoginIn(ctx context.Context, request request.LoginRequest) (*response.AuthResponse, error)
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Register(ctx context.Context, request request.RegisterRequest) (*response.AuthResponse, error) {
	encryptedPassword, err := utils.EncryptPassword(request.Password)
	if err != nil {
		return nil, err
	}
	request.Password = encryptedPassword
	auth, err := s.repo.Register(ctx, request)
	if err != nil {
		return nil, err
	}
	return &response.AuthResponse{
		ID:       utils.UUIDToString(auth.ID),
		Username: auth.Username,
		Email:    auth.Email,
	}, nil
}

func (s *authService) LoginIn(ctx context.Context, request request.LoginRequest) (*response.AuthResponse, error) {
	auth, err := s.repo.GetAuthByUsername(ctx, request.Username)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	err = utils.ComparePassword(auth.Password, request.Password)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return &response.AuthResponse{
		ID:       utils.UUIDToString(auth.ID),
		Username: auth.Username,
		Email:    auth.Email,
	}, nil
}
