package responses

import (
	appError "github.com/MingPV/clean-go-template/pkg/apperror"
	"github.com/gofiber/fiber/v2"
)

// ErrorResponse represents the standard error response
type ErrorResponse struct {
	Error string `json:"error" example:"example error"`
}

func Error(c *fiber.Ctx, err error) error {
	return c.Status(appError.StatusCode(err)).JSON(ErrorResponse{Error: err.Error()})
}

func ErrorWithMessage(c *fiber.Ctx, err error, message string) error {
	return c.Status(appError.StatusCode(err)).JSON(ErrorResponse{Error: message})
}
