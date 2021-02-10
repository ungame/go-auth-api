package routes

import (
	"auth-api/api/controllers"
	"auth-api/api/middlewares"

	"github.com/gofiber/fiber/v2"
)

type authRoutes struct {
	authControllers controllers.AuthControllers
}

func NewAuthRoutes(authControllers controllers.AuthControllers) Routes {
	return &authRoutes{authControllers: authControllers}
}

func (r *authRoutes) Install(app *fiber.App) {
	app.Post(WithBasePath("/signup"), r.authControllers.SignUp)
	app.Post(WithBasePath("/signin"), r.authControllers.SignIn)
	app.Get(WithBasePath("/users"), middlewares.Authenticate, r.authControllers.GetUsers)
	app.Get(WithBasePath("/users/:id"), middlewares.Authenticate, r.authControllers.GetUser)
	app.Put(WithBasePath("/users/:id"), middlewares.Authenticate, r.authControllers.UpdateUser)
	app.Delete(WithBasePath("/users/:id"), middlewares.Authenticate, r.authControllers.DeleteUser)
}
