package adapters

import (
	"fmt"
	"strconv"

	"github.com/MingPV/clean-go-template/entities"
	usecases "github.com/MingPV/clean-go-template/usecases/user"
	"github.com/gofiber/fiber/v3"
)

type HttpUserHandler struct {
	userUseCase usecases.UserUseCase
}

func NewHttpUserHandler(useCase usecases.UserUseCase) *HttpUserHandler {
	return &HttpUserHandler{userUseCase: useCase}
}

// Register new user
func (h *HttpUserHandler) Register(c fiber.Ctx) error {
	user := &entities.User{}
	if err := c.Bind().Body(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if err := h.userUseCase.Register(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Clear password before returning
	user.Password = ""

	return c.Status(fiber.StatusCreated).JSON(user)
}

// Authenticate user (login)
func (h *HttpUserHandler) Login(c fiber.Ctx) error {
	loginReq := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := c.Bind().Body(&loginReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	token, user, err := h.userUseCase.Login(loginReq.Email, loginReq.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid email or password"})
	}

	user.Password = ""

	return c.JSON(fiber.Map{
		"user":  user,
		"token": token,
	})
}

// Get my user
func (h *HttpUserHandler) GetUser(c fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	id, err := strconv.Atoi(fmt.Sprint(userID))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
	}

	user, err := h.userUseCase.FindUserByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	user.Password = ""
	return c.JSON(user)
}

// Get user by ID
func (h *HttpUserHandler) FindUserByID(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id is required"})
	}

	// Convert id to int
	userID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	user, err := h.userUseCase.FindUserByID(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	user.Password = ""
	return c.JSON(user)
}

// Get all users
func (h *HttpUserHandler) FindAllUsers(c fiber.Ctx) error {
	users, err := h.userUseCase.FindAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	for i := range users {
		users[i].Password = ""
	}

	return c.JSON(users)
}
