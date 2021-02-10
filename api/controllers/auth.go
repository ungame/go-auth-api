package controllers

import (
	"auth-api/api/errs"
	"auth-api/api/models"
	"auth-api/api/repository"
	"auth-api/api/security"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type AuthControllers interface {
	SignUp(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	GetUser(c *fiber.Ctx) error
	GetUsers(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
}

type authControllers struct {
	usersRepository repository.UsersRepository
}

func NewAuthControllers(usersRepository repository.UsersRepository) AuthControllers {
	return &authControllers{usersRepository: usersRepository}
}

func (a *authControllers) SignUp(c *fiber.Ctx) error {
	var user models.User
	err := c.BodyParser(&user)
	if err != nil {
		return c.
			Status(http.StatusBadRequest).
			SendString(err.Error())
	}
	user.Email = NormalizeEmail(user.Email)
	if user.Email == "" {
		return c.
			Status(http.StatusBadRequest).
			SendString("email can't be empty")
	}
	exists, err := a.usersRepository.GetByEmail(user.Email)
	if err != nil && err != mgo.ErrNotFound {
		return c.
			Status(http.StatusBadRequest).
			SendString(err.Error())
	}
	if exists != nil {
		return c.
			Status(http.StatusBadRequest).
			SendString("email already exists")
	}
	if strings.TrimSpace(user.Password) == "" {
		return c.
			Status(http.StatusBadRequest).
			SendString("password can't be empty")
	}
	user.Password, err = security.EncryptPassword(user.Password)
	if err != nil {
		return c.
			Status(http.StatusBadRequest).
			SendString(err.Error())
	}
	user.Id = bson.NewObjectId()
	user.Created = time.Now()
	user.Updated = user.Created
	err = a.usersRepository.Save(&user)
	if err != nil {
		return c.
			Status(http.StatusUnprocessableEntity).
			SendString(err.Error())
	}
	return c.Status(http.StatusCreated).JSON(user)
}

func (a *authControllers) SignIn(c *fiber.Ctx) error {
	var input SignInInput
	err := c.BodyParser(&input)
	if err != nil {
		return c.
			Status(http.StatusBadRequest).
			SendString(err.Error())
	}
	user, err := a.usersRepository.GetByEmail(input.Email)
	if err != nil {
		return c.
			Status(http.StatusBadRequest).
			SendString(err.Error())
	}
	if !input.Authenticated(user.Password) {
		return c.
			Status(http.StatusBadRequest).
			SendString(errs.ErrUnauthorized.Error())
	}
	token, err := security.NewToken(user.Id.Hex())
	if err != nil {
		return c.
			Status(http.StatusBadRequest).
			SendString(err.Error())
	}
	return c.
		Status(http.StatusOK).
		JSON(&SignInOutput{User: user, Token: token})
}

func (a *authControllers) GetUser(c *fiber.Ctx) error {
	id, err := AuthRequestWithId(c)
	if err != nil {
		return c.
			Status(http.StatusUnauthorized).
			SendString(err.Error())
	}
	user, err := a.usersRepository.GetById(id)
	if err != nil {
		return c.
			Status(http.StatusBadRequest).
			SendString(err.Error())
	}
	return c.
		Status(http.StatusOK).
		JSON(user)
}
func (a *authControllers) GetUsers(c *fiber.Ctx) error {
	users, err := a.usersRepository.GetAll()
	if err != nil {
		return c.
			Status(http.StatusBadRequest).
			SendString(err.Error())
	}
	return c.
		Status(http.StatusOK).
		JSON(users)
}

func (a *authControllers) UpdateUser(c *fiber.Ctx) error {
	id, err := AuthRequestWithId(c)
	if err != nil {
		return c.
			Status(http.StatusUnauthorized).
			SendString(err.Error())
	}
	var update models.User
	err = c.BodyParser(&update)
	if err != nil {
		return c.
			Status(http.StatusBadRequest).
			SendString(err.Error())
	}
	user, err := a.usersRepository.GetById(id)
	if err != nil {
		return c.
			Status(http.StatusBadRequest).
			SendString(err.Error())
	}
	update.Email = NormalizeEmail(update.Email)
	if update.Email == "" {
		return c.
			Status(http.StatusBadRequest).
			SendString("email can't be empty")
	}
	if update.Email == user.Email {
		return c.
			Status(http.StatusOK).
			JSON(user)
	}
	exists, err := a.usersRepository.GetByEmail(update.Email)
	if err != nil && err != mgo.ErrNotFound {
		return c.
			Status(http.StatusBadRequest).
			SendString(err.Error())
	}
	if exists != nil && exists.Id != user.Id {
		return c.
			Status(http.StatusBadRequest).
			SendString("email already exists")
	}
	user.Email = update.Email
	user.Updated = time.Now()
	err = a.usersRepository.Update(user)
	if err != nil {
		return c.
			Status(http.StatusUnprocessableEntity).
			SendString(err.Error())
	}
	return c.Status(http.StatusOK).JSON(user)
}

func (a *authControllers) DeleteUser(c *fiber.Ctx) error {
	id, err := AuthRequestWithId(c)
	if err != nil {
		return c.
			Status(http.StatusUnauthorized).
			SendString(err.Error())
	}
	err = a.usersRepository.Delete(id)
	if err != nil {
		return c.
			Status(http.StatusBadRequest).
			SendString(err.Error())
	}
	c.Set("Entity", id)
	return c.SendStatus(http.StatusNoContent)
}

func AuthRequestWithId(c *fiber.Ctx) (string, error) {
	id := c.Params("id")
	if !bson.IsObjectIdHex(id) {
		return "", errs.ErrUnauthorized
	}
	tokenString := security.GetTokenFromHeader(c)
	payload, err := security.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	if id != payload.Id {
		return "", errs.ErrUnauthorized
	}
	return payload.Id, nil
}
