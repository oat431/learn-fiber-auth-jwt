package utils

import (
	"github.com/gofiber/fiber/v3/log"
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword, password string) error {
	log.Info("encryped" + hashedPassword + ": password :" + password)
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
