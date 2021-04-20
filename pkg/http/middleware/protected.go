package middleware

import (
	"net/http"
	"strings"

	"github.com/class-manager/api/pkg/util/security"
	"github.com/gofiber/fiber/v2"
)

func Protected(c *fiber.Ctx) error {
	// Get authorization header
	hv := c.Get("Authorization")

	// If header is empty, return 401
	if hv == "" {
		return c.SendStatus(http.StatusUnauthorized)
	}

	t := strings.TrimPrefix(hv, "Bearer ")
	if t == hv {
		return c.SendStatus(http.StatusUnauthorized)
	}

	// Validate token
	valid, claims := security.ValidateJWT([]byte(t))
	if !valid {
		return c.SendStatus(http.StatusForbidden)
	}

	uid, exists := claims["uid"]
	if !exists {
		return c.SendStatus(http.StatusInternalServerError)
	}

	c.Locals("uid", uid)

	return c.Next()
}
