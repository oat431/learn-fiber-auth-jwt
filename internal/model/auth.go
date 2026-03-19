package model

import (
	"oat431/learn-fiber-auth-jwt/pkg/common"
)

type Auth struct {
	common.BaseEntity

	Username string `db:"username" json:"username"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"-"`
}
