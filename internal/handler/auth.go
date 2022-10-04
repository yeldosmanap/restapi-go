package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"

	"gorest-api/internal/logs"
	"gorest-api/internal/model"

	"github.com/gofiber/fiber/v2"
)

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body model.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *fiber.Ctx) error {
	logs.Log().Info("Signing up... ")

	var input model.User

	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	validate := validator.New()

	if validationErr := validate.Struct(&input); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid input body",
		})
	}

	id, err := h.services.Authorization.CreateUser(c.UserContext(), input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})

	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"id": id,
	})
}

type signInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body signInInput true "credentials"
// @Success 200 {string} string "token"
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *fiber.Ctx) error {
	logs.Log().Info("Signing in...")

	var input signInInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	token, err := h.services.Authorization.GenerateToken(c.UserContext(), input.Email, input.Password)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"token": token})
}
