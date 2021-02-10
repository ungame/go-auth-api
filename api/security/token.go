package security

import (
	"fmt"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

// https://sodocumentation.net/go/topic/10161/jwt-authorization-in-go

var JwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
var JwtSigningMethod = jwt.SigningMethodHS256

const jwtExpiration = time.Minute * 5

func NewToken(userId string) (string, error) {
	claims := &jwt.StandardClaims{
		Id:        userId,
		Issuer:    userId,
		ExpiresAt: time.Now().Add(jwtExpiration).Unix(),
		IssuedAt:  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(JwtSigningMethod, claims)
	return token.SignedString(JwtSecretKey)
}

func validateSigningMethod(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	return JwtSecretKey, nil
}

func ParseToken(tokenString string) (*jwt.StandardClaims, error) {
	claims := new(jwt.StandardClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, validateSigningMethod)
	if err != nil {
		return nil, err
	}
	var ok bool
	claims, ok = token.Claims.(*jwt.StandardClaims)
	if !token.Valid || !ok {
		return nil, fmt.Errorf("invalid token: %v", tokenString)
	}
	return claims, nil
}

func GetTokenFromHeader(c *fiber.Ctx) string {
	header := c.Get("Authorization") // Bearer Token...
	return strings.Split(header, " ")[1]
}
