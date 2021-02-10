package middlewares

import (
	"auth-api/api/security"
	"net/http"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

func Authenticate(c *fiber.Ctx) error {
	cfg := jwtware.Config{
		SigningKey:    security.JwtSecretKey,
		SigningMethod: security.JwtSigningMethod.Name,
		TokenLookup:   "header:Authorization", // Authorization: Bearer Token...
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.
				Status(http.StatusUnauthorized).
				SendString(err.Error())
		},
	}
	return jwtware.New(cfg)(c)
}
