package service

import (
	"context"
	"oat431/learn-fiber-auth-jwt/internal/model"
	"oat431/learn-fiber-auth-jwt/internal/payload/request"
	"oat431/learn-fiber-auth-jwt/internal/payload/response"
	"oat431/learn-fiber-auth-jwt/internal/repository"
	"oat431/learn-fiber-auth-jwt/pkg/common"
	"oat431/learn-fiber-auth-jwt/pkg/utils"
	"time"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v3/log"
)

type authService struct {
	repo             repository.AuthRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

type AuthService interface {
	Register(ctx context.Context, request request.RegisterRequest) (*response.AuthResponse, error)
	LoginIn(ctx context.Context, request request.LoginRequest) (*response.JWTResponse, error)
	RevokeAccess(ctx context.Context, refreshToken string) error
	GetUserDetails(ctx context.Context, authID uuid.UUID) (*response.AuthResponse, error)
}

func NewAuthService(repo repository.AuthRepository, refreshTokenRepo repository.RefreshTokenRepository) AuthService {
	return &authService{
		repo:             repo,
		refreshTokenRepo: refreshTokenRepo,
	}
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

func (s *authService) LoginIn(ctx context.Context, request request.LoginRequest) (*response.JWTResponse, error) {
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

	accessToken, err := utils.GenerateAccessToken(auth.ID)
	if err != nil {
		return nil, err
	}

	refreshTokenStr, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Save refresh token
	refreshToken := model.RefreshToken{
		BaseEntity: common.BaseEntity{
			ID:        utils.GetUUIDFromString(utils.GenerateUUID()),
			CreatedAt: utils.GetTimeFromString(utils.GetCurrentTime()),
			UpdatedAt: utils.GetTimeFromString(utils.GetCurrentTime()),
		},
		AuthID:    auth.ID,
		Token:     refreshTokenStr,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7), // 7 days
		Revoked:   false,
	}

	err = s.refreshTokenRepo.Save(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return &response.JWTResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
	}, nil
}

func (s *authService) RevokeAccess(ctx context.Context, refreshToken string) error {
	return s.refreshTokenRepo.Revoke(ctx, refreshToken)
}

func (s *authService) GetUserDetails(ctx context.Context, authID uuid.UUID) (*response.AuthResponse, error) {
	auth, err := s.repo.GetAuthByID(ctx, authID)
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		ID:       utils.UUIDToString(auth.ID),
		Username: auth.Username,
		Email:    auth.Email,
	}, nil
}
