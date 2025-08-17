package responses

import "github.com/gofiber/fiber/v2"

// Message Response represents the standard message response
type MessageResponse struct {
	Message string `json:"message" example:"example message"`
}

func Message(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(MessageResponse{Message: message})
}
