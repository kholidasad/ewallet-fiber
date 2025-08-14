package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		h := c.Get("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error":"missing token"})
		}
		token := strings.TrimPrefix(h, "Bearer ")
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) { return []byte(secret), nil })
		if err != nil { return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error":"invalid token"}) }
		v, ok := claims["sub"].(float64); if !ok { return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error":"invalid token"}) }
		c.Locals("userID", uint(v))
		return c.Next()
	}
}
