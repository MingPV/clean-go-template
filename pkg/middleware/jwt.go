package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}
		tokenStr := auth[len("Bearer "):]

		// tokenStr := c.Cookies("token") // Assuming the token is stored in a cookie named "token"
		// if tokenStr == "" {
		// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		// }

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := claims["user_id"]
		c.Locals("user_id", userID)

		return c.Next()
	}
}
