package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

const BasePath = "/api/v1"

type Routes interface {
	Install(app *fiber.App)
}

func WithBasePath(path string) string {
	if path[0] == '/' {
		return fmt.Sprintf("%s%s", BasePath, path)
	}
	return fmt.Sprintf("%s/%s", BasePath, path)
}

