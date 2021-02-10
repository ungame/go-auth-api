package main

import (
	"auth-api/api/controllers"
	"auth-api/api/db"
	"auth-api/api/repository"
	"auth-api/api/routes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

var port int

func init() {
	flag.IntVar(&port, "port", 8080, "set authentication api port")
	flag.Parse()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	cfg := db.NewConfig()
	fmt.Println(cfg.Dsn())

	conn, err := db.NewConnection(cfg)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	usersRepository := repository.NewUsersRepository(conn)
	authControllers := controllers.NewAuthControllers(usersRepository)

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/", apiInfo)

	routes.
		NewAuthRoutes(authControllers).
		Install(app)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}

func apiInfo(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"ip":        c.IP(),
		"host":      c.Hostname(),
		"listening": port,
		"date":      time.Now().String(),
		"info":      api,
	})
}

type Route struct {
	URI         string `json:"uri"`
	Method      string `json:"method"`
	Description string `json:"description"`
}

type Api struct {
	Description string  `json:"description"`
	Version     float64 `json:"version"`
	Routes      []*Route
}

var api = Api{
	Description: "Authentication API",
	Version:     1.0,
	Routes: []*Route{
		{
			URI:         routes.WithBasePath("/signup"),
			Method:      http.MethodPost,
			Description: "Create one user with email and password",
		},
		{
			URI:         routes.WithBasePath("/signin"),
			Method:      http.MethodPost,
			Description: "Login with email and password",
		},
		{
			URI:         routes.WithBasePath("/users"),
			Method:      http.MethodGet,
			Description: "Show all users",
		},
		{
			URI:         routes.WithBasePath("/users/:id"),
			Method:      http.MethodGet,
			Description: "Show authenticated user info by id",
		},
		{
			URI:         routes.WithBasePath("/users/:id"),
			Method:      http.MethodPut,
			Description: "Update authenticated user",
		},
		{
			URI:         routes.WithBasePath("/users/:id"),
			Method:      http.MethodDelete,
			Description: "Delete authenticated user",
		},
	},
}
