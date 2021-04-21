package api_v1

import (
	"net/http"

	"github.com/class-manager/api/pkg/util/security"
	"github.com/gofiber/fiber/v2"
)

func Reauth(c *fiber.Ctx) error {
	crt := c.Cookies("crt_")

	if crt == "" {
		return c.SendStatus(http.StatusUnauthorized)
	}

	// Cookie exists, validate
	valid, acc := security.ValidateRefreshToken(crt, true)
	if !valid {
		security.ClearRefreshCookie(c)
		return c.SendStatus(http.StatusForbidden)
	}

	return authorize(c, acc)
}
