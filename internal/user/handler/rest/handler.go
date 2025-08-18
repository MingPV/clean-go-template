package handler

import (
	"fmt"

	"github.com/MingPV/clean-go-template/internal/entities"
	"github.com/MingPV/clean-go-template/internal/user/dto"
	"github.com/MingPV/clean-go-template/internal/user/usecase"
	"github.com/MingPV/clean-go-template/pkg/apperror"
	"github.com/MingPV/clean-go-template/pkg/responses"
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
	req := new(dto.RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return responses.Error(c, apperror.ErrInvalidData)
	}

	userEntity := dto.ToUserEntity(req)
	if err := h.userUseCase.Register(userEntity); err != nil {
		return responses.Error(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ToUserResponse(userEntity))
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
	loginReq := new(dto.LoginRequest)
	if err := c.BodyParser(loginReq); err != nil {
		return responses.Error(c, apperror.ErrInvalidData)
	}

	token, userEntity, err := h.userUseCase.Login(loginReq.Email, loginReq.Password)
	if err != nil {
		return responses.ErrorWithMessage(c, apperror.ErrUnauthorized, "invalid email or password")
	}

	return c.JSON(fiber.Map{
		"user":  dto.ToUserResponse(userEntity),
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
		return responses.Error(c, apperror.ErrInvalidData)
	}

	userEntity, err := h.userUseCase.FindUserByID(fmt.Sprint(userID))
	if err != nil {
		return responses.Error(c, err)
	}

	return c.JSON(dto.ToUserResponse(userEntity))
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
		return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "id is required")
	}

	userEntity, err := h.userUseCase.FindUserByID(id)
	if err != nil {
		return responses.Error(c, err)
	}

	return c.JSON(dto.ToUserResponse(userEntity))
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
		return responses.Error(c, err)
	}

	return c.JSON(dto.ToUserResponseList(users))
}

// PatchUser godoc
// @Summary Update an user partially
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body entities.User true "User update payload"
// @Success 200 {object} entities.User
// @Router /users/{id} [patch]
func (h *HttpUserHandler) PatchUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.PatchUserRequest
	if err := c.BodyParser(&req); err != nil {
		return responses.ErrorWithMessage(c, err, "invalid request")
	}

	user := &entities.User{Name: req.Name}

	msg, err := validatePatchUser(user)
	if err != nil {
		return responses.ErrorWithMessage(c, err, msg)
	}

	updatedUser, err := h.userUseCase.PatchUser(id, user)
	if err != nil {
		return responses.Error(c, err)
	}

	return c.JSON(dto.ToUserResponse(updatedUser))
}

// DeleteUser godoc
// @Summary Delete an user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.MessageResponse
// @Router /users/{id} [delete]
func (h *HttpUserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.userUseCase.DeleteUser(id); err != nil {
		return responses.Error(c, err)
	}

	return responses.Message(c, fiber.StatusOK, "user deleted")
}

func validatePatchUser(user *entities.User) (string, error) {

	if user.Name == "" {
		return "username is invalid", apperror.ErrInvalidData
	}

	return "", nil
}
