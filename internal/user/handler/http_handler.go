package handler

import (
	"fmt"
	"strconv"

	"github.com/MingPV/clean-go-template/internal/entities"
	"github.com/MingPV/clean-go-template/internal/user/usecase"
	response "github.com/MingPV/clean-go-template/pkg/responses"
	"github.com/gofiber/fiber/v2"
)

type HttpUserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewHttpUserHandler(useCase usecase.UserUseCase) *HttpUserHandler {
	return &HttpUserHandler{userUseCase: useCase}
}

// Register godoc
// @Summary Register a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body entities.User true "User registration payload"
// @Success 201 {object} entities.User
// @Router /auth/signup [post]
func (h *HttpUserHandler) Register(c *fiber.Ctx) error {
	user := &entities.User{}
	if err := c.BodyParser(&user); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request")

	}

	if err := h.userUseCase.Register(user); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())

	}

	// Clear password before returning
	user.Password = ""

	return c.Status(fiber.StatusCreated).JSON(user)
}

// Login godoc
// @Summary Authenticate user and return token
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body map[string]string true "Login credentials (email & password)"
// @Success 200 {object} map[string]interface{} "Authenticated user and JWT token"
// @Router /auth/signin [post]
func (h *HttpUserHandler) Login(c *fiber.Ctx) error {
	loginReq := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := c.BodyParser(&loginReq); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request")

	}

	token, user, err := h.userUseCase.Login(loginReq.Email, loginReq.Password)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "invalid email or password")
	}

	user.Password = ""

	return c.JSON(fiber.Map{
		"user":  user,
		"token": token,
	})
}

// GetUser godoc
// @Summary Get currently authenticated user
// @Tags users
// @Produce json
// @Success 200 {object} entities.User
// @Router /users/me [get]
func (h *HttpUserHandler) GetUser(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return response.Error(c, fiber.StatusUnauthorized, "invalid email or password")
	}

	id, err := strconv.Atoi(fmt.Sprint(userID))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid user id")
	}

	user, err := h.userUseCase.FindUserByID(uint(id))
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "user not found")
	}

	user.Password = ""
	return c.JSON(user)
}

// FindUserByID godoc
// @Summary Get user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} entities.User
// @Router /users/{id} [get]
func (h *HttpUserHandler) FindUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(c, fiber.StatusBadRequest, "id is required")
	}

	// Convert id to int
	userID, err := strconv.Atoi(id)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid user id")
	}

	user, err := h.userUseCase.FindUserByID(uint(userID))
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "user not found")
	}

	user.Password = ""
	return c.JSON(user)
}

// FindAllUsers godoc
// @Summary Get all users
// @Tags users
// @Produce json
// @Success 200 {array} entities.User
// @Router /users [get]
func (h *HttpUserHandler) FindAllUsers(c *fiber.Ctx) error {
	users, err := h.userUseCase.FindAllUsers()
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "internal server error")
	}

	for i := range users {
		users[i].Password = ""
	}

	return c.JSON(users)
}
