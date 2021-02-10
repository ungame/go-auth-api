package controllers

import (
	"auth-api/api/models"
	"auth-api/api/security"
	"log"
	"strings"
)

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (i *SignInInput) Cleanup() {
	i.Email = strings.TrimSpace(strings.ToLower(i.Email))
}

func (i *SignInInput) Authenticated(hashedPassword string) bool {
	err := security.VerifyPassword(hashedPassword, i.Password)
	if err != nil {
		log.Println("verify password failed:", err)
		return false
	}
	return true
}

type SignInOutput struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

func NormalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}
