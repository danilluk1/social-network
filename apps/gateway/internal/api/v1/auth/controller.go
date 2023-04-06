package auth

import (
	"github.com/danilluk1/social-network/apps/gateway/internal/middlewares"
	"github.com/danilluk1/social-network/apps/gateway/internal/types"
	"github.com/gofiber/fiber/v2"
)

func Setup(router fiber.Router, services types.Services) fiber.Router {
	mw := router.Group("auth")
	mw.Post("register", handleRegister(services))
	mw.Post("login", handleLogin(services))

	return mw
}

func handleRegister(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		dto := &createUserDto{}
		err := middlewares.ValidateBody(
			c,
			services.Validator,
			services.ValidatorTranslator,
			dto,
		)
		if err != nil {
			return nil
		}

		user, err := createUserService(dto)
		if err != nil {
			return err
		}

		return c.JSON(user)
	}
}

func handleLogin(services types.Services) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var req types.LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}

		if err := services.Login(c.Context(), &req); err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).SendString("ok")
	}
}
