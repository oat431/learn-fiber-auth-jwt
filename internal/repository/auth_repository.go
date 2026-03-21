package repository

import (
	"context"
	"oat431/learn-fiber-auth-jwt/internal/model"
	"oat431/learn-fiber-auth-jwt/internal/payload/request"
	"oat431/learn-fiber-auth-jwt/pkg/common"
	"oat431/learn-fiber-auth-jwt/pkg/utils"

	"github.com/gofiber/fiber/v3/log"
	"github.com/jmoiron/sqlx"
)

type authRepository struct {
	db *sqlx.DB
}

type AuthRepository interface {
	Register(ctx context.Context, request request.RegisterRequest) (*model.Auth, error)
	GetAuthByUsername(ctx context.Context, username string) (*model.Auth, error)
	GetAuthByID(ctx context.Context, id int) (*model.Auth, error)
	GetAuthByEmail(ctx context.Context, email string) (*model.Auth, error)
}

func NewAuthRepository(db *sqlx.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) Register(ctx context.Context, request request.RegisterRequest) (*model.Auth, error) {
	query := `INSERT INTO tb_auth (
				id,
				created_at,
				updated_at,
				deleted_at,
				username,
				email,
				"password"
			) VALUES ($1, $2, $3, $4, $5, $6, $7) 
	RETURNING id, created_at, updated_at, deleted_at, username, email, "password"`
	id := utils.GenerateUUID()
	currentTime := utils.GetCurrentTime()
	_, err := r.db.ExecContext(ctx, query, id, currentTime, currentTime, nil, request.Username, request.Email, request.Password)
	if err != nil {
		return nil, err
	}
	return &model.Auth{
		BaseEntity: common.BaseEntity{
			ID:        utils.GetUUIDFromString(id),
			CreatedAt: utils.GetTimeFromString(currentTime),
			UpdatedAt: utils.GetTimeFromString(currentTime),
			DeletedAt: nil,
		},
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}, nil
}

func (r *authRepository) GetAuthByUsername(ctx context.Context, username string) (*model.Auth, error) {
	query := `
		SELECT 
			id,
			created_at,
			updated_at,
			deleted_at,
			username,
			email,
			"password"
		FROM
			tb_auth
		WHERE
			username = $1`
	var auth model.Auth
	err := r.db.GetContext(ctx, &auth, query, username)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return &auth, nil
}

func (r *authRepository) GetAuthByID(ctx context.Context, id int) (*model.Auth, error) {
	return nil, nil
}

func (r *authRepository) GetAuthByEmail(ctx context.Context, email string) (*model.Auth, error) {
	return nil, nil
}
